package distributionstrategy

import (
	"fmt"
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
	fmt.Println(submissions)
	fmt.Println("Submissions: ", submissions)
	return nil
}
