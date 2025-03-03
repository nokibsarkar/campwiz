package services

import (
	"fmt"
	"nokib/campwiz/database"
	idgenerator "nokib/campwiz/services/idGenerator"
)

// UserName is a type for jury user name

type CampaignService struct{}
type CampaignCreateRequest struct {
	database.CampaignWithWriteableFields
	CreatedByID  database.IDType     `json:"-"`
	Coordinators []database.UserName `json:"coordinators"`
	Organizers   []database.UserName `json:"organizers"`
	database.RoundRestrictions
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
		CampaignID: idgenerator.GenerateID("c"),
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
	// user_repo := database.NewUserRepository()
	role_service := NewRoleService()
	conn, close := database.GetDB()
	defer close()
	tx := conn.Begin()
	err := campaign_repo.Create(tx, campaign)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = role_service.FetchChangeRoles(tx, database.RoleTypeCoordinator, campaign.CampaignID, "", campaignRequest.Coordinators)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = role_service.FetchChangeRoles(tx, database.RoleTypeOrganizer, campaign.CampaignID, "", campaignRequest.Organizers)
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
func (service *CampaignService) GetCampaignByID(id database.IDType) (*database.Campaign, error) {
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

func (service *CampaignService) UpdateCampaign(ID database.IDType, campaignRequest *CampaignUpdateRequest) (*database.Campaign, error) {
	conn, close := database.GetDB()
	defer close()
	campaign_repo := database.NewCampaignRepository()
	campaign, err := campaign_repo.FindByID(conn, ID)
	if err != nil {
		return nil, err
	}
	roleService := NewRoleService()
	campaign.Name = campaignRequest.Name
	campaign.Description = campaignRequest.Description
	campaign.StartDate = campaignRequest.StartDate
	campaign.EndDate = campaignRequest.EndDate
	// campaign.Language = campaignRequest.Language
	campaign.Rules = campaignRequest.Rules
	campaign.Image = campaignRequest.Image
	tx := conn.Begin()
	err = campaign_repo.Update(tx, campaign)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = roleService.FetchChangeRoles(tx, database.RoleTypeOrganizer, campaign.CampaignID, "", campaignRequest.Organizers)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = roleService.FetchChangeRoles(tx, database.RoleTypeCoordinator, campaign.CampaignID, "", campaignRequest.Coordinators)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return campaign, nil
}
