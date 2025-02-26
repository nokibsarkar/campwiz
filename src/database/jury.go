package database

type Jury struct {
	JuryID         IDType    `json:"juryId" gorm:"primaryKey"`
	UserID         IDType    `json:"userId"`
	CampaignID     IDType    `json:"campaignId"`
	TotalAssigned  int       `json:"totalAssigned"`
	TotalEvaluated int       `json:"totalEvaluated"`
	TotalScore     int       `json:"totalScore"`
	Campaign       *Campaign `json:"-"`
	User           *User     `json:"-"`
}
