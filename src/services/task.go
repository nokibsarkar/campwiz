package services

import "nokib/campwiz/database"

type TaskService struct{}

func NewTaskService() *TaskService {
	return &TaskService{}
}

type TaskRequest struct {
	Type                 string `json:"type"`
	AssociatedRoundID    string `json:"roundId"`
	AssociatedCampaignID string `json:"campaignId"`
	AssociatedUserID     string `json:"userId"`
	Data                 any    `json:"data"`
	CreatedByID          string `json:"createdById"`
	HandlerFunc          func(c ...any)
}
type TaskResponse struct {
	database.Task
	handlerFunc func(c ...any)
}
type TaskFilter struct {
	Type                 string `json:"type"`
	AssociatedRoundID    string `json:"roundId"`
	AssociatedCampaignID string `json:"campaignId"`
	AssociatedUserID     string `json:"userId"`
	Status               string `json:"status"`
}

func (t *TaskResponse) Start() {
	go t.handlerFunc()
}

func (t *TaskService) CreateTask(TaskRequest, handlerFunc func(c ...any)) (*TaskResponse, error) {

	return nil, nil
}
func (t *TaskService) GetTask(taskId string) (*database.Task, error) {

	return nil, nil
}
func (t *TaskService) ListTasks(filter *TaskFilter) ([]database.Task, error) {
	return nil, nil
}
