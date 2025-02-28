package services

import "nokib/campwiz/database"

type EvaluationService struct{}

func NewEvaluationService() *EvaluationService {
	return &EvaluationService{}
}
func (e *EvaluationService) Evaluate() {
}
func (e *EvaluationService) GetEvaluationById() {
}

func (e *EvaluationService) ListEvaluations(filter *database.EvaluationFilter) ([]database.Evaluation, error) {
	ev_repo := database.NewEvaluationRepository()
	conn, close := database.GetDB()
	defer close()
	return ev_repo.ListAllEvaluations(conn, filter)
}
