package services

import "nokib/campwiz/database"

type CampaignService struct{}

func NewCampaignService() *CampaignService {
	return &CampaignService{}
}

func (service *CampaignService) CreateCampaign(campaignRequest *database.CampaignRequest) (*database.Campaign, error) {
	// Create a new campaign
	campaign := &database.Campaign{
		Name:        campaignRequest.Name,
		Description: campaignRequest.Description,
	}
	conn, close := database.GetDB()
	defer close()
	result := conn.Create(campaign)
	if result.Error != nil {
		return nil, result.Error
	}
	return campaign, nil
}
