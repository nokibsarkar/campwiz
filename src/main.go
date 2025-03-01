package main

import (
	"net/http"
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

}
func postRun() {
}
func beforeSetupRouter(testing bool) {
	database.InitDB(testing)
	cache.InitCacheDB(testing)
}
func afterSetupRouter(testing bool) {
}
func SetupRouter(testing bool) *gin.Engine {
	beforeSetupRouter(testing)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.StaticFS("/static", http.Dir("static"))
	routes.NewRoutes(r.Group("/"))
	afterSetupRouter(testing)
	return r
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

// @license.name GPL-3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
// @contact.name Nokib Sarkar
// @contact.email nokibsarkar@gmail.com
// @contact.url https://github.com/nokibsarkar
func main() {
	preRun()
	r := SetupRouter(false)
	r.Run()
	postRun()
}
