package services

import (
	"nokib/campwiz/database"
)

type SubmissionService struct{}

func NewSubmissionService() *SubmissionService {
	return &SubmissionService{}
}

func (s *SubmissionService) ListAllSubmissions(filter *database.SubmissionListFilter) ([]database.Submission, error) {
	conn, close := database.GetDB()
	defer close()
	submission_repo := database.NewSubmissionRepository()
	submissions, err := submission_repo.ListAllSubmissions(conn, filter)
	if err != nil {
		return nil, err
	}
	return submissions, nil
}
