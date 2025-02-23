package services

import (
	"fmt"
	"nokib/campwiz/database"
)

type CampaignService struct{}
type CampaignRequest struct {
	database.CampaignWithWriteableFields
	CreatedBy string `json:"createdBy"`
	Jury      []uint `json:"jury"`
}

func NewCampaignService() *CampaignService {
	return &CampaignService{}
}

// CreateCampaign creates a new campaign
// @summary Create a new campaign
// @description Create a new campaign
// @tags Campaign
// @param campaignRequest body CampaignRequest true "The campaign request"
// @produce json
// @success 200 {object} database.Campaign
// @router /campaign/ [post]
func (service *CampaignService) CreateCampaign(campaignRequest *CampaignRequest) (*database.Campaign, error) {
	// Create a new campaign
	campaign := &database.Campaign{
		CampaignWithWriteableFields: database.CampaignWithWriteableFields{
			Name:        campaignRequest.Name,
			Description: campaignRequest.Description,
			ID:          GenerateID(),
			StartDate:   campaignRequest.StartDate,
			EndDate:     campaignRequest.EndDate,
			Language:    campaignRequest.Language,
			Rules:       campaignRequest.Rules,
			Image:       campaignRequest.Image,
		},
		CreatedBy: campaignRequest.CreatedBy,
	}
	campaign_repo := database.NewCampaignRepository()
	round_repo := database.NewCampaignRoundRepository()
	conn, close := database.GetDB()
	defer close()
	tx := conn.Begin()
	err := campaign_repo.Create(tx, campaign)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	round := []database.CampaignRound{
		{
			CreatedByID: campaign.CreatedBy,
			CampaignID:  campaign.ID,
			ID:          GenerateID(),
			CampaignRoundWritable: database.CampaignRoundWritable{
				Name:             "Round 0",
				Description:      "The system round of the campaign",
				StartDate:        campaign.StartDate,
				EndDate:          campaign.EndDate,
				IsOpen:           true,
				IsPublic:         false,
				Serial:           0,
				DependsOnRoundID: nil,
				MediaCampaignRestrictions: database.MediaCampaignRestrictions{
					ImageCampaignRestrictions: database.ImageCampaignRestrictions{
						MaximumSubmissionOfSameImage: 1,
						MinimumTotalImageSize:        1024,
					},
				},
			},
		},
	}
	err = round_repo.Create(tx, round)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return campaign, nil
}
func (service *CampaignService) GetAllCampaigns(query *database.CampaignFilter) []database.Campaign {
	conn, close := database.GetDB()
	defer close()
	campaign_repo := database.NewCampaignRepository()

	campaigns, err := campaign_repo.ListAllCampaigns(conn, query)
	if err != nil {
		fmt.Println("Error: ", err)
		return []database.Campaign{}
	}
	return campaigns
}
