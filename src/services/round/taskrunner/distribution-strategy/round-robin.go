package distributionstrategy

import (
	"log"
	"nokib/campwiz/database"

	"gorm.io/gorm"
)

type RoundRobinDistributionStrategy struct{}

func NewRoundRobinDistributionStrategy() *RoundRobinDistributionStrategy {
	return &RoundRobinDistributionStrategy{}
}
func (r *RoundRobinDistributionStrategy) AssignJuries(conn *gorm.DB, round *database.Round, juries []database.Role) error {
	submission_repo := database.NewSubmissionRepository()
	submissions, err := submission_repo.ListAllSubmissions(conn, &database.SubmissionListFilter{})
	if err != nil {
		return err
	}
	submission_ids := make([]database.IDType, len(submissions))
	totalAssignments := 0
	for i, submission := range submissions {
		submission_ids[i] = submission.SubmissionID
		totalAssignments++
	}
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
	for i, jury := range juries {
		remainingQuta := highLimitPerJury - jury.TotalAssigned
		last := min(startIndex+remainingQuta-1, totalAssignments-1)
		log.Printf("Jury %d would have %d remaining assignments", i, last-startIndex)
		juryPosition2Range[i] = append(juryPosition2Range[i], Range{Start: startIndex, End: last})
		startIndex = last + 1
	}
	log.Println("Jury Position to Range: ", juryPosition2Range)
	return nil
}
