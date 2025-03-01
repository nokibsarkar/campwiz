package types

import (
	"strings"
)

type RoleIDType string

const ROLE_ID_PREFIX = "j"

func (id *RoleIDType) GormDataType() string {
	return "varchar(255)"
}

func NewRoleIDType(id string) (RoleIDType, error) {
	if !strings.HasPrefix(id, ROLE_ID_PREFIX) {
		return "", ErrorType
	}
	roleIDType := RoleIDType(id)
	return roleIDType, nil
}
