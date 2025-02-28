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

type EvaluationFilter struct {
	Type                 database.EvaluationType `form:"type"`
	AssociatedRoundID    database.IDType         `form:"roundId"`
	AssociatedCampaignID database.IDType         `form:"campaignId"`
	AssociatedUserID     database.IDType         `form:"userId"`
	Status               bool                    `form:"status"`
}

func (e *EvaluationService) ListEvaluations(filter *EvaluationFilter) ([]database.Evaluation, error) {
	return nil, nil
}
