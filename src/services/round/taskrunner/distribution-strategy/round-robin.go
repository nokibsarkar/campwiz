package distributionstrategy

import (
	"log"
	"nokib/campwiz/database"
	idgenerator "nokib/campwiz/services/idGenerator"

	"gorm.io/gorm"
)

type RoundRobinDistributionStrategy struct {
	TaskId database.IDType
}

func NewRoundRobinDistributionStrategy(taskId database.IDType) *RoundRobinDistributionStrategy {
	return &RoundRobinDistributionStrategy{
		TaskId: taskId,
	}
}
func (strategy *RoundRobinDistributionStrategy) AssignJuries(conn *gorm.DB, round *database.Round, juries []database.Role) error {
	submission_repo := database.NewSubmissionRepository()
	submissions, err := submission_repo.ListAllSubmissions(conn, &database.SubmissionListFilter{
		RoundID:    round.RoundID,
		CampaignID: round.CampaignID,
	})
	if err != nil {
		return err
	}
	totalAssignments := len(submissions)
	grandTotalAssigned := totalAssignments
	juryCount := 0
	for _, jury := range juries {
		// Assign the juries to the submissions
		grandTotalAssigned += jury.TotalAssigned
		juryCount++
	}
	highLimitPerJury := grandTotalAssigned / juryCount
	// if grandTotalAssigned%juryCount > 0 {
	// 	highLimitPerJury++
	// }
	log.Printf("Each Jury would have total %d asignments (including already assigned)", highLimitPerJury)
	juryPosition2Range := make([][]Range, juryCount)
	startIndex := 0
	const distributionRound = 2
	for distributionIteration := range distributionRound {
		log.Printf("Distribution Iteration %d", distributionIteration)
		for i, jury := range juries {
			remainingQuta := highLimitPerJury - jury.TotalAssigned
			last := min(startIndex+remainingQuta-1, totalAssignments-1)
			log.Printf("Jury %d would have %d remaining assignments", i, last-startIndex)
			juryPosition2Range[i] = append(juryPosition2Range[i], Range{Start: startIndex, End: last})
			startIndex = last + 1
		}
	}
	evaluations := []database.Evaluation{}
	for i, jury := range juries {
		for _, rn := range juryPosition2Range[i] {
			for j := rn.Start; j <= rn.End; j++ {
				submission := submissions[j]
				evaluation := database.Evaluation{
					SubmissionID:       submission.SubmissionID,
					EvaluationID:       idgenerator.GenerateID("e"),
					JudgeID:            jury.RoleID,
					ParticipantID:      submission.ParticipantID,
					DistributionTaskID: strategy.TaskId,
				}
				evaluations = append(evaluations, evaluation)
			}
		}
	}
	res := conn.Create(&evaluations)
	if res.Error != nil {
		return res.Error
	}
	// log.Println("Evaluations: ", evaluations)
	log.Println("Jury Position to Range: ", juryPosition2Range)
	return nil
}
