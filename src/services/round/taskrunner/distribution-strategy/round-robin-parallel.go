package distributionstrategy

import (
	"fmt"
	"log"
	"nokib/campwiz/consts"
	"nokib/campwiz/database"
	idgenerator "nokib/campwiz/services/idGenerator"
	"sync"

	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/util/sets"
)

type RoundRobinDistributionStrategyV2 struct {
	TaskId database.IDType
}
type DistributionResultV2 struct {
	TotalWorkLoad             WorkLoadTypeV2                   `json:"totalWorkLoad"`
	AvergaeWorkLoad           WorkLoadTypeV2                   `json:"averageWorkLoad"`
	TotalWorkloadDistribution map[database.IDType]WorkLoadType `json:"totalworkloadDistribution"`
}
type WorkLoadTypeV2 int64

// Juror represents a juror with workload and ID
type JurorV2 struct {
	ID       int
	Workload WorkLoadTypeV2
}
type TaskDistributionResultV2 struct {
	TotalWorkLoad        WorkLoadTypeV2 `json:"totalWorkLoad"`
	AvergaeWorkLoad      WorkLoadTypeV2 `json:"averageWorkLoad"`
	WorkloadDistribution []WorkLoadType `json:"workloadDistribution"`
}

// MinimumWorkloadHeap is a priority queue of Jurors
type MinimumWorkloadHeapV2 []JurorV2

func (h MinimumWorkloadHeapV2) Len() int           { return len(h) }
func (h MinimumWorkloadHeapV2) Less(i, j int) bool { return h[i].Workload < h[j].Workload }
func (h MinimumWorkloadHeapV2) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinimumWorkloadHeapV2) Push(x any) {
	*h = append(*h, x.(JurorV2))
}
func (h *MinimumWorkloadHeapV2) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
func NewRoundRobinDistributionStrategyV2(taskId database.IDType) *RoundRobinDistributionStrategyV2 {
	return &RoundRobinDistributionStrategyV2{
		TaskId: taskId,
	}
}

// distributeImagesBalanced distributes images among jurors while balancing workload
func distributeImagesBalancedConcurrent(numberOfImages, numberOfJury, distinctJuryPerImage int, alreadyWorkload []WorkLoadTypeV2) ([]sets.Set[int], error) {
	if numberOfJury < distinctJuryPerImage {
		return nil, fmt.Errorf("number of jurors must be greater than or equal to distinct jurors per image")
	}
	log.Println("Concurrently: Number of images: ", numberOfImages)
	log.Println("Number of jurors: ", numberOfJury)
	log.Println("Distinct jurors per image: ", distinctJuryPerImage)
	log.Println("Already workload: ", alreadyWorkload)
	// Min-heap to track least-loaded jurors
	hp := NewNokibPriorityQueue(&MinimumWorkloadHeapV2{})
	// Initialize the heap with jurors and their workloads
	for i := range numberOfJury {
		hp.Push(JurorV2{i, alreadyWorkload[i]})
	}
	assignments := make([]sets.Set[int], numberOfImages)
	batchCount := consts.Config.Distribution.MaximumBatchCount
	batchSize := numberOfImages / batchCount
	if batchSize < consts.Config.Distribution.MinimumBatchSize {
		batchSize = consts.Config.Distribution.MinimumBatchSize
		batchCount = numberOfImages / batchSize
	}
	if numberOfImages%batchSize != 0 {
		batchCount++
	}
	log.Printf("Batch count: %d\n", batchCount)
	var wg sync.WaitGroup
	// Mutual exclusion lock to protect the heap
	// This protects the heap for
	var perImageMutex sync.RWMutex
	wg.Add(batchCount)
	for currentBatch := range batchCount {
		start := currentBatch * batchSize
		end := min((currentBatch+1)*batchSize, numberOfImages)
		log.Println("Processing batch of images from ", start, " to ", end)
		go func() {
			defer wg.Done()
			for i := start; i < end; i++ {
				perImageMutex.Lock()
				// log.Println("Processing image: ", i)
				selectedJurySet := sets.Set[int]{} // To keep track of distinct jurors
				// Pick K distinct least-loaded jurors
				for len(selectedJurySet) < distinctJuryPerImage {
					juror := hp.Pop().(JurorV2)
					// log.Println("Juror: ", juror)
					selectedJurySet.Insert(juror.ID)
				}
				assignments[i] = selectedJurySet
				// Store the assigned jurors for this image
				for juror := range selectedJurySet {
					alreadyWorkload[juror]++
					newJuror := JurorV2{juror, alreadyWorkload[juror]}
					hp.Push(newJuror)
				}
				perImageMutex.Unlock()
			}
		}()
	}
	wg.Wait()
	hp.StopWatchHeapOps()
	return assignments, nil
}

func (strategy *RoundRobinDistributionStrategyV2) AssignJuries(tx *gorm.DB, round *database.Round, juries []database.Role) (evaluationCount int, err error) {
	submission_repo := database.NewSubmissionRepository()
	submissions, err := submission_repo.ListAllSubmissions(tx, &database.SubmissionListFilter{
		RoundID:    round.RoundID,
		CampaignID: round.CampaignID,
	})
	if err != nil {
		return 0, err
	}
	numberOfImages := len(submissions)
	numberOfJury := len(juries)
	alreadyWorkload := make([]WorkLoadTypeV2, numberOfJury)
	distinctJuryPerImage := 1
	for index, jury := range juries {
		alreadyWorkload[index] = WorkLoadTypeV2(jury.TotalAssigned)
	}
	assignments, err := distributeImagesBalancedConcurrent(numberOfImages, numberOfJury, distinctJuryPerImage, alreadyWorkload)
	if err != nil {
		return 0, err
	}
	evaluations := []database.Evaluation{}
	updatedJury := make([]database.Role, numberOfJury)
	for i, assignment := range assignments {
		for juryIndex := range assignment {
			jury := juries[juryIndex]
			submission := submissions[i]
			evaluation := database.Evaluation{
				SubmissionID:       submission.SubmissionID,
				EvaluationID:       idgenerator.GenerateID("e"),
				JudgeID:            jury.RoleID,
				ParticipantID:      submission.ParticipantID,
				DistributionTaskID: strategy.TaskId,
				Type:               round.Type,
				Serial:             uint(i),
			}
			juryRole := updatedJury[juryIndex]
			juryRole.RoleID = jury.RoleID
			juryRole.TotalAssigned++
			updatedJury[juryIndex] = juryRole
			evaluations = append(evaluations, evaluation)
		}
	}
	for _, jury := range updatedJury {
		st := tx.Updates(&jury)
		if st.Error != nil {
			return 0, st.Error
		}
	}
	res := tx.Create(&evaluations)
	return len(evaluations), res.Error
}
