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
	Judge              *Jury       `json:"-"`
	CreatedAt          *time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt          *time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	EvaluatedAt        *time.Time  `json:"evaluatedAt" gorm:"type:datetime"`
	DistributionTaskID IDType      `json:"distributionTaskId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
type EvaluationFilter struct {
	Type                 EvaluationType `form:"type"`
	AssociatedRoundID    IDType         `form:"roundId"`
	AssociatedCampaignID IDType         `form:"campaignId"`
	AssociatedUserID     IDType         `form:"userId"`
	Status               bool           `form:"status"`
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
	if filter != nil {

	}
	where := tx.Where(condition)
	if filter.ContinueToken != "" {
		where = where.Where("evaluation_id > ?", filter.ContinueToken)
	}
	stmt := where
	if filter.Limit > 0 {
		stmt = stmt.Limit(max(100, filter.Limit))
	}
	result := stmt.Find(&evaluations)
	return evaluations, result.Error
}
func (r *EvaluationRepository) UpdateEvaluation(tx *gorm.DB, evaluation *Evaluation) error {
	result := tx.Save(evaluation)
	return result.Error
}
