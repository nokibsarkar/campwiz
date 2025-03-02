package distributionstrategy

import (
	"container/heap"
	"fmt"
	"log"
	"nokib/campwiz/database"
	idgenerator "nokib/campwiz/services/idGenerator"

	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/util/sets"
)

type RoundRobinDistributionStrategy struct {
	TaskId database.IDType
}

type WorkLoadType int64

// Juror represents a juror with workload and ID
type Juror struct {
	ID       int
	Workload WorkLoadType
}
type TaskDistributionResult struct {
	TotalWorkLoad        WorkLoadType   `json:"totalWorkLoad"`
	AvergaeWorkLoad      WorkLoadType   `json:"averageWorkLoad"`
	WorkloadDistribution []WorkLoadType `json:"workloadDistribution"`
}

// MinimumWorkloadHeap is a priority queue of Jurors
type MinimumWorkloadHeap []Juror

func (h MinimumWorkloadHeap) Len() int           { return len(h) }
func (h MinimumWorkloadHeap) Less(i, j int) bool { return h[i].Workload < h[j].Workload }
func (h MinimumWorkloadHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinimumWorkloadHeap) Push(x interface{}) {
	*h = append(*h, x.(Juror))
}
func (h *MinimumWorkloadHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
func NewRoundRobinDistributionStrategy(taskId database.IDType) *RoundRobinDistributionStrategy {
	return &RoundRobinDistributionStrategy{
		TaskId: taskId,
	}
}

// distributeImagesBalanced distributes images among jurors while balancing workload
func distributeImagesBalanced(numberOfImages, numberOfJury, distinctJuryPerImage int, alreadyWorkload []WorkLoadType) ([]sets.Set[int], error) {
	if numberOfJury < distinctJuryPerImage {
		return nil, fmt.Errorf("number of jurors must be greater than or equal to distinct jurors per image")
	}
	log.Println("Number of images: ", numberOfImages)
	log.Println("Number of jurors: ", numberOfJury)
	log.Println("Distinct jurors per image: ", distinctJuryPerImage)
	log.Println("Already workload: ", alreadyWorkload)
	// Min-heap to track least-loaded jurors
	h := &MinimumWorkloadHeap{}

	// Initialize the heap with jurors and their workloads
	for i := range numberOfJury {
		h.Push(Juror{i, alreadyWorkload[i]})
	}
	heap.Init(h)
	assignments := make([]sets.Set[int], numberOfImages)
	for i := range numberOfImages {
		selectedJurySet := sets.Set[int]{} // To keep track of distinct jurors
		// Pick K distinct least-loaded jurors
		for len(selectedJurySet) < distinctJuryPerImage {
			juror := heap.Pop(h).(Juror)
			selectedJurySet.Insert(juror.ID)
		}
		assignments[i] = selectedJurySet
		// Store the assigned jurors for this image
		for juror := range selectedJurySet {
			alreadyWorkload[juror]++
			newJuror := Juror{juror, alreadyWorkload[juror]}
			heap.Push(h, newJuror)
		}
	}

	return assignments, nil
}
func simulateDistributeImagesBalanced(numberOfImages int, numberOfJury int, distinctJuryPerImage int, alreadyWorkload []WorkLoadType) (totalWorkLoad, avergaeWorkLoad WorkLoadType, workloadDistribution []WorkLoadType, err error) {
	if numberOfJury < distinctJuryPerImage {
		return 0, 0, nil, fmt.Errorf("number of jurors must be greater than or equal to distinct jurors per image")
	}
	var previousWorkload WorkLoadType = 0
	for i := range numberOfJury {
		previousWorkload += alreadyWorkload[i]
	}
	newWorkload := WorkLoadType(numberOfImages) * WorkLoadType(distinctJuryPerImage)
	totalWorkLoad = previousWorkload + newWorkload
	avergaeWorkLoad = totalWorkLoad / WorkLoadType(numberOfJury)
	workloadDistribution = make([]WorkLoadType, numberOfJury)
	extra := totalWorkLoad % WorkLoadType(numberOfJury)
	for i := range numberOfJury {
		workloadDistribution[i] = max(avergaeWorkLoad-alreadyWorkload[i], 0)
		for k := extra; k > 0; k-- {
			if workloadDistribution[i]+k+alreadyWorkload[i] <= avergaeWorkLoad+1 {
				workloadDistribution[i] += k
				extra -= k
				break
			}
		}
	}
	return
}
func (strategy *RoundRobinDistributionStrategy) AssignJuries(tx *gorm.DB, round *database.Round, juries []database.Role) error {
	submission_repo := database.NewSubmissionRepository()
	submissions, err := submission_repo.ListAllSubmissions(tx, &database.SubmissionListFilter{
		RoundID:    round.RoundID,
		CampaignID: round.CampaignID,
	})
	if err != nil {
		return err
	}
	numberOfImages := len(submissions)
	numberOfJury := len(juries)
	alreadyWorkload := make([]WorkLoadType, numberOfJury)
	distinctJuryPerImage := 1
	for index, jury := range juries {
		alreadyWorkload[index] = WorkLoadType(jury.TotalAssigned)
	}
	totalWorkload, averageWorkload, distribution, err := simulateDistributeImagesBalanced(numberOfImages, numberOfJury, distinctJuryPerImage, alreadyWorkload)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	fmt.Println("Total Workload: ", totalWorkload)
	fmt.Println("Average Workload: ", averageWorkload)
	fmt.Println("Workload Distribution: ", distribution)
	assignments, err := distributeImagesBalanced(numberOfImages, numberOfJury, distinctJuryPerImage, alreadyWorkload)
	if err != nil {
		return err
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
			return st.Error
		}
	}
	tx.Create(&evaluations)
	tx.Rollback()
	// for i, jury := range juries {
	// 	for _, rn := range juryPosition2Range[i] {
	// 		for j := rn.Start; j <= rn.End; j++ {
	// 			submission := submissions[j]
	// 			evaluation := database.Evaluation{
	// 				SubmissionID:       submission.SubmissionID,
	// 				EvaluationID:       idgenerator.GenerateID("e"),
	// 				JudgeID:            jury.RoleID,
	// 				ParticipantID:      submission.ParticipantID,
	// 				DistributionTaskID: strategy.TaskId,
	// 				Type:               round.Type,
	// 			}
	// 			evaluations = append(evaluations, evaluation)
	// 		}
	// 	}
	// }
	// res := conn.Create(&evaluations)
	// if res.Error != nil {
	// 	return res.Error
	// }
	return nil
}
