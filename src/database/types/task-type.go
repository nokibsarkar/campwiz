package types

import "strings"

type TaskIDType string

const TASK_ID_PREFIX = "t"

func (id *TaskIDType) GormDataType() string {
	return "varchar(255)"
}

func NewTaskIDType(id string) (TaskIDType, error) {
	if !strings.HasPrefix(id, TASK_ID_PREFIX) {
		return "", ErrorType
	}
	taskIDType := TaskIDType(id)
	return taskIDType, nil
}
