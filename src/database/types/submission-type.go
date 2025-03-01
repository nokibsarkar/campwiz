package types

import "strings"

type SubmissionIDType string

const SUBMISSION_ID_PREFIX = "s"

func (id *SubmissionIDType) GormDataType() string {
	return "varchar(255)"
}

func NewSubmissionIDType(id string) (SubmissionIDType, error) {
	if !strings.HasPrefix(id, SUBMISSION_ID_PREFIX) {
		return "", ErrorType
	}
	submissionIDType := SubmissionIDType(id)
	return submissionIDType, nil
}
