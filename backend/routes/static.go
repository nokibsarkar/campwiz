package routes

import (
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Campwiz",
	})
}

// NewStaticRouter creates a new router for static files
func NewStaticRouter(parent *gin.RouterGroup) {
	parent.GET("/", Home)
	parent.GET(("favicon.ico"), func(c *gin.Context) {
		c.File("static/favicon.ico")
	})
}
