package consts

import (
	"database/sql/driver"
)

// Permission is a type to represent a permission
type Permission uint64

// PermissionGroup is a type to represent a group of permissions
type PermissionGroup Permission

const (
	// Permission to login
	PermissionLogin Permission = 1 << iota
	// Permission to Create a Campaign
	PermissionCreateCampaign
	// Permission to Update an existing Campaign (but it's not possible to change the status of the campaign)
	PermissionUpdateCampaignDetails
	// Permission to Change the status of a Campaign, but not to update the details of the campaign
	PermissionUpdateCampaignStatus
	// Permission to Delete an existing Campaign
	PermissionDeleteCampaign
	// Permission to Create a Round on an existing Campaign
	PermissionCreateRound
	// Permission to update the details of a Round
	PermissionUpdateRoundDetails
	// Permission to update the status of a Round
	PermissionUpdateRoundStatus
	// Permission to Delete a Round
	PermissionDeleteRound
	// Permission to see the list of users
	PermissionSeeAllUsers
	// Permission to see the details of a user
	PermissionSeeUserDetails
	// Permission to See self Evaluation Result
	PermissionSeeSelfEvaluationResult
	// Permission to See Evaluation Result of others
	PermissionSeeOthersEvaluationResult
	// Permission to Approve or Reject a Campaign
	PermissionApproveRejectCampaign
	// Permission to evaluate a submission
	PermissionEvaluateSubmission
	// Permission to submit a submission
	PermissionSubmitSubmission
)
const (
	// PermissionGroups
	PermissionGroupBanned     = PermissionGroup(0)
	PermissionGroupUSER       = PermissionGroup(PermissionLogin)
	PermissionGroupADMIN      = PermissionGroupUSER | PermissionGroup(PermissionCreateCampaign) | PermissionGroup(PermissionUpdateCampaignDetails) | PermissionGroup(PermissionUpdateCampaignStatus) | PermissionGroup(PermissionDeleteCampaign) | PermissionGroup(PermissionSeeUserDetails)
	PermissionGroupSuperAdmin = PermissionGroupADMIN | PermissionGroup(PermissionCreateRound)
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
