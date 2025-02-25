package database

import (
	"time"

	"gorm.io/gorm"
)

// These are the restrictions that are applied to the articles that are submitted to the campaign
type CampaignRoundCommonRestrictions struct {
	AllowJuryToParticipate bool `json:"allowJuryToParticipate"`
	AllowMultipleJudgement bool `json:"allowMultipleJudgement"`
}

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
type CampaignRoundArticleRestrictions struct {
	MaximumSubmissionOfSameArticle int    `json:"maximumSubmissionOfSameArticle"`
	AllowExpansions                bool   `json:"allowExpansions"`
	AllowCreations                 bool   `json:"allowCreations"`
	MinimumTotalBytes              int    `json:"minimumTotalBytes"`
	MinimumTotalWords              int    `json:"minimumTotalWords"`
	MinimumAddedBytes              int    `json:"minimumAddedBytes"`
	MinimumAddedWords              int    `json:"minimumAddedWords"`
	SecretBallot                   bool   `json:"secretBallot"`
	Blacklist                      string `json:"blacklist"`
}
type CampaignRoundMediaRestrictions struct {
	CampaignRoundImageRestrictions
	CampaignRoundAudioVideoRestrictions
}

// these are the restrictions that are applied to
type CampaignRoundRestrictions struct {
	CampaignRoundCommonRestrictions
	CampaignRoundMediaRestrictions
	CampaignRoundArticleRestrictions
	AllowedMediaTypes MediaTypeSet `json:"allowedMediaTypes" gore:"type:text;not null;default:'ARTICLE'"`
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
	CampaignRoundRestrictions
}
type CampaignRound struct {
	RoundID          string     `json:"roundId" gorm:"primaryKey"`
	CampaignID       string     `json:"campaignId" gorm:"index"`
	CreatedAt        *time.Time `json:"createdAt" gorm:"-<-:create"`
	CreatedByID      string     `json:"createdById"`
	TotalSubmissions int        `json:"totalSubmissions" gorm:"default:0"`
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
func (r *CampaignRoundRepository) Create(conn *gorm.DB, round *CampaignRound) (*CampaignRound, error) {
	result := conn.Create(round)
	if result.Error != nil {
		return nil, result.Error
	}
	return round, nil
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
	where := &CampaignRound{RoundID: id}
	result := conn.First(round, where)
	return round, result.Error
}
