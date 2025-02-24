package routes

import (
	"nokib/campwiz/consts"
	"nokib/campwiz/database"
	"nokib/campwiz/database/cache"
	"nokib/campwiz/services"

	"github.com/gin-gonic/gin"
)

// CreateBatchFromCommonsCategory godoc
// @Summary Create a batch of images from commons category
// @Description The user would provide a list of commons categories and the system would create a batch from the images in those categories
// @Produce  json
// @Success 200 {object} ResponseSingle[services.BatchCreationResult]
// @Router /batch/create/commons [post]
// @Param CreateFromCommons body services.CreateFromCommons true "The batch creation request"
// @Tags Batch
// @Error 400 {object} ResponseError
func CreateBatchFromCommonsCategory(c *gin.Context, sess *cache.Session) {
	req := &services.CreateFromCommons{
		CreatedBy: sess.UserID,
	}
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Invalid request : " + err.Error()})
		return
	}
	req.Categories = []string{"Category:Quality images from Wiki Loves Folklore 2025"}
	batch_service := services.NewBatchService()
	batch, err := batch_service.CreateBatchFromCommonsCategory(req)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Failed to create batch : " + err.Error()})
		return
	}
	c.JSON(200, ResponseSingle[*services.BatchCreationResult]{Data: batch})
}

// CreateBatchFromCsv godoc
// @Summary Create a batch of images from a CSV file
// @Description The user would provide a CSV file with a list of image names and the system would create a batch from the images in those URLs
// @Produce  json
// @Success 200 {object} ResponseSingle[services.BatchCreationResult]
// @Router /batch/create/csv [post]
// @Param file formData file true "The CSV file"
// @Tags Batch
// @Error 400 {object} ResponseError
func CreateBatchFromCsv(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

// CreateBatchFromPagePile godoc
// @Summary Create a batch of images from a pagepile
// @Description The user would provide a pagepile ID and the system would create a batch from the images in that pagepile
// @Produce  json
// @Success 200 {object} ResponseSingle[services.BatchCreationResult]
// @Router /batch/create/pagepile/{pagepileId} [post]
// @Param pagepileId path string true "The pagepile ID"
// @Tags Batch
// @Error 400 {object} ResponseError
func CreateBatchFromPagePile(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

// CreateBatchFromWikipageList godoc
// @Summary Create a batch of images from a list of wikipages
// @Description The user would provide a wikipage and the system would create a batch from the images listed in that wikipage
// @Produce  json
// @Success 200 {object} ResponseSingle[services.BatchCreationResult]
// @Router /batch/create/wikipage [post]
// @Param wikipage formData string true "The wikipage"
// @Param language formData string true "The language of the wikipage"
// @Tags Batch
// @Error 400 {object} ResponseError
func CreateBatchFromWikipageList(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

// CreateBatchFromAnotherRoundOutput god
// @Summary Create a batch of images from another round output
// @Description The user would provide a round ID and the system would create a batch from the images in that round
// @Produce  json
// @Success 200 {object} ResponseSingle[services.BatchCreationResult]
// @Router /batch/create/round/{roundId} [post]
// @Param roundId path string true "The round ID"
// @Tags Batch
// @Error 400 {object} ResponseError
func CreateBatchFromAnotherRoundOutput(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
func DistributeTasks(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

// Get a batch
// @Summary Get a batch
// @Description Get a batch
// @Produce  json
// @Success 200 {object} ResponseSingle[database.Batch]
// @Router /batch/{batchId} [get]
// @Param batchId path string true "The batch ID"
// @Tags Batch
// @Error 404 {object} ResponseError
func GetBatch(c *gin.Context) {
	batchId := c.Param("batchId")
	if batchId == "" {
		c.JSON(400, ResponseError{Detail: "Invalid request : Batch ID is required"})
	}
	batch_service := services.NewBatchService()
	batch, err := batch_service.GetBatchByID(batchId)
	if err != nil {
		c.JSON(404, ResponseError{Detail: "Failed to get batch : " + err.Error()})
		return
	}
	c.JSON(200, ResponseSingle[*database.Batch]{Data: batch})
}

// GetBatchImages godoc
// @Summary Get images of a batch
// @Description Get images of a batch
// @Produce  json
// @Success 200 {object} ResponseList[database.Image]
// @Router /batch/{batchId}/images [get]
// @Param batchId path string true "The batch ID"
// @Tags Batch
// @Error 400 {object} ResponseError
func GetBatchImages(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

// ListBatches godoc
// @Summary List batches
// @Description List batches
// @Produce  json
// @Success 200 {object} ResponseList[database.Batch]
// @Router /batch/ [get]
// @Param filter query database.BatchFilter false "The batch filter"
// @Tags Batch
// @Error 400 {object} ResponseError
func ListBatches(c *gin.Context) {
	filter := &database.BatchFilter{}
	err := c.ShouldBindQuery(filter)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Invalid request : " + err.Error()})
		return
	}
	batch_service := services.NewBatchService()
	batches, err := batch_service.GetAllBatches(filter)
	if err != nil {
		c.JSON(400, ResponseError{Detail: "Failed to get batches : " + err.Error()})
		return
	}
	c.JSON(200, ResponseList[database.Batch]{Data: batches})
}
func NewBatchRouter(parent *gin.RouterGroup) {
	router := parent.Group("/batch")
	router.POST("/create/commons", WithPermission(consts.PermissionUpdateCampaign, CreateBatchFromCommonsCategory))
	router.POST("/create/csv", CreateBatchFromCsv)
	router.POST("/create/pagepile/:pagepileId", CreateBatchFromPagePile)
	router.POST("/create/wikipage", CreateBatchFromWikipageList)
	router.POST("/create/round/:roundId", CreateBatchFromAnotherRoundOutput)
	router.POST("/distribute/:batchId", DistributeTasks)
	router.GET("/:batchId", GetBatch)
	router.GET("/:batchId/images", GetBatchImages)
	router.GET("/", ListBatches)
}
