package services

import "nokib/campwiz/database"

type CampaignService struct{}
type CampaignRequest struct {
	database.Campaign
	Jury []uint `json:"jury"`
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
		Name:        campaignRequest.Name,
		Description: campaignRequest.Description,
		ID:          GenerateID(),
		StartDate:   campaignRequest.StartDate,
		EndDate:     campaignRequest.EndDate,
		Language:    campaignRequest.Language,
		Rules:       campaignRequest.Rules,
		Image:       campaignRequest.Image,
	}
	conn, close := database.GetDB()
	defer close()
	result := conn.Create(campaign)
	if result.Error != nil {
		return nil, result.Error
	}
	return campaign, nil
}
func (service *CampaignService) GetAllCampaigns() []database.Campaign {
	conn, close := database.GetDB()
	defer close()
	var campaigns []database.Campaign
	conn.Find(&campaigns)
	return campaigns
}
