package services

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"nokib/campwiz/database"
	idgenerator "nokib/campwiz/services/idGenerator"
	importservice "nokib/campwiz/services/round/taskrunner"
	importsources "nokib/campwiz/services/round/taskrunner/import-sources"
	"slices"

	"gorm.io/datatypes"
)

type RoundService struct {
}
type RoundRequest struct {
	CampaignID  database.IDType `json:"campaignId"`
	CreatedByID database.IDType `json:"-"`
	database.RoundWritable
}

type DistributionRequest struct {
	AmongJuries []database.IDType `json:"juries"`
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
		RoundID:       idgenerator.GenerateID("r"),
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
		TaskID:               idgenerator.GenerateID("t"),
		Type:                 database.TaskTypeImportFromCommons,
		Status:               database.TaskStatusPending,
		AssociatedRoundID:    &roundId,
		AssociatedUserID:     &round.CreatedByID,
		CreatedByID:          round.CreatedByID,
		AssociatedCampaignID: &round.CampaignID,
		SuccessCount:         0,
		FailedCount:          0,
		FailedIds:            &datatypes.JSONType[map[string]string]{},
		RemainingCount:       0,
	}
	task, err := task_repo.Create(tx, taskReq)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	fmt.Println("Task created with ID: ", task.TaskID)
	commonsCategorySource := importsources.NewCommonsCategoryListSource(categories)
	batch_processor := importservice.NewTaskRunner(task.TaskID, commonsCategorySource)
	go batch_processor.Run()
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
func (r *RoundService) UpdateRoundDetails(roundID database.IDType, req *RoundRequest) (*database.Round, error) {
	round_repo := database.NewRoundRepository()
	conn, close := database.GetDB()
	defer close()
	tx := conn.Begin()
	round, err := round_repo.FindByID(tx, roundID)
	if err != nil {
		tx.Rollback()
		return nil, err
	} else if round == nil {
		tx.Rollback()
		return nil, fmt.Errorf("round not found")
	}
	// round.CampaignID = req.CampaignID
	round.RoundWritable = req.RoundWritable
	round, err = round_repo.Update(tx, round)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return round, nil
}
func (r *RoundService) DistributeEvaluations(roundId database.IDType, distributionReq *DistributionRequest) (*database.Task, error) {
	return nil, nil
}

/*
NewImages := []database.ImageResult{}
NewImageCount := len(NewImages)
Jury[i] is already assigned count of the jury
TotalEvaluation := NewImageCount * evaluationCountRequired + sum(Jury[i])
TotalJury := len(Jury)
Average := TotalEvaluation / TotalJury
Now the goal is to distribute in such way so that each jury would be assigned Average +/- toleranceCount
If the toleranceCount is 0, then it will be set to 1
Difference[i] = Average - Jury[i] // if Difference[i] is positive, then the jury has to be assigned more images, if negative or zero then discard the jury because he already has more images
Now Difference[i] images will be assigned to the jury[i]
Adjusted Jury Count = len(Difference) after removing the juries with Difference[i] <= 0
Now the goal is to distribute the remaining images among the Adjusted Jury Count
Give each image a serial number from 0 to NewImageCount
Now sort the juries based on the Difference[i]
Now start assigning the images to the juries in a round-robin fashion
ImageCount 16
AlreadyAssigned 4
JuryCount 5
EvaluationCountRequired 1
ToleranceCount 1
TotalEvaluation 16 * 1 + 4 = 20
TotalJury 5
Average 19 / 5 = 4
------------------------
Jury 0: 13	14	15	16
Jury 1:	9	10	11	12
Jury 2: * 	6	7	8
Jury 3: * 	3	4	5
Jury 4: * 	* 	1 	2
------------------------
Now for the next iteration,
adjusted set = images
jurylist = previous jurylist - 4 (because jury 4 already have been assigned 1 and 2 once)
Images: 3 4 5 6 7 8 9 10 11 12 13 14 15 16
Total Evaluation = 16 * 1 + 4 = 20
Total Jury = 4
Average = 20 / 4 = 5
Difference = 0 0 0 0
---------------------
Jury 0: 13	14	15	16 8
Jury 1:	9	10	11	12
Jury 2: * 	6	7	8	3	4	5	6	7
Jury 3: * 	3	4	5
// Jury 4: * 	* 	1 	2
*/

func Distribute() {

}
