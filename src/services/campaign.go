package services

import (
	"fmt"
	"nokib/campwiz/database"
)

type JuryUserName string
type CampaignService struct{}
type CampaignCreateRequest struct {
	database.CampaignWithWriteableFields
	CreatedByID string         `json:"-"`
	Jury        []JuryUserName `json:"jury"`
	database.CampaignRoundRestrictions
}
type CampaignUpdateRequest struct {
	CampaignCreateRequest
}

func NewCampaignService() *CampaignService {
	return &CampaignService{}
}

func (service *CampaignService) CreateCampaign(campaignRequest *CampaignCreateRequest) (*database.Campaign, error) {
	// Create a new campaign
	campaign := &database.Campaign{
		CampaignID: GenerateID("c"),
		CampaignWithWriteableFields: database.CampaignWithWriteableFields{
			Name:        campaignRequest.Name,
			Description: campaignRequest.Description,
			StartDate:   campaignRequest.StartDate,
			EndDate:     campaignRequest.EndDate,
			Language:    campaignRequest.Language,
			Rules:       campaignRequest.Rules,
			Image:       campaignRequest.Image,
		},
		CreatedByID: campaignRequest.CreatedByID,
	}
	campaign_repo := database.NewCampaignRepository()
	conn, close := database.GetDB()
	defer close()
	tx := conn.Begin()
	err := campaign_repo.Create(tx, campaign)
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
func (service *CampaignService) UpdateCampaign(ID string, campaignRequest *CampaignUpdateRequest) (*database.Campaign, error) {
	conn, close := database.GetDB()
	defer close()
	campaign_repo := database.NewCampaignRepository()
	campaign, err := campaign_repo.FindByID(conn, ID)
	if err != nil {
		return nil, err
	}
	campaign.Name = campaignRequest.Name
	campaign.Description = campaignRequest.Description
	campaign.StartDate = campaignRequest.StartDate
	campaign.EndDate = campaignRequest.EndDate
	// campaign.Language = campaignRequest.Language
	campaign.Rules = campaignRequest.Rules
	campaign.Image = campaignRequest.Image
	err = campaign_repo.Update(conn, campaign)
	if err != nil {
		return nil, err
	}
	return campaign, nil
}
