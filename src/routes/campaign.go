package routes

import "github.com/gin-gonic/gin"

func ListAllCampaigns(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

/*
This function will return all the timelines of all the campaigns
*/
func GetAllCampaignTimeLine(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
func GetSingleCampaign(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
func ListAllJury(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
func CreateCampaign(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
func UpdateCampaign(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
func GetCampaignResult(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
func GetCampaignSubmissions(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
func GetNextSubmission(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
func ApproveCampaign(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
func ImportEvaluationFromFountain(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

/*
NewCampaignRoutes will create all the routes for the /campaign endpoint
*/
func NewCampaignRoutes(parent *gin.RouterGroup) {
	defer HandleError("/campaign")
	r := parent.Group("/campaign")
	r.GET("/", ListAllCampaigns)
	r.GET("/timeline2", GetAllCampaignTimeLine)
	r.GET("/:id", GetSingleCampaign)
	r.GET("/jury", ListAllJury)
	r.POST("/", CreateCampaign)
	r.POST("/:id", UpdateCampaign)
	r.GET("/:id/result", GetCampaignResult)
	r.GET("/:id/submissions", GetCampaignSubmissions)
	r.GET("/:id/next", GetNextSubmission)
	r.POST("/:id/status", ApproveCampaign)
	r.POST("/:id/fountain", ImportEvaluationFromFountain)

}
