package database

type Jury struct {
	ID             uint64    `json:"id" gorm:"primaryKey"`
	UserID         uint64    `json:"userId"`
	CampaignID     uint64    `json:"campaignId"`
	Campaign       *Campaign `json:"-"`
	User           *User     `json:"-"`
	TotalAssigned  int       `json:"totalAssigned"`
	TotalEvaluated int       `json:"totalEvaluated"`
	TotalScore     int       `json:"totalScore"`
}
