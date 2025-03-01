package types

import "strings"

type EvaluationIDType string

const EVALUATION_ID_PREFIX = "e"

func (id *EvaluationIDType) GormDataType() string {
	return "varchar(255)"
}
func NewEvaluationIDType(id string) (EvaluationIDType, error) {
	if !strings.HasPrefix(id, EVALUATION_ID_PREFIX) {
		return "", ErrorType
	}
	evaluationIDType := EvaluationIDType(id)
	return evaluationIDType, nil
}
