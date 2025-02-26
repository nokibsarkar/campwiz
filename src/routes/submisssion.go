package routes

import (
	"nokib/campwiz/database"
	"nokib/campwiz/services"

	"github.com/gin-gonic/gin"
)

// ListAllSubmissions godoc
// @Summary List all submissions
// @Description get all submissions
// @Produce  json
// @Success 200 {object} ResponseList[database.Submission]
// @Router /submission/ [get]
// @param SubmissionListFilter query database.SubmissionListFilter false "Filter the submissions"
// @Tags Submission
// @Error 400 {object} ResponseError
func ListAllSubmissions(c *gin.Context) {
	filter := &database.SubmissionListFilter{}
	err := c.ShouldBindQuery(filter)
	if err != nil {
		c.JSON(400, ResponseError{
			Detail: "Invalid query",
		})
		return
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	submission_service := services.NewSubmissionService()
	submissions, err := submission_service.ListAllSubmissions(filter)

	if err != nil {
		c.JSON(400, ResponseError{
			Detail: "Error listing submissions",
		})
		return
	}
	continueToken := ""
	previousToken := ""
	if len(submissions) > 0 {
		continueToken = string(submissions[len(submissions)-1].SubmissionID)
		previousToken = string(submissions[0].SubmissionID)
	}
	c.JSON(200, ResponseList[database.Submission]{
		Data:          submissions,
		ContinueToken: continueToken,
		PreviousToken: previousToken,
	})

}
func CreateDraftSubmission(c *gin.Context) {
	// ...
}
func CreateLateDraftSubmission(c *gin.Context) {
	// ...
}
func GetDraftSubmission(c *gin.Context) {
	// ...
}
func GetSubmission(c *gin.Context) {
	// ...
}
func DeleteSubmission(c *gin.Context) {
	// ...
}
func GetEvaluation(c *gin.Context) {
	// ...
}
func CreateSubmission(c *gin.Context) {
	// ...
}
func CreateLateSubmission(c *gin.Context) {
	// ...
}
func EvaluateSubmission(c *gin.Context) {
	// ...
}

func NewSubmissionRoutes(parent *gin.RouterGroup) {
	r := parent.Group("/submission")
	r.GET("/", ListAllSubmissions)
	r.POST("/draft", CreateDraftSubmission)
	r.POST("/draft/late", CreateLateDraftSubmission)
	r.GET("/draft/:id", GetDraftSubmission)
	r.GET("/:id", GetSubmission)
	r.DELETE("/:id", DeleteSubmission)
	r.GET("/:id/judge", GetEvaluation)
	r.POST("/", CreateSubmission)
	r.POST("/late", CreateLateSubmission)
	r.POST("/:id/judge", EvaluateSubmission)
}
