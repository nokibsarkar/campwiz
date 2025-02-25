package routes

import (
	"nokib/campwiz/consts"
	"nokib/campwiz/database"
	"nokib/campwiz/database/cache"
	"nokib/campwiz/services"

	"github.com/gin-gonic/gin"
)

// CreateRound godoc
// @Summary Create a new round
// @Description Create a new round for a campaign
// @Produce  json
// @Success 200 {object} ResponseSingle[database.CampaignRound]
// @Router /round/ [post]
// @Param roundRequest body services.RoundRequest true "The round request"
// @Tags Round
// @Error 400 {object} ResponseError
func CreateRound(c *gin.Context, sess *cache.Session) {
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
	round, err := round_service.CreateRound(&requestedRounds)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Error creating round : " + err.Error()})
		return
	}
	c.JSON(200, ResponseSingle[database.CampaignRound]{Data: *round})

}

// ListAllRounds godoc
// @Summary List all rounds
// @Description get all rounds
// @Produce  json
// @Success 200 {object} ResponseList[database.CampaignRound]
// @Router /round/ [get]
// @param RoundFilter query database.RoundFilter false "Filter the rounds"
// @Tags Round
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

// ImportFromCommons godoc
// @Summary Import images from commons
// @Description The user would provide a round ID and a list of commons categories and the system would import images from those categories
// @Produce  json
// @Success 200 {object} ResponseSingle[services.RoundImportSummary]
// @Router /round/import/{roundId} [post]
// @Param roundId path string true "The round ID"
// @Param ImportFromCommons body services.ImportFromCommonsPayload true "The import from commons request"
// @Tags Round
// @Error 400 {object} ResponseError
func ImportFromCommons(c *gin.Context, sess *cache.Session) {
	roundId := c.Param("roundId")
	if roundId == "" {
		c.JSON(400, ResponseError{Detail: "Invalid request : Round ID is required"})
	}
	req := &services.ImportFromCommonsPayload{}
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Error Decoding : " + err.Error()})
		return
	}
	round_service := services.NewRoundService()
	if len(req.Categories) == 0 {
		c.JSON(400, ResponseError{Detail: "Invalid request : No categories provided"})
		return
	}
	round, err := round_service.ImportFromCommons(roundId, req.Categories)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Failed to import images : " + err.Error()})
		return
	}
	c.JSON(200, ResponseSingle[*services.RoundImportSummary]{Data: round})
}

// GetImportStatus godoc
// @Summary Get the import status about a round
// @Description It would be used as a server sent event stream to broadcast on the frontend about current status of the round
// @Produce  json
// @Success 200 {object} ResponseSingle[services.RoundImportSummary]
// @Router /round/import/{roundId} [get]
// @Param roundId path string true "The round ID"
// @Tags Round
// @Error 400 {object} ResponseError
func GetImportStatus(c *gin.Context, sess *cache.Session) {

}

func NewRoundRoutes(parent *gin.RouterGroup) {
	r := parent.Group("/round")
	r.POST("/", WithPermission(consts.PermissionCreateRound, CreateRound))
	r.GET("/", WithSession(ListAllRounds))
	r.GET("/import/:roundId", SSEHeadersMiddleware(), WithPermission(consts.PermissionCreateRound, GetImportStatus))
	r.POST("/import/:roundId/commons", WithPermission(consts.PermissionCreateRound, ImportFromCommons))
}
