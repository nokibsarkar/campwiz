package main

import (
	"nokib/campwiz/database"
	"nokib/campwiz/database/cache"
	"nokib/campwiz/routes"

	"github.com/gin-gonic/gin"
)

func preRun() {
	database.InitDB()
	cache.InitCacheDB()
}
func postRun() {
}
func main() {
	preRun()
	r := gin.Default()
	routes.NewRoutes(r.Group("/"))
	r.Run()
	postRun()
}
