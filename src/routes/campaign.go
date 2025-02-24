package routes

import (
	"nokib/campwiz/consts"
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
// @param CampaignFilter query database.CampaignFilter false "Filter the campaigns"
// @Tags Campaign
// @Error 400 {object} ResponseError
func ListAllCampaigns(c *gin.Context) {
	query := &database.CampaignFilter{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Invalid request : " + err.Error()})
		return
	}
	campaignService := services.NewCampaignService()
	campaignList := campaignService.GetAllCampaigns(query)
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

// CreateCampaign creates a new campaign
// @summary Create a new campaign
// @description Create a new campaign
// @tags Campaign
// @param campaignRequest body services.CampaignCreateRequest true "The campaign request"
// @produce json
// @success 200 {object} database.Campaign
// @router /campaign/ [post]
func CreateCampaign(c *gin.Context, sess *cache.Session) {
	createRequest := &services.CampaignCreateRequest{}
	err := c.ShouldBindBodyWithJSON(createRequest)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Invalid request : " + err.Error()})
		return
	}
	createRequest.CreatedBy = sess.UserID
	camapign_service := services.NewCampaignService()
	campaign, err := camapign_service.CreateCampaign(createRequest)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Failed to create campaign : " + err.Error()})
		return
	}
	c.JSON(200, ResponseSingle[*database.Campaign]{Data: campaign})
}

// UpdateCampaign godoc
// @Summary Update a campaign
// @Description Update a campaign
// @Produce  json
// @Success 200 {object} ResponseSingle[database.Campaign]
// @Router /campaign/{id} [post]
// @Tags Campaign
// @Param id path string true "The campaign ID"
// @Param campaignRequest body services.CampaignUpdateRequest true "The campaign request"

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
	r.POST("/", WithPermission(consts.PermissionCreateCampaign, CreateCampaign))
	r.POST("/:id", UpdateCampaign)
	r.GET("/:id/result", GetCampaignResult)
	r.GET("/:id/submissions", GetCampaignSubmissions)
	r.GET("/:id/next", GetNextSubmission)
	r.POST("/:id/status", ApproveCampaign)
	r.POST("/:id/fountain", ImportEvaluationFromFountain)
}
