package database

import (
	"time"
)

type CampaignRoundRestrictions struct {
}
type CampaignRound struct {
	ID               uint64         `json:"id" gorm:"primaryKey"`
	CampaignID       uint64         `json:"campaignId"`
	Name             string         `json:"name"`
	Description      string         `json:"description" gorm:"type:text"`
	StartDate        time.Time      `json:"startDate" gorm:"type:datetime"`
	EndDate          time.Time      `json:"endDate" gorm:"type:datetime"`
	IsOpen           bool           `json:"isOpen" gorm:"default:true"`
	IsPublic         bool           `json:"isPublic" gorm:"default:false"`
	DependsOnRoundID uint64         `json:"dependsOnRoundId" gorm:"default:null"`
	CreatedByID      uint64         `json:"createdById"`
	DependsOnRound   *CampaignRound `json:"-"`
	Campaign         *Campaign      `json:"-"`
	// CreatedBy        *User          `json:"-" gorm:"foreignKey:CreatedByID"`
}
