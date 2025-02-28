package services

import (
	"fmt"
	"log"
	"nokib/campwiz/database"
	idgenerator "nokib/campwiz/services/idGenerator"
)

// UserName is a type for jury user name
type UserName string
type CampaignService struct{}
type CampaignCreateRequest struct {
	database.CampaignWithWriteableFields
	CreatedByID  database.IDType `json:"-"`
	Coordinators []UserName      `json:"coordinators"`
	Organizers   []UserName      `json:"organizers"`
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
	user_repo := database.NewUserRepository()
	role_repo := database.NewJuryRepository()
	conn, close := database.GetDB()
	defer close()
	tx := conn.Begin()
	err := campaign_repo.Create(tx, campaign)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	username2IDMap := map[string]database.IDType{}
	for _, coordinatorUsername := range campaignRequest.Coordinators {
		username2IDMap[string(coordinatorUsername)] = idgenerator.GenerateID("u")
	}
	for _, organizerUsername := range campaignRequest.Organizers {
		username2IDMap[string(organizerUsername)] = idgenerator.GenerateID("u")
	}
	username2IDMap, err = user_repo.EnsureExists(tx, username2IDMap)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	roles := []database.Role{}
	for _, userName := range campaignRequest.Coordinators {
		userID, ok := username2IDMap[string(userName)]
		if !ok {
			log.Println("User not found: ", userName)
			continue
		}
		if userID == "" {
			continue
		}
		roles = append(roles, database.Role{
			RoleID:     idgenerator.GenerateID("j"),
			UserID:     userID,
			CampaignID: campaign.CampaignID,
			Type:       database.RoleTypeCoordinator,
			RoundID:    nil,
		})
	}
	for _, userName := range campaignRequest.Organizers {
		userID, ok := username2IDMap[string(userName)]
		if !ok {
			log.Println("User not found: ", userName)
			continue
		}
		if userID == "" {
			continue
		}
		roles = append(roles, database.Role{
			RoleID:     idgenerator.GenerateID("j"),
			UserID:     userID,
			CampaignID: campaign.CampaignID,
			Type:       database.RoleTypeOrganizer,
			RoundID:    nil,
		})
	}
	if len(roles) == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("no valid coordinators or organizers found")
	}
	err = role_repo.CreateRoles(tx, roles)
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

// UpdateCampaign updates a campaign
// @summary Update a campaign
// @description Update a campaign
// @tags Campaign
// @param id path string true "The campaign ID"
// @param campaignRequest body CampaignUpdateRequest true "The campaign request"
// @produce json
// @success 200 {object} database.Campaign
// @router /campaign/{id} [post]
func (service *CampaignService) UpdateCampaign(ID database.IDType, campaignRequest *CampaignUpdateRequest) (*database.Campaign, error) {
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
