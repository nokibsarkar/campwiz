package routes

import (
	"nokib/campwiz/database"
	"nokib/campwiz/database/cache"
	"nokib/campwiz/services"

	"github.com/gin-gonic/gin"
)

// BulkAddRound godoc
// @Summary Add multiple rounds to a campaign
// @Description Add multiple rounds to a campaign
// @Produce  json
// @Success 200 {object} ResponseList[database.CampaignRound]
// @Router /round/bulk-add [post]
// @param roundRequest body services.RoundRequest true "The round request"
// @Tags CampaignRound
// @Error 400 {object} ResponseError
func BulkAddRound(c *gin.Context, sess *cache.Session) {
	defer HandleError("BulkAddRound")
	requestedRounds := services.RoundRequest{
		CreatedByID: sess.UserID,
	}
	err := c.ShouldBindJSON(&requestedRounds)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Invalid request : " + err.Error()})
		return
	}
	round_service := services.NewRoundService()
	rounds, err := round_service.CreateRound(&requestedRounds)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Error creating round : " + err.Error()})
		return
	}
	c.JSON(200, ResponseList[database.CampaignRound]{Data: rounds})

}

// ListAllRounds godoc
// @Summary List all rounds
// @Description get all rounds
// @Produce  json
// @Success 200 {object} ResponseList[database.CampaignRound]
// @Router /round/ [get]
// @param RoundFilter query database.RoundFilter false "Filter the rounds"
// @Tags CampaignRound
// @Error 400 {object} ResponseError
func ListAllRounds(c *gin.Context, sess *cache.Session) {
	defer HandleError("ListAllRounds")
	filter := &database.RoundFilter{}
	err := c.ShouldBindQuery(filter)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Invalid request : " + err.Error()})
		return
	}
	round_service := services.NewRoundService()
	rounds, err := round_service.ListAllRounds(filter)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Error listing rounds : " + err.Error()})
		return
	}
	c.JSON(200, ResponseList[database.CampaignRound]{Data: rounds})
}
func NewRoundRoutes(parent *gin.RouterGroup) {
	r := parent.Group("/round")
	r.POST("/bulk-add", WithSession(BulkAddRound))
	r.GET("/", WithSession(ListAllRounds))
}
