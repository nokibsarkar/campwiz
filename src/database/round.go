package database

import (
	"time"
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
	MediaCampaignRestrictions
}
