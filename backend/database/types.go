package database

// This represents the username provided by wikimedia, it should be normalized already by wikimedia
type WikimediaUsernameType string
type UserIDType string
type CampaignIDType string
type RoundIDType string
type TaskIDType string
type SubmissionIDType string
type RoleIDType string
type EvaluationIDType string
type IDType string

func (i *IDType) String() string {
	return string(*i)
}
func (i *IDType) GormDataType() string {
	return "varchar(255)"
}
