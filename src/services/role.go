package services

import (
	"log"
	"nokib/campwiz/database"
	idgenerator "nokib/campwiz/services/idGenerator"

	"gorm.io/gorm"
)

type RoleService struct{}

func NewRoleService() *RoleService {
	return &RoleService{}
}
func (r *RoleService) CalculateRoleDifference(tx *gorm.DB, Type database.RoleType, filter *database.RoleFilter, updatedRoleUsernames []database.UserName) (addedRoles []database.Role, removedRoles []database.IDType, err error) {
	role_repo := database.NewRoleRepository()
	user_repo := database.NewUserRepository()
	existingRoles, err := role_repo.ListAllRoles(tx, filter)
	if err != nil {
		return nil, nil, err
	}
	username2IDMap := map[database.UserName]database.IDType{}
	for _, username := range updatedRoleUsernames {
		username2IDMap[username] = idgenerator.GenerateID("u")
	}
	username2IDMap, err = user_repo.EnsureExists(tx, username2IDMap)
	if err != nil {
		return nil, nil, err
	}
	userId2NameMap := map[database.IDType]database.UserName{}
	for username := range username2IDMap {
		userId := username2IDMap[username]
		userId2NameMap[userId] = username
	}

	id2RoleMap := map[database.IDType]database.Role{}
	for _, existingRole := range existingRoles {
		id2RoleMap[existingRole.UserID] = existingRole
		if !existingRole.IsAllowed {
			// Already soft deleted, so either way, pretend it does not exist
			log.Printf("Soft %s", existingRole.UserID)
			continue
		}
		_, ok := userId2NameMap[existingRole.UserID]
		if !ok {
			// not exist in updated roles
			// remove the role by adding isAllowed = false and permission would be null

			log.Printf("%s is being removed as %s", existingRole.RoleID, Type)
			removedRoles = append(removedRoles, existingRole.RoleID)
		} else {
			// Remain unchanged
			log.Printf("%s is unchanged as %s", userId2NameMap[existingRole.UserID], Type)
		}
	}

	for userId := range userId2NameMap {
		log.Printf("%v", userId)
		role, ok := id2RoleMap[userId]
		if !ok || !role.IsAllowed {
			// Newly added
			newRole := database.Role{
				RoleID:     idgenerator.GenerateID("j"),
				Type:       Type,
				UserID:     userId,
				CampaignID: filter.CampaignID,
				IsAllowed:  true,
			}
			if filter.RoundID != "" {
				newRole.RoundID = &filter.RoundID
			}
			if ok {
				log.Printf("%s was banned earlier. Unbanning ", role.RoleID)
				// already exisiting but was banned
				newRole.RoleID = role.RoleID
			}
			log.Printf("adding %s was banned earlier.", newRole.RoleID)
			addedRoles = append(addedRoles, newRole)
		} else {
			//remain unchanged
		}
	}
	log.Println("Add: ", addedRoles)
	log.Println("Remove ", removedRoles)
	return addedRoles, removedRoles, nil
}
func (service *RoleService) FetchChangeRoles(tx *gorm.DB, roleType database.RoleType, campaignId database.IDType, updatedRoleUsernames []database.UserName) error {
	filter := &database.RoleFilter{
		CampaignID: campaignId,
		Type:       &roleType,
	}
	addedRoles, removedRoles, err := service.CalculateRoleDifference(tx, roleType, filter, updatedRoleUsernames)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	if len(addedRoles) > 0 {
		res := tx.Save(addedRoles)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
	}
	if len(removedRoles) > 0 {
		for _, roleID := range removedRoles {
			log.Println("Banning role: ", roleID)
			res := tx.Model(&database.Role{RoleID: roleID}).Update("is_allowed", false)
			if res.Error != nil {
				tx.Rollback()
				return res.Error
			}
		}
	}
	return nil
}
