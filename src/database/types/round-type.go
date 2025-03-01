package types

import "strings"

type RoundIDType string

const ROUND_ID_PREFIX = "r"

func (id *RoundIDType) GormDataType() string {
	return "varchar(255)"
}

func NewRoundIDType(id string) (RoundIDType, error) {
	if !strings.HasPrefix(id, ROUND_ID_PREFIX) {
		return "", ErrorType
	}
	roundIDType := RoundIDType(id)
	return roundIDType, nil
}
