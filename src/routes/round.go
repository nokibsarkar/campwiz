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
// @Success 200 {object} ResponseSingle[database.Round]
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
	c.JSON(200, ResponseSingle[database.Round]{Data: *round})
}

// ListAllRounds godoc
// @Summary List all rounds
// @Description get all rounds
// @Produce  json
// @Success 200 {object} ResponseList[database.Round]
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
	c.JSON(200, ResponseList[database.Round]{Data: rounds})
}

// ImportFromCommons godoc
// @Summary Import images from commons
// @Description The user would provide a round ID and a list of commons categories and the system would import images from those categories
// @Produce  json
// @Success 200 {object} ResponseSingle[database.Task]
// @Router /round/import/{roundId}/commons [post]
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
	task, err := round_service.ImportFromCommons(database.IDType(roundId), req.Categories)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Failed to import images : " + err.Error()})
		return
	}
	c.JSON(200, ResponseSingle[*database.Task]{Data: task})
}

// UpdateRoundDetails godoc
// @Summary Update the details of a round
// @Description Update the details of a round
// @Produce  json
// @Success 200 {object} ResponseSingle[database.Round]
// @Router /round/{roundId} [post]
// @Param roundId path string true "The round ID"
// @Param roundRequest body services.RoundRequest true "The round request"
// @Tags Round
// @Error 400 {object} ResponseError
func UpdateRoundDetails(c *gin.Context, sess *cache.Session) {
	roundId := c.Param("roundId")
	if roundId == "" {
		c.JSON(400, ResponseError{Detail: "Invalid request : Round ID is required"})
	}
	req := &services.RoundRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Error Decoding : " + err.Error()})
		return
	}
	round_service := services.NewRoundService()
	round, err := round_service.UpdateRoundDetails(database.IDType(roundId), req)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Failed to update round : " + err.Error()})
		return
	}
	c.JSON(200, ResponseSingle[*database.Round]{Data: round})
}

// DistributeEvaluations godoc
// @Summary Distribute evaluations to juries
// @Description Distribute evaluations to juries
// @Produce  json
// @Success 200 {object} ResponseSingle[database.Task]
// @Router /round/distribute/{roundId} [post]
// @Param roundId path string true "The round ID"
// @Param DistributionRequest body services.DistributionRequest true "The distribution request"
// @Tags Round
// @Error 400 {object} ResponseError
func DistributeEvaluations(c *gin.Context, sess *cache.Session) {
	roundId := c.Param("roundId")
	if roundId == "" {
		c.JSON(400, ResponseError{Detail: "Invalid request : Round ID is required"})
	}
	distributionReq := &services.DistributionRequest{}
	err := c.ShouldBind(distributionReq)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Error Decoding : " + err.Error()})
		return
	}
	round_service := services.NewRoundService()
	task, err := round_service.DistributeEvaluations(sess.UserID, database.IDType(roundId), distributionReq)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Failed to distribute evaluations : " + err.Error()})
		return
	}
	c.JSON(200, ResponseSingle[*database.Task]{Data: task})
}
func NewRoundRoutes(parent *gin.RouterGroup) {
	r := parent.Group("/round")
	r.GET("/", WithSession(ListAllRounds))
	r.POST("/", WithPermission(consts.PermissionCreateCampaign, CreateRound))
	r.POST("/:roundId", WithPermission(consts.PermissionCreateCampaign, UpdateRoundDetails))
	r.POST("/import/:roundId/commons", WithPermission(consts.PermissionCreateCampaign, ImportFromCommons))
	r.POST("/distribute/:roundId", WithPermission(consts.PermissionCreateCampaign, DistributeEvaluations))
}
