package database

import (
	"time"

	"gorm.io/gorm"
)

// These are the restrictions that are applied to the articles that are submitted to the campaign
type CampaignRoundCommonRestrictions struct{}

// These are the restrictions that are applied to the audio and video that are submitted to the campaign
type CampaignRoundAudioVideoRestrictions struct {
	MinimumDurationMilliseconds int `json:"minimumDurationMilliseconds" gorm:"default:0"`
}

// These are the restrictions that are applied to the images that are submitted to the campaign
type CampaignRoundImageRestrictions struct {
	MinimumHeight     int `json:"minimumHeight" gorm:"default:0"`
	MinimumWidth      int `json:"minimumWidth" gorm:"default:0"`
	MinimumResolution int `json:"minimumResolution" gorm:"default:0"`
}

// these are the restrictions that are applied to
type CampaignRoundRestrictions struct {
	CampaignRoundCommonRestrictions
	CampaignRoundAudioVideoRestrictions
	CampaignRoundImageRestrictions
}
type CampaignRound struct {
	ID               string         `json:"id" gorm:"primaryKey"`
	CampaignID       string         `json:"campaignId" gorm:"index"`
	Name             string         `json:"name"`
	Description      string         `json:"description" gorm:"type:text"`
	StartDate        time.Time      `json:"startDate" gorm:"type:datetime"`
	EndDate          time.Time      `json:"endDate" gorm:"type:datetime"`
	IsOpen           bool           `json:"isOpen" gorm:"default:true"`
	IsPublic         bool           `json:"isPublic" gorm:"default:false"`
	DependsOnRoundID *string        `json:"dependsOnRoundId" gorm:"default:null"`
	CreatedByID      string         `json:"createdById"`
	DependsOnRound   *CampaignRound `json:"-" gorm:"foreignKey:DependsOnRoundID"`
	Campaign         *Campaign      `json:"-" gorm:"foreignKey:CampaignID"`
	MediaCampaignRestrictions
}
type CampaignRoundRepository struct{}

func NewCampaignRoundRepository() *CampaignRoundRepository {
	return &CampaignRoundRepository{}
}
func (r *CampaignRoundRepository) Create(conn *gorm.DB, campaignRound *CampaignRound) error {
	result := conn.Create(campaignRound)
	return result.Error
}
