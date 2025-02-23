package database

type Jury struct {
	JuryID         uint64    `json:"juryId" gorm:"primaryKey"`
	UserID         uint64    `json:"userId"`
	CampaignID     uint64    `json:"campaignId"`
	TotalAssigned  int       `json:"totalAssigned"`
	TotalEvaluated int       `json:"totalEvaluated"`
	TotalScore     int       `json:"totalScore"`
	Campaign       *Campaign `json:"-"`
	User           *User     `json:"-"`
}
