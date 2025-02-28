package database

import (
	"time"

	"gorm.io/gorm"
)

// These are the restrictions that are applied to the articles that are submitted to the campaign
type RoundCommonRestrictions struct {
	AllowJuryToParticipate bool `json:"allowJuryToParticipate"`
	AllowMultipleJudgement bool `json:"allowMultipleJudgement"`
}
type RoundStatus string

const (
	// RoundStatusPending is the status when the round is created but not yet approved by the admin
	RoundStatusPending RoundStatus = "PENDING"
	// RoundStatusImporting is the status when the round is approved and importing images from commons is in progress
	RoundStatusImporting RoundStatus = "IMPORTING"
	// RoundStatusDistributing is the status when the images are imported and being distributed among juries
	RoundStatusDistributing RoundStatus = "DISTRIBUTING"
	// RoundStatusEvaluating is the status when the images are distributed and juries are evaluating the images.
	RoundStatusEvaluating RoundStatus = "EVALUATING"
	// RoundStatusRejected is the status when the round is rejected by the admin
	RoundStatusRejected RoundStatus = "REJECTED"
	// RoundStatusCancelled is the status when the round is cancelled by the admin or the creator
	RoundStatusCancelled RoundStatus = "CANCELLED"
	// RoundStatusPaused is the status when the round is paused by the admin
	RoundStatusPaused RoundStatus = "PAUSED"
	// RoundStatusScheduled is the status when the round is scheduled to start at a later time
	RoundStatusScheduled RoundStatus = "SCHEDULED"
	// RoundStatusActive is the status when the round is active and juries are evaluating the images
	RoundStatusActive RoundStatus = "ACTIVE"
	// RoundStatusCompleted is the status when the round is completed and the results are ready
	RoundStatusCompleted RoundStatus = "COMPLETED"
)

// These are the restrictions that are applied to the audio and video that are submitted to the campaign
type RoundAudioVideoRestrictions struct {
	MinimumDurationMilliseconds int `json:"minimumDurationMilliseconds" gorm:"default:0"`
}

// These are the restrictions that are applied to the images that are submitted to the campaign
type RoundImageRestrictions struct {
	MinimumHeight     int `json:"minimumHeight" gorm:"default:0"`
	MinimumWidth      int `json:"minimumWidth" gorm:"default:0"`
	MinimumResolution int `json:"minimumResolution" gorm:"default:0"`
}
type RoundArticleRestrictions struct {
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
type RoundMediaRestrictions struct {
	RoundImageRestrictions
	RoundAudioVideoRestrictions
}

// these are the restrictions that are applied to
type RoundRestrictions struct {
	RoundCommonRestrictions
	RoundMediaRestrictions
	RoundArticleRestrictions
	AllowedMediaTypes MediaTypeSet `json:"allowedMediaTypes" gore:"type:text;not null;default:'ARTICLE'"`
}
type RoundWritable struct {
	Name             string    `json:"name"`
	Description      string    `json:"description" gorm:"type:text"`
	StartDate        time.Time `json:"startDate" gorm:"type:datetime"`
	EndDate          time.Time `json:"endDate" gorm:"type:datetime"`
	IsOpen           bool      `json:"isOpen" gorm:"default:true"`
	IsPublic         bool      `json:"isPublic" gorm:"default:false"`
	DependsOnRoundID *string   `json:"dependsOnRoundId" gorm:"default:null"`
	DependsOnRound   *Round    `json:"-" gorm:"foreignKey:DependsOnRoundID"`
	Serial           int       `json:"serial" gorm:"default:0"`
	RoundRestrictions
}
type Round struct {
	RoundID          IDType      `json:"roundId" gorm:"primaryKey"`
	CampaignID       IDType      `json:"campaignId" gorm:"index;cascade:OnUpdate,OnDelete"`
	CreatedAt        *time.Time  `json:"createdAt" gorm:"-<-:create"`
	CreatedByID      IDType      `json:"createdById" gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TotalSubmissions int         `json:"totalSubmissions" gorm:"default:0"`
	Status           RoundStatus `json:"status"`
	Campaign         *Campaign   `json:"-"`
	Creator          *User       `json:"-" gorm:"foreignKey:CreatedByID"`
	LatestTaskID     *IDType     `json:"latestTaskId" gorm:"default:null"`
	RoundWritable
	Roles []Role `json:"roles"`
}
type RoundFilter struct {
	CampaignID IDType      `form:"campaignId"`
	Status     RoundStatus `form:"status"`
	CommonFilter
}
type RoundRepository struct{}

func NewRoundRepository() *RoundRepository {
	return &RoundRepository{}
}
func (r *RoundRepository) Create(conn *gorm.DB, round *Round) (*Round, error) {
	result := conn.Create(round)
	if result.Error != nil {
		return nil, result.Error
	}
	return round, nil
}
func (r *RoundRepository) Update(conn *gorm.DB, round *Round) (*Round, error) {
	result := conn.Save(round)
	if result.Error != nil {
		return nil, result.Error
	}
	return round, nil
}
func (r *RoundRepository) FindAll(conn *gorm.DB, filter *RoundFilter) ([]Round, error) {
	var rounds []Round
	where := &Round{}
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
func (r *RoundRepository) FindByID(conn *gorm.DB, id IDType) (*Round, error) {
	round := &Round{}
	where := &Round{RoundID: IDType(id)}
	result := conn.First(round, where)
	return round, result.Error
}
