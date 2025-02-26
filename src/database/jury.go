package database

type Jury struct {
	JuryID         IDType    `json:"juryId" gorm:"primaryKey"`
	UserID         IDType    `json:"userId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CampaignID     IDType    `json:"campaignId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TotalAssigned  int       `json:"totalAssigned"`
	TotalEvaluated int       `json:"totalEvaluated"`
	TotalScore     int       `json:"totalScore"`
	Campaign       *Campaign `json:"-"`
	User           *User     `json:"-"`
}
