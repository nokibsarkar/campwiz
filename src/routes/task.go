package routes

import (
	"io"
	"nokib/campwiz/database"
	"nokib/campwiz/database/cache"
	"nokib/campwiz/services"

	"github.com/gin-gonic/gin"
)

// GetTaskById godoc
// @Summary Get a task by ID
// @Description The task represents a background job that can be run by the system
// @Produce  json
// @Success 200 {object} ResponseSingle[services.TaskResponse]
// @Router /task/{taskId} [get]
// @Tags Task
// @Param taskId path string true "The task ID"
// @Error 400 {object} ResponseError
// @Error 404 {object} ResponseError
func GetTaskById(c *gin.Context, sess *cache.Session) {
	defer HandleError("GetTaskById")
	taskId := c.Param("taskId")
	if taskId == "" {
		c.JSON(400, ResponseError{Detail: "Invalid request : Task ID is required"})
		return
	}
	task_service := services.NewTaskService()
	task, err := task_service.GetTask(database.IDType(taskId))
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Error getting task : " + err.Error()})
		return
	}
	if task == nil {
		c.JSON(404, ResponseError{Detail: "Task not found"})
		return
	}
	c.JSON(200, ResponseSingle[database.Task]{Data: *task})
}

// GetTaskByIDStream godoc
// @Summary Get a task by ID but stream the response
// @Description The task represents a background job that can be run by the system. This endpoint streams the response
// @Produce  json
// @Router /task/{taskId}/stream [get]
// @Tags Task
// @Param taskId path string true "The task ID"
// @Error 400 {object} ResponseError
// @Error 404 {object} ResponseError
func GetTaskByIDStream(c *gin.Context, sess *cache.Session) {
	defer HandleError("GetTaskById")
	taskId := c.Param("taskId")
	if taskId == "" {
		c.JSON(400, ResponseError{Detail: "Invalid request : Task ID is required"})
		return
	}
	task_service := services.NewTaskService()
	c.Stream(func(w io.Writer) bool {
		task, err := task_service.GetTask(database.IDType(taskId))
		if err != nil {
			c.SSEvent("error", err.Error())
			return false
		}
		if task == nil {
			c.SSEvent("error", "Task not found")
			return false
		}
		c.SSEvent("task", task)
		if task.Status == "completed" || task.Status == "failed" {
			return false
		}
		return true
	})
}
func NewTaskRoutes(p *gin.RouterGroup) {
	task := p.Group("/task")
	task.GET("/:taskId", WithSession(GetTaskById))
	task.GET("/:taskId/stream", WithSession(GetTaskByIDStream))

}
