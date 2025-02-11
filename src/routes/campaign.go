package routes

import (
	"nokib/campwiz/database"
	"nokib/campwiz/database/cache"
	"nokib/campwiz/services"

	"github.com/gin-gonic/gin"
)

// ListAllCampaigns godoc
// @Summary List all campaigns
// @Description get all campaigns
// @Produce  json
// @Success 200 {object} ResponseList[database.Campaign]
// @Router /campaign/ [get]
// @Tags Campaign
// @Error 400 {object} ResponseError
func ListAllCampaigns(c *gin.Context) {
	campaignService := services.NewCampaignService()
	campaignList := campaignService.GetAllCampaigns()
	c.JSON(200, ResponseList[database.Campaign]{Data: campaignList})
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
func CreateCampaign(c *gin.Context, sess *cache.Session) {
	createRequest := &services.CampaignRequest{}
	err := c.BindJSON(createRequest)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Invalid request"})
		return
	}
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
	r.POST("/", WithSession(CreateCampaign))
	r.POST("/:id", UpdateCampaign)
	r.GET("/:id/result", GetCampaignResult)
	r.GET("/:id/submissions", GetCampaignSubmissions)
	r.GET("/:id/next", GetNextSubmission)
	r.POST("/:id/status", ApproveCampaign)
	r.POST("/:id/fountain", ImportEvaluationFromFountain)

}
