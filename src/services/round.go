package services

import (
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"nokib/campwiz/database"
	"slices"
	"strings"
)

type RoundService struct {
}
type RoundRequest struct {
	CampaignID  string `json:"campaignId"`
	CreatedByID string `json:"-"`
	database.CampaignRoundWritable
}
type ImportStatus string

const (
	ImportStatusSuccess ImportStatus = "success"
	ImportStatusFailed  ImportStatus = "failed"
	ImportStatusPending ImportStatus = "pending"
)

type RoundImportSummary struct {
	Status       ImportStatus `json:"status"`
	SuccessCount int          `json:"successCount"`
	FailedCount  int          `json:"failedCount"`
	FailedIds    []string     `json:"failedIds"`
}
type RoundCreateRequest struct {
	CreatedBy string `json:"-"`       // User ID who created the batch
	RoundId   string `json:"roundId"` // Round ID to which the batch belongs
}
type ImportFromCommonsPayload struct {
	// Categories from which images will be fetched
	Categories []string `json:"categories"`
}

type Jury struct {
	ID            uint64 `json:"id" gorm:"primaryKey"`
	totalAssigned int
}
type Evaluation struct {
	JuryID            uint64 `json:"juryId"`
	ImageID           uint64 `json:"imageId"`
	DistributionRound int    `json:"distributionRound"`
	Name              string `json:"name"`
}
type ByAssigned []*Jury

func (a ByAssigned) Len() int           { return len(a) }
func (a ByAssigned) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAssigned) Less(i, j int) bool { return a[i].totalAssigned < a[j].totalAssigned }

func NewRoundService() *RoundService {
	return &RoundService{}
}
func (s *RoundService) CreateRound(request *RoundRequest) (*database.CampaignRound, error) {
	round_repo := database.NewCampaignRoundRepository()
	campaign_repo := database.NewCampaignRepository()
	conn, close := database.GetDB()
	defer close()
	tx := conn.Begin()
	campaign, err := campaign_repo.FindByID(tx, request.CampaignID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("campaign not found")
	}
	round := &database.CampaignRound{
		RoundID:               GenerateID("r"),
		CreatedByID:           request.CreatedByID,
		CampaignID:            campaign.CampaignID,
		CampaignRoundWritable: request.CampaignRoundWritable,
	}
	round, err = round_repo.Create(tx, round)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return round, nil
}
func (s *RoundService) ListAllRounds(filter *database.RoundFilter) ([]database.CampaignRound, error) {
	round_repo := database.NewCampaignRoundRepository()
	conn, close := database.GetDB()
	defer close()
	rounds, err := round_repo.FindAll(conn, filter)
	if err != nil {
		return nil, err
	}
	return rounds, nil
}

func (b *RoundService) ImportFromCommons(roundId string, categories []string) (*RoundImportSummary, error) {
	commons_repo := database.NewCommonsRepository()
	round_repo := database.NewCampaignRoundRepository()
	conn, close := database.GetDB()
	defer close()
	tx := conn.Begin()
	round, err := round_repo.FindByID(tx, roundId)
	if err != nil {
		tx.Rollback()
		return nil, err
	} else if round == nil {
		tx.Rollback()
		return nil, fmt.Errorf("round not found")
	}
	submissions := []database.Submission{}
	successCount := 0
	failedCount := 0
	failedimages := []string{}
	for _, category := range categories {
		images, failedImages := commons_repo.GetImagesFromCommonsCategories(category)
		if images == nil {
			log.Println("No images found in the category: ", category)
			return nil, fmt.Errorf("no images found in the categories")
		}
		if failedImages != nil {
			failedCount += len(failedImages)
		}
		failedimages = append(failedimages, failedImages...)
		successCount += len(images)
		// submission_repo := database.NewSubmissionRepository()
		participants := map[string]string{}
		for _, image := range images {
			participants[image.UploaderUsername] = GenerateID("user")
		}
		part_repo := database.NewParticipantRepository()
		username2IdMap, err := part_repo.EnsureExists(tx, participants)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, image := range images {
			uploaderId := username2IdMap[image.UploaderUsername]
			submission := database.Submission{
				SubmissionID:  GenerateID("sub"),
				Name:          image.Name,
				CampaignID:    round.CampaignID,
				URL:           image.URL,
				Author:        image.UploaderUsername,
				SubmittedByID: uploaderId,
				ParticipantID: uploaderId,
				SubmittedAt:   image.SubmittedAt,
				MediaSubmission: database.MediaSubmission{
					MediaType:   database.MediaType(image.MediaType),
					ThumbURL:    image.URL,
					ThumbWidth:  image.Width,
					ThumbHeight: image.Height,
					License:     strings.ToUpper(image.License),
					CreditHTML:  image.CreditHTML,
					Description: image.Description,
					AudioVideoSubmission: database.AudioVideoSubmission{
						Duration: image.Duration,
						Size:     image.Size,
						Bitrate:  0,
					},
					ImageSubmission: database.ImageSubmission{
						Width:  image.Width,
						Height: image.Height,
					},
				},
			}
			submissions = append(submissions, submission)
		}
	}
	res := tx.Create(submissions)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}
	round.TotalSubmissions += successCount
	tx.Save(round)
	tx.Commit()
	return &RoundImportSummary{
		SuccessCount: round.TotalSubmissions,
		FailedCount:  failedCount,
		FailedIds:    failedimages,
	}, nil
}
func (b *RoundService) DistributeTaskAmongExistingJuries(images []database.ImageResult) {
	juries := []*Jury{}
	for i := 1; i <= 100; i++ {
		juries = append(juries, &Jury{ID: uint64(i), totalAssigned: rand.IntN(100)})
	}
	evaluations := []Evaluation{}
	imageCount, juryCount, evaluationCountRequired := len(images), len(juries), 10
	// datasetIndex := 0
	toleranceCount := 100
	if toleranceCount == 0 {
		fmt.Println("Tolerance count cannot be zero. Setting it to 1")
		toleranceCount = 1
	}
	sortedJuryByAssigned := ByAssigned(juries)
	slices.SortStableFunc(sortedJuryByAssigned, func(a, b *Jury) int {
		if a.totalAssigned < b.totalAssigned {
			return -1
		}
		if a.totalAssigned > b.totalAssigned {
			return 1
		}
		return 0
	})
	for i := 0; i < imageCount; i++ {
		// check if the last considered jury has been assigned the maximum number of images
		if evaluationCountRequired < juryCount && i%toleranceCount == 0 {
			firstUnassignedJuryIndex := evaluationCountRequired
			swapped := false
			for pivot := firstUnassignedJuryIndex; pivot > 0; pivot-- {
				for k := pivot; k < juryCount; k++ {
					if sortedJuryByAssigned[k-1].totalAssigned < sortedJuryByAssigned[k].totalAssigned {
						break
					}
					// swap the juries
					sortedJuryByAssigned[k-1], sortedJuryByAssigned[k] = sortedJuryByAssigned[k], sortedJuryByAssigned[k-1]
					swapped = true
				}
				if !swapped {
					break
				}
			}
		}
		for j := 0; j < evaluationCountRequired; j++ {
			evaluations = append(evaluations, Evaluation{
				JuryID:            sortedJuryByAssigned[j].ID,
				ImageID:           images[i].ID,
				Name:              images[i].Name,
				DistributionRound: j + 1,
			})
			sortedJuryByAssigned[j].totalAssigned++
		}
	}
	groupByJuryID := make(map[uint64][]Evaluation)
	for _, evaluation := range evaluations {
		groupByJuryID[evaluation.JuryID] = append(groupByJuryID[evaluation.JuryID], evaluation)
	}
	for j := range juryCount {
		fmt.Printf("Jury %d has %d images\n", sortedJuryByAssigned[j].ID, len(groupByJuryID[sortedJuryByAssigned[j].ID]))
	}
}
