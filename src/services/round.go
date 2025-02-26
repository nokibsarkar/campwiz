package services

import (
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"nokib/campwiz/database"
	"slices"
	"strings"

	"gorm.io/datatypes"
)

type RoundService struct {
}
type RoundRequest struct {
	CampaignID  database.IDType `json:"campaignId"`
	CreatedByID database.IDType `json:"-"`
	database.RoundWritable
}

type RoundImportSummary struct {
	Status       database.RoundStatus `json:"status"`
	SuccessCount int                  `json:"successCount"`
	FailedCount  int                  `json:"failedCount"`
	FailedIds    []string             `json:"failedIds"`
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
func (s *RoundService) CreateRound(request *RoundRequest) (*database.Round, error) {
	round_repo := database.NewRoundRepository()
	campaign_repo := database.NewCampaignRepository()
	conn, close := database.GetDB()
	defer close()
	tx := conn.Begin()
	campaign, err := campaign_repo.FindByID(tx, request.CampaignID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("campaign not found")
	}
	round := &database.Round{
		RoundID:       GenerateID("r"),
		CreatedByID:   request.CreatedByID,
		CampaignID:    campaign.CampaignID,
		RoundWritable: request.RoundWritable,
	}
	round, err = round_repo.Create(tx, round)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return round, nil
}
func (s *RoundService) ListAllRounds(filter *database.RoundFilter) ([]database.Round, error) {
	round_repo := database.NewRoundRepository()
	conn, close := database.GetDB()
	defer close()
	rounds, err := round_repo.FindAll(conn, filter)
	if err != nil {
		return nil, err
	}
	return rounds, nil
}
func importImagesFromCommons(taskId database.IDType, categories []string) {
	submissions := []database.Submission{}
	successCount := 0
	failedCount := 0
	failedimages := []string{}
	commons_repo := database.NewCommonsRepository()
	round_repo := database.NewRoundRepository()
	task_repo := database.NewTaskRepository()
	conn, close := database.GetDbWithoutDefaultTransaction()
	defer close()
	task, err := task_repo.FindByID(conn, taskId)
	if err != nil {
		log.Printf("Error finding task with ID: %s\n", taskId)
		return
	} else if task == nil {
		return
	}
	round, err := round_repo.FindByID(conn, *task.AssociatedRoundID)
	if err != nil {
		log.Printf("Error finding round with ID: %s\n", *task.AssociatedRoundID)
		return
	}
	task.Status = database.TaskStatusRunning
	conn.Save(task)
	defer func() {
		conn.Save(task)
	}()
	for _, category := range categories {
		images, failedImages := commons_repo.GetImagesFromCommonsCategories(category)
		if images == nil {
			log.Println("No images found in the category: ", category)
			continue
		}
		if failedImages != nil {
			failedCount += len(failedImages)
		}
		failedimages = append(failedimages, failedImages...)
		successCount += len(images)
		// submission_repo := database.NewSubmissionRepository()
		participants := map[string]database.IDType{}
		for _, image := range images {
			participants[image.UploaderUsername] = GenerateID("user")
		}
		part_repo := database.NewUserRepository()
		perCategoryTx := conn.Begin()
		username2IdMap, err := part_repo.EnsureExists(perCategoryTx, participants)
		fmt.Println("Username to ID map: ", username2IdMap)
		if err != nil {
			log.Println("Error ensuring users exist: ", err)
			perCategoryTx.Rollback()
			task.Status = database.TaskStatusFailed
			return
		}
		submissionCount := 0
		for _, image := range images {
			uploaderId := username2IdMap[image.UploaderUsername]
			submission := database.Submission{
				SubmissionID:   GenerateID("s"),
				Name:           image.Name,
				CampaignID:     *task.AssociatedCampaignID,
				URL:            image.URL,
				Author:         image.UploaderUsername,
				SubmittedByID:  uploaderId,
				ParticipantID:  uploaderId,
				SubmittedAt:    image.SubmittedAt,
				CurrentRoundID: round.RoundID,
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
			submissionCount++
		}
		task.SuccessCount = successCount
		task.FailedCount = failedCount
		*task.FailedIds = datatypes.JSONSlice[string](failedimages)
		res := perCategoryTx.Save(task)
		if res.Error != nil {
			log.Println("Error saving task: ", res.Error)
			task.Status = database.TaskStatusFailed
			perCategoryTx.Rollback()
			return
		}
		perCategoryTx.Commit()
	}
	if len(submissions) == 0 {
		task.Status = database.TaskStatusFailed
		return
	}
	tx := conn.Begin()
	res := tx.Create(submissions)
	if res.Error != nil {
		task.Status = database.TaskStatusFailed
		log.Println("Error saving submissions: ", res.Error)
		tx.Rollback()
		return
	}
	round.Status = database.RoundStatusCompleted
	task.Status = database.TaskStatusSuccess
	round.TotalSubmissions += successCount
	tx.Save(task)
	tx.Save(round)
	tx.Commit()
}
func (b *RoundService) ImportFromCommons(roundId database.IDType, categories []string) (*database.Task, error) {
	round_repo := database.NewRoundRepository()
	task_repo := database.NewTaskRepository()
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
	taskReq := &database.Task{
		TaskID:               GenerateID("t"),
		Type:                 database.TaskTypeImportFromCommons,
		Status:               database.TaskStatusPending,
		AssociatedRoundID:    &roundId,
		AssociatedUserID:     &round.CreatedByID,
		CreatedByID:          round.CreatedByID,
		AssociatedCampaignID: &round.CampaignID,
		SuccessCount:         0,
		FailedCount:          0,
		FailedIds:            &datatypes.JSONSlice[string]{},
		RemainingCount:       0,
	}
	task, err := task_repo.Create(tx, taskReq)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	fmt.Println("Task created with ID: ", task.TaskID)
	go importImagesFromCommons(task.TaskID, categories)
	return task, nil
}
func (b *RoundService) GetById(roundId database.IDType) (*database.Round, error) {
	round_repo := database.NewRoundRepository()
	conn, close := database.GetDB()
	defer close()
	return round_repo.FindByID(conn, roundId)
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
