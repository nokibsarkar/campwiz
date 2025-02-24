package services

import (
	"errors"
	"nokib/campwiz/database"
)

type RoundService struct {
}
type RoundRequest struct {
	CampaignID  string                           `json:"campaignId"`
	CreatedByID string                           `json:"-"`
	Rounds      []database.CampaignRoundWritable `json:"rounds"`
}

func NewRoundService() *RoundService {
	return &RoundService{}
}
func (s *RoundService) CreateRound(request *RoundRequest) ([]database.CampaignRound, error) {
	round_repo := database.NewCampaignRoundRepository()
	campaign_repo := database.NewCampaignRepository()
	conn, close := database.GetDB()
	defer close()
	tx := conn.Begin()
	campaign, err := campaign_repo.FindByID(tx, request.CampaignID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("campaign not found")
	}
	rounds := make([]database.CampaignRound, len(request.Rounds))
	for i, round := range request.Rounds {
		rounds[i] = database.CampaignRound{
			ID:                    GenerateID(),
			CreatedByID:           request.CreatedByID,
			CampaignID:            campaign.CampaignID,
			CampaignRoundWritable: round,
		}
	}
	err = round_repo.Create(tx, rounds)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return rounds, nil
}
func (s *RoundService) ListAllRounds(filter *database.RoundFilter) ([]database.CampaignRound, error) {
	round_repo := database.NewCampaignRoundRepository()
	conn, close := database.GetDB()
	defer close()
	rounds, err := round_repo.FindAll(conn, filter)
	if err != nil {
		return nil, err
	}
	return rounds, nil
}
