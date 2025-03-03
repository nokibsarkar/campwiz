package distributionstrategy

import (
	"bytes"
	"fmt"
	"nokib/campwiz/database"

	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RoundRobinDistributionStrategySimulator struct {
	TaskId database.IDType
}

func NewRoundRobinDistributionStrategySimulator(taskId database.IDType) *RoundRobinDistributionStrategySimulator {
	return &RoundRobinDistributionStrategySimulator{
		TaskId: taskId,
	}
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
func (r *RoundRobinDistributionStrategySimulator) AssignJuries(tx *gorm.DB, round *database.Round, juries []database.Role) (int, error) {
	submission_repo := database.NewSubmissionRepository()
	numberOfImages, err := submission_repo.GetSubmissionCount(tx, &database.SubmissionListFilter{
		RoundID:    round.RoundID,
		CampaignID: round.CampaignID,
	})
	result := &DistributionResult{}
	if err != nil {
		return 0, err
	}
	numberOfJury := len(juries)
	alreadyWorkload := make([]WorkLoadType, numberOfJury)
	distinctJuryPerImage := 1
	for index, jury := range juries {
		alreadyWorkload[index] = WorkLoadType(jury.TotalAssigned)
	}
	totalWorkload, averageWorkload, distribution, err := simulateDistributeImagesBalanced(int(numberOfImages), numberOfJury, distinctJuryPerImage, alreadyWorkload)
	if err != nil {
		fmt.Println("Error: ", err)
		return 0, err
	}
	fmt.Println("Total Workload: ", totalWorkload)
	fmt.Println("Average Workload: ", averageWorkload)
	fmt.Println("Workload Distribution: ", distribution)
	result.TotalWorkLoad = totalWorkload
	result.AvergaeWorkLoad = averageWorkload
	result.TotalWorkloadDistribution = make(map[database.IDType]WorkLoadType)
	for i, jury := range juries {
		result.TotalWorkloadDistribution[jury.RoleID] = distribution[i]
	}
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(result)
	if err != nil {
		return 0, err
	}
	res := buf.Bytes()
	tx.Updates(&database.Task{
		TaskID: r.TaskId,
		Data:   (*datatypes.JSON)(&res),
	})
	return int(result.TotalWorkLoad), nil
}
