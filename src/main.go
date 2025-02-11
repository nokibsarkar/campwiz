package main

import (
	"nokib/campwiz/database"
	"nokib/campwiz/database/cache"
	"nokib/campwiz/routes"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "nokib/campwiz/docs"

	"github.com/gin-gonic/gin"
)

// preRun is a function that will be called before the main function
func preRun() {
	database.InitDB()
	cache.InitCacheDB()
}
func postRun() {
}

// @title Campwiz API
// @version 1
// @description This is the API documentation for Campwiz
// @host localhost:8080
// @BasePath /api/v2
// @schemes http https
// @produce json
// @consumes json
// @securitydefinitions.oauth2 implicit
// @type oauth2
// @authorizationurl https://meta.wikimedia.org/w/rest.php/oauth2/authorize
// @tokenurl https://meta.wikimedia.org/w/rest.php/oauth2/access_token

// @license.name LGPL-3.0
// @license.url http://www.gnu.org/licenses/lgpl-3.0.html
func main() {
	// batch_service := services.NewBatchService()
	// batch_service.CreateBatchFromCommonsCategory()
	preRun()
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	routes.NewRoutes(r.Group("/"))
	r.Run()
	postRun()
}
