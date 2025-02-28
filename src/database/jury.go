package database

import "gorm.io/gorm"

type RoleType string

const (
	RoleTypeJury        RoleType = "jury"
	RoleTypeAdmin       RoleType = "admin"
	RoleTypeParticipant RoleType = "participant"
	RoleTypeCoordinator RoleType = "coordinator"
	RoleTypeOrganizer   RoleType = "organizer"
)

type Role struct {
	RoleID         IDType    `json:"roleId" gorm:"primaryKey"`
	Type           RoleType  `json:"type" gorm:"not null"`
	UserID         IDType    `json:"userId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CampaignID     IDType    `json:"campaignId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RoundID        *IDType   `json:"roundId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsAllowed      bool      `json:"isAllowed" gorm:"default:true"`
	TotalAssigned  int       `json:"totalAssigned"`
	TotalEvaluated int       `json:"totalEvaluated"`
	TotalScore     int       `json:"totalScore"`
	Campaign       *Campaign `json:"-"`
	User           *User     `json:"-"`
	Round          *Round    `json:"-"`
}
type RoleFilter struct {
	CommonFilter
	CampaignID IDType    `form:"campaignId"`
	RoundID    IDType    `form:"roundId"`
	Type       *RoleType `form:"type"`
}
type RoleRepository struct{}

func NewJuryRepository() *RoleRepository {
	return &RoleRepository{}
}
func (r *RoleRepository) CreateRole(tx *gorm.DB, jury *Role) error {
	result := tx.Create(jury)
	return result.Error
}
func (r *RoleRepository) CreateRoles(tx *gorm.DB, juries []Role) error {
	result := tx.Create(juries)
	return result.Error
}
func (r *RoleRepository) FindRoleByID(tx *gorm.DB, juryID IDType) (*Role, error) {
	jury := &Role{}
	result := tx.First(jury, &Role{RoleID: juryID})
	return jury, result.Error
}
func (r *RoleRepository) ListAllRoles(tx *gorm.DB, filter *RoleFilter) ([]Role, error) {
	var juries []Role
	stmt := tx
	if filter != nil {
		if filter.CampaignID != "" {
			stmt = stmt.Where(&Role{CampaignID: filter.CampaignID})
		}
		if filter.RoundID != "" {
			stmt = stmt.Where(&Role{RoundID: &filter.RoundID})
		}
		if filter.Type != nil {
			stmt = stmt.Where(&Role{Type: *filter.Type})
		}

	}
	result := stmt.Find(&juries)
	return juries, result.Error
}
