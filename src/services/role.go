package services

import (
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
			continue
		}
		_, ok := userId2NameMap[existingRole.UserID]
		if !ok {
			// not exist in updated roles
			// remove the role by adding isAllowed = false and permission would be null
			removedRole := &database.Role{
				RoleID: existingRole.RoleID,
			}
			removedRoles = append(removedRoles, removedRole.RoleID)
		} else {
			// Remain unchanged
		}
	}
	for userId := range userId2NameMap {
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
				// already exisiting but was banned
				newRole.RoleID = role.RoleID
			}
			addedRoles = append(addedRoles, newRole)
		} else {
			//remain unchanged
		}
	}
	return addedRoles, removedRoles, nil
}
