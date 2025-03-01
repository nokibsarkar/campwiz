package database

import (
	"time"

	"gorm.io/gorm"
)

type EvaluationType string

const (
	EvaluationTypeRanking EvaluationType = "ranking"
	EvaluationTypeScore   EvaluationType = "score"
	EvaluationTypeBinary  EvaluationType = "binary"
)

type Evaluation struct {
	EvaluationID  IDType         `json:"evaluationId" gorm:"primaryKey"`
	SubmissionID  IDType         `json:"submissionId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	JudgeID       IDType         `json:"judgeId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;uniqueIndex:idx_unique_vote_position"`
	ParticipantID IDType         `json:"participantId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Type          EvaluationType `json:"type"`
	// Applicable if the evaluation type is score, it would be between 0-100
	VoteScore *int `json:"score" gorm:"null;:null"`
	// Applicable if the evaluation type is binary, it would be 0 to Number of submissions in this round
	// The pair (JudgeID, VotePosition) should be unique (i.e. a judge can't vote for the same position twice)
	VotePosition *int `json:"votePosition" gorm:"null;default:null;uniqueIndex:idx_unique_vote_position"`
	// Applicable if the evaluation type is binary, it would be true or false
	VotePassed         *bool       `json:"votePassed" gorm:"null;default:null;"`
	Comment            string      `json:"comment" gorm:"default:null"`
	Serial             uint        `json:"serial"`
	Submission         *Submission `json:"-"`
	Participant        *User       `json:"-" gorm:"foreignKey:ParticipantID"`
	Judge              *Role       `json:"-"`
	CreatedAt          *time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt          *time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	EvaluatedAt        *time.Time  `json:"evaluatedAt" gorm:"type:datetime"`
	DistributionTaskID IDType      `json:"distributionTaskId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
type EvaluationFilter struct {
	Type          EvaluationType `form:"type"`
	RoundID       IDType         `form:"roundId"`
	CampaignID    IDType         `form:"campaignId"`
	ParticipantID IDType         `form:"userId"`
	Evaluated     *bool          `form:"status"`
	SubmissionID  IDType         `form:"submissionId"`
	JuryID        IDType         `form:"juryId"`
	CommonFilter
}
type EvaluationRepository struct{}

func NewEvaluationRepository() *EvaluationRepository {
	return &EvaluationRepository{}
}
func (r *EvaluationRepository) CreateEvaluation(tx *gorm.DB, evaluation *Evaluation) error {
	result := tx.Create(evaluation)
	return result.Error
}
func (r *EvaluationRepository) FindEvaluationByID(tx *gorm.DB, evaluationID IDType) (*Evaluation, error) {
	evaluation := &Evaluation{}
	result := tx.First(evaluation, &Evaluation{EvaluationID: evaluationID})
	return evaluation, result.Error
}
func (r *EvaluationRepository) ListAllEvaluations(tx *gorm.DB, filter *EvaluationFilter) ([]Evaluation, error) {
	var evaluations []Evaluation
	condition := &Evaluation{}
	stmt := tx
	if filter != nil {
		s := &Submission{}
		if filter.RoundID != "" || filter.CampaignID != "" {
			if filter.RoundID != "" {
				s.CurrentRoundID = filter.RoundID
			}
			if filter.CampaignID != "" {
				s.CampaignID = filter.CampaignID
			}
			stmt = tx.Joins("Submission", tx.Where(s))
		}
		if filter.ParticipantID != "" {
			condition.ParticipantID = filter.ParticipantID
		}
		if filter.Type != "" {
			condition.Type = filter.Type
		}
		if filter.Evaluated != nil {
			if *filter.Evaluated {
				stmt = stmt.Where("evaluated_at IS NOT NULL")
			} else {
				stmt = stmt.Where("evaluated_at IS NULL")
			}
		}
		if filter.SubmissionID != "" {
			condition.SubmissionID = filter.SubmissionID
		}
		if filter.JuryID != "" {
			condition.JudgeID = filter.JuryID
		}
		if filter.ContinueToken != "" {
			stmt = stmt.Where("evaluation_id > ?", filter.ContinueToken)
		}
		if filter.PreviousToken != "" {
			stmt = stmt.Where("evaluation_id < ?", filter.PreviousToken)
		}
		if filter.Limit > 0 {
			stmt = stmt.Limit(max(5, filter.Limit))
		}
	}
	result := stmt.Where(condition).Find(&evaluations)
	return evaluations, result.Error
}
func (r *EvaluationRepository) UpdateEvaluation(tx *gorm.DB, evaluation *Evaluation) error {
	result := tx.Save(evaluation)
	return result.Error
}
