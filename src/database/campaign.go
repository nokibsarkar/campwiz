package database

import (
	"nokib/campwiz/consts"
	"time"
)

/*
name : str
language : Language
start_at : datetime | date
end_at : datetime | date
description : str | None
rules : str | list[str]
blacklist : list[str] | str | None
image : str | None
maximumSubmissionOfSameArticle : int
allowExpansions : bool
minimumTotalBytes : int
minimumTotalWords : int
minimumAddedBytes : int
minimumAddedWords : int
secretBallot : bool
allowJuryToParticipate : bool
allowMultipleJudgement : bool
*/
type ArticleCampaignRestrictions struct {
	MaximumSubmissionOfSameArticle int    `json:"maximum_submission_of_same_article"`
	AllowExpansions                bool   `json:"allow_expansions"`
	AllowCreations                 bool   `json:"allow_creations"`
	MinimumTotalBytes              int    `json:"minimum_total_bytes"`
	MinimumTotalWords              int    `json:"minimum_total_words"`
	MinimumAddedBytes              int    `json:"minimum_added_bytes"`
	MinimumAddedWords              int    `json:"minimum_added_words"`
	SecretBallot                   bool   `json:"secret_ballot"`
	AllowJuryToParticipate         bool   `json:"allow_jury_to_participate"`
	AllowMultipleJudgement         bool   `json:"allow_multiple_judgement"`
	Blacklist                      string `json:"blacklist"`
}
type ImageCampaignRestrictions struct {
	MaximumSubmissionOfSameImage int `json:"maximum_submission_of_same_image"`
	MinimumTotalImageSize        int `json:"minimum_total_image_size"`
}
type MediaCampaignRestrictions struct {
	ImageCampaignRestrictions
}
type CampaignWithWriteableFields struct {
	ID          string          `gorm:"primaryKey" json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	StartDate   time.Time       `json:"start_date"`
	EndDate     time.Time       `json:"end_date"`
	Language    consts.Language `json:"language"`
	Rules       string          `json:"rules"`
	Image       string          `json:"image"`
	CreatedBy   string          `json:"created_by"`
}
type Campaign struct {
	// read only
	CreatedAt *time.Time `json:"created_at,omitEmpty" gorm:"-<-:create"`
	CampaignWithWriteableFields
}
type CampaignFilter struct {
	IDs       []uint `form:"ids,omitEmpty"`
	Limit     int    `form:"limit,omitEmpty"`
	NextToken string `form:"nextToken,omitEmpty"`
}
