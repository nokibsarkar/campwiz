package cache

import (
	"nokib/campwiz/database"
	idgenerator "nokib/campwiz/services/idGenerator"

	"gorm.io/gorm"
)

type Assignments struct {
	EvaluationID       database.IDType ` gorm:"primaryKey"`
	SubmissionID       database.IDType `gorm:"uniqueIndex:idx_submission_judge"`
	JudgeID            database.IDType `gorm:"uniqueIndex:idx_submission_judge"`
	DistributionTaskID database.IDType
}

func ExportToCache(tx *gorm.DB, cacheTx *gorm.DB, filter *database.EvaluationFilter) ([]*Assignments, error) {
	var evaluations []*Assignments
	condition := &database.Evaluation{}
	stmt := tx.Model(&database.Evaluation{})
	if filter != nil {
		if filter.RoundID != "" {
			condition.SubmissionID = filter.RoundID
		}
		if filter.ParticipantID != "" {
			condition.ParticipantID = filter.ParticipantID
		}
		if filter.JuryRoleID != "" {
			condition.JudgeID = filter.JuryRoleID
		}
		if filter.Limit > 0 {
			stmt = stmt.Limit(max(5, filter.Limit))
		}
	}
	result := stmt.Where(condition).Find(&evaluations)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(evaluations) == 0 {
		return evaluations, nil
	}
	for _, evaluation := range evaluations {
		evaluation.DistributionTaskID = idgenerator.GenerateID("t")
	}
	cacheTx.Create(evaluations)
	return evaluations, result.Error
}
