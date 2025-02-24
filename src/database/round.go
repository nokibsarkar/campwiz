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
type CampaignRoundWritable struct {
	Name             string         `json:"name"`
	Description      string         `json:"description" gorm:"type:text"`
	StartDate        time.Time      `json:"startDate" gorm:"type:datetime"`
	EndDate          time.Time      `json:"endDate" gorm:"type:datetime"`
	IsOpen           bool           `json:"isOpen" gorm:"default:true"`
	IsPublic         bool           `json:"isPublic" gorm:"default:false"`
	DependsOnRoundID *string        `json:"dependsOnRoundId" gorm:"default:null"`
	DependsOnRound   *CampaignRound `json:"-" gorm:"foreignKey:DependsOnRoundID"`
	Campaign         *Campaign      `json:"-" gorm:"foreignKey:CampaignID"`
	Serial           int            `json:"serial" gorm:"default:0"`
	MediaCampaignRestrictions
}
type CampaignRound struct {
	CampaignID  string     `json:"campaignId" gorm:"index"`
	ID          string     `json:"id" gorm:"primaryKey"`
	CreatedAt   *time.Time `json:"createdAt" gorm:"-<-:create"`
	CreatedByID string     `json:"createdById"`
	CampaignRoundWritable
}
type RoundFilter struct {
	CampaignID string `form:"campaignId"`
	Limit      int    `form:"limit"`
}
type CampaignRoundRepository struct{}

func NewCampaignRoundRepository() *CampaignRoundRepository {
	return &CampaignRoundRepository{}
}
func (r *CampaignRoundRepository) Create(conn *gorm.DB, rounds []CampaignRound) error {
	result := conn.Create(rounds)
	return result.Error
}
func (r *CampaignRoundRepository) FindAll(conn *gorm.DB, filter *RoundFilter) ([]CampaignRound, error) {
	var rounds []CampaignRound
	where := &CampaignRound{}
	if filter != nil {
		if filter.CampaignID != "" {
			where.CampaignID = filter.CampaignID
		}
	}
	stmt := conn.Where(where)
	if filter.Limit > 0 {
		stmt = stmt.Limit(filter.Limit)
	}
	result := stmt.Find(&rounds)
	return rounds, result.Error
}
func (r *CampaignRoundRepository) FindByID(conn *gorm.DB, id string) (*CampaignRound, error) {
	round := &CampaignRound{}
	where := &CampaignRound{ID: id}
	result := conn.First(round, where)
	return round, result.Error
}
