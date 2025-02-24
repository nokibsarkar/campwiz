package database

import (
	"nokib/campwiz/consts"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ArticleCampaignRestrictions struct {
	MaximumSubmissionOfSameArticle int    `json:"maximumSubmissionOfSameArticle"`
	AllowExpansions                bool   `json:"allowExpansions"`
	AllowCreations                 bool   `json:"allowCreations"`
	MinimumTotalBytes              int    `json:"minimumTotalBytes"`
	MinimumTotalWords              int    `json:"minimumTotalWords"`
	MinimumAddedBytes              int    `json:"minimumAddedBytes"`
	MinimumAddedWords              int    `json:"minimumAddedWords"`
	SecretBallot                   bool   `json:"secretBallot"`
	AllowJuryToParticipate         bool   `json:"allowJuryToParticipate"`
	AllowMultipleJudgement         bool   `json:"allowMultipleJudgement"`
	Blacklist                      string `json:"blacklist"`
}
type ImageCampaignRestrictions struct {
	MaximumSubmissionOfSameImage int `json:"maximumSubmissionOfSameImage"`
	MinimumTotalImageSize        int `json:"minimumTotalImageSize"`
}
type MediaCampaignRestrictions struct {
	ImageCampaignRestrictions
}
type CampaignWithWriteableFields struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	StartDate   time.Time       `json:"startDate"`
	EndDate     time.Time       `json:"endDate"`
	Language    consts.Language `json:"language"`
	Rules       string          `json:"rules"`
	Image       string          `json:"image"`
}
type Campaign struct {
	ID string `gorm:"primaryKey" json:"id"`
	// read only
	CreatedAt *time.Time `json:"createdAt" gorm:"-<-:create"`
	CreatedBy string     `json:"createdBy"`
	CampaignWithWriteableFields
}
type CampaignFilter struct {
	IDs       []string `form:"ids,omitEmpty"`
	Limit     int      `form:"limit,omitEmpty"`
	NextToken string   `form:"nextToken,omitEmpty"`
}
type CampaignRepository struct{}

func NewCampaignRepository() *CampaignRepository {
	return &CampaignRepository{}
}
func (c *CampaignRepository) Create(conn *gorm.DB, campaign *Campaign) error {
	result := conn.Create(campaign)
	return result.Error
}
func (c *CampaignRepository) FindByID(conn *gorm.DB, id string) (*Campaign, error) {
	campaign := &Campaign{}
	where := &Campaign{ID: id}
	result := conn.First(campaign, where)
	return campaign, result.Error
}
func (c *CampaignRepository) ListAllCampaigns(conn *gorm.DB, query *CampaignFilter) ([]Campaign, error) {
	var campaigns []Campaign
	stmt := conn
	if query != nil {
		if query.Limit > 0 {
			stmt = stmt.Limit(query.Limit)
		}
		if len(query.IDs) > 0 {
			idCopies := []string{}
			for _, id := range query.IDs {
				if id != "" && strings.Contains(id, ",") {
					idCopies = append(idCopies, strings.Split(id, ",")...)
				} else {
					idCopies = append(idCopies, id)
				}
			}
			stmt = stmt.Where("id IN (?)", idCopies)
		}
	}
	result := stmt.Find(&campaigns)
	return campaigns, result.Error
}
func (c *CampaignRepository) Update(conn *gorm.DB, campaign *Campaign) error {
	result := conn.Save(campaign)
	return result.Error
}
