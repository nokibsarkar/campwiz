package consts

import (
	"database/sql/driver"
)

// Permission is a type to represent a permission
type Permission uint64

// PermissionGroup is a type to represent a group of permissions
type PermissionGroup uint64

const (
	// Permissions
	PermissionCreateCampaign Permission = 1 << iota
	PermissionUpdateCampaign
	PermissionChangeCampaignStatus
	PermissionDeleteCampaign
	PermissionSeeUser
	PermissionLogin
)
const (
	// PermissionGroups
	PermissionGroupUSER  = PermissionGroup(PermissionLogin)
	PermissionGroupADMIN = PermissionGroupUSER | PermissionGroup(PermissionCreateCampaign) | PermissionGroup(PermissionUpdateCampaign) | PermissionGroup(PermissionChangeCampaignStatus) | PermissionGroup(PermissionDeleteCampaign) | PermissionGroup(PermissionSeeUser)
)

// HasPermission checks if a permission has a permission
func (pg PermissionGroup) HasPermission(p Permission) bool {
	return pg&PermissionGroup(p) == PermissionGroup(p)
}

// AddPermission adds a permission to a permission group
func (pg *PermissionGroup) AddPermission(p Permission) {
	*pg |= PermissionGroup(p)
}

// RemovePermission removes a permission from a permission group
func (pg *PermissionGroup) RemovePermission(p Permission) {
	*pg &= ^PermissionGroup(p)
}
func (p *PermissionGroup) Scan(value interface{}) error {
	*p = PermissionGroup(value.(int64))
	return nil

}

func (p *PermissionGroup) Value() (driver.Value, error) {
	return int64(*p), nil
}
