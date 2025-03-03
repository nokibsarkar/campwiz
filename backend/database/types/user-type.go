package types

import "strings"

type UserIDType string

const USER_ID_PREFIX = "u"

func (id *UserIDType) GormDataType() string {
	return "varchar(255)"
}

func NewUserIDType(id string) (UserIDType, error) {
	if !strings.HasPrefix(id, USER_ID_PREFIX) {
		return "", ErrorType
	}
	userIDType := UserIDType(id)
	return userIDType, nil
}
