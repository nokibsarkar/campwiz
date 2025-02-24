package services

import (
	"fmt"
	"math/rand/v2"
	"nokib/campwiz/database"
	"slices"
)

type BatchService struct{}
type BatchCreationResult struct {
	database.Batch
	SuccessCount int      `json:"successCount"`
	FailedCount  int      `json:"failedCount"`
	FailedIds    []string `json:"failedIds"`
}
type CreateFromCommons struct {
	Categories []string `json:"categories"` // Categories from which images will be fetched
	CreatedBy  string   `json:"-"`          // User ID who created the batch
	RoundId    string   `json:"roundId"`    // Round ID to which the batch belongs
}

func NewBatchService() *BatchService {
	return &BatchService{}
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
func (b *BatchService) CreateBatchFromCommonsCategory(req *CreateFromCommons) (*BatchCreationResult, error) {
	categories := req.Categories
	commons_repo := database.NewCommonsRepository()
	images, failedImages := commons_repo.GetImagesFromCommonsCategories(categories)
	if images == nil {
		return nil, fmt.Errorf("no images found in the categories")
	}
	successCount := len(images)
	failedCount := 0
	if failedImages != nil {
		failedCount = len(failedImages)
	}
	batch := &database.Batch{
		BatchID:          GenerateID(),
		CreatedByID:      req.CreatedBy,
		TotalSubmissions: successCount,
	}
	fmt.Println("Creating batch with ID: ", batch.BatchID)
	conn, close := database.GetDB()
	defer close()
	tx := conn.Begin()
	round_repo := database.NewCampaignRoundRepository()
	round, err := round_repo.FindByID(tx, req.RoundId)
	if err != nil {
		tx.Rollback()
		return nil, err
	} else if round == nil {
		tx.Rollback()
		return nil, fmt.Errorf("round not found")
	}
	batch_repo := database.NewBatchRepository()
	batch, err = batch_repo.CreateBatch(tx, batch)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	submission_repo := database.NewSubmissionRepository()
	for _, image := range images {
		submission := &database.Submission{
			SubmissionID: GenerateID(),
			Name:         image.Name,
			CampaignID:   round.CampaignID,
		}
		err = submission_repo.CreateSubmission(tx, submission)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &BatchCreationResult{
		Batch:        *batch,
		SuccessCount: successCount,
		FailedCount:  failedCount,
		FailedIds:    failedImages,
	}, nil

}
func (b *BatchService) distributeTaskAmongExistingJuries(images []database.Image) {
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

func (b *BatchService) GetBatchByID(id string) (*database.Batch, error) {
	// Get batch by ID
	repo := database.NewBatchRepository()
	conn, close := database.GetDB()
	defer close()
	batch, err := repo.GetBatchByID(conn, id)
	if err != nil {
		return nil, err
	}
	return batch, nil
}
func (b *BatchService) GetAllBatches(filter *database.BatchFilter) ([]database.Batch, error) {
	// Get all batches
	repo := database.NewBatchRepository()
	conn, close := database.GetDB()
	defer close()
	batches, err := repo.GetBatches(conn, filter)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return batches, nil
}
