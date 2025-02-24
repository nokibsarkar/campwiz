package services

import (
	"fmt"
	"nokib/campwiz/database"
)

type CampaignService struct{}
type CampaignCreateRequest struct {
	database.CampaignWithWriteableFields
	CreatedBy string `json:"-"`
	Jury      []uint `json:"jury"`
}
type CampaignUpdateRequest struct {
	CampaignCreateRequest
	ID string `json:"campaignId"`
}

func NewCampaignService() *CampaignService {
	return &CampaignService{}
}

func (service *CampaignService) CreateCampaign(campaignRequest *CampaignCreateRequest) (*database.Campaign, error) {
	// Create a new campaign
	campaign := &database.Campaign{
		ID: GenerateID(),
		CampaignWithWriteableFields: database.CampaignWithWriteableFields{
			Name:        campaignRequest.Name,
			Description: campaignRequest.Description,
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
func (service *CampaignService) GetCampaignByID(id string) (*database.Campaign, error) {
	conn, close := database.GetDB()
	defer close()
	campaign_repo := database.NewCampaignRepository()
	campaign, err := campaign_repo.FindByID(conn, id)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return campaign, nil
}

// UpdateCampaign updates a campaign
// @summary Update a campaign
// @description Update a campaign
// @tags Campaign
// @param id path string true "The campaign ID"
// @param campaignRequest body CampaignUpdateRequest true "The campaign request"
// @produce json
// @success 200 {object} database.Campaign
// @router /campaign/{id} [post]
func (service *CampaignService) UpdateCampaign(campaignRequest *CampaignUpdateRequest) (*database.Campaign, error) {
	conn, close := database.GetDB()
	defer close()
	campaign_repo := database.NewCampaignRepository()
	campaign, err := campaign_repo.FindByID(conn, campaignRequest.ID)
	if err != nil {
		return nil, err
	}
	campaign.Name = campaignRequest.Name
	campaign.Description = campaignRequest.Description
	campaign.StartDate = campaignRequest.StartDate
	campaign.EndDate = campaignRequest.EndDate
	campaign.Language = campaignRequest.Language
	campaign.Rules = campaignRequest.Rules
	campaign.Image = campaignRequest.Image

	err = campaign_repo.Update(conn, campaign)
	if err != nil {
		return nil, err
	}
	return campaign, nil
}
