package routes

import (
	"nokib/campwiz/database"
	"nokib/campwiz/database/cache"
	"nokib/campwiz/services"

	"github.com/gin-gonic/gin"
)

// List Evaluations godoc
// @Summary List all evaluations
// @Description get all evaluations
// @Produce  json
// @Success 200 {object} ResponseList[database.Evaluation]
// @Router /evaluation/ [get]
// @param EvaluationFilter query services.EvaluationFilter false "Filter the evaluations"
// @Tags Evaluation
// @Security ApiKeyAuth
// @Error 400 {object} ResponseError
func ListEvaluations(c *gin.Context, sess *cache.Session) {
	filter := &services.EvaluationFilter{}
	err := c.ShouldBindQuery(filter)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Invalid request : " + err.Error()})
		return
	}
	evaluation_service := services.NewEvaluationService()
	evaluations, err := evaluation_service.ListEvaluations(filter)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Error listing evaluations : " + err.Error()})
		return
	}
	c.JSON(200, ResponseList[database.Evaluation]{Data: evaluations})
}
func NewEvaluationRoutes(r *gin.RouterGroup) {
	route := r.Group("/evaluation")
	route.GET("/", WithSession(ListEvaluations))
}
