package routes

import "github.com/gin-gonic/gin"

func ListAllSubmissions(c *gin.Context) {
	// ...
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
