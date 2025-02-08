package routes

import "github.com/gin-gonic/gin"

func ListUsers(c *gin.Context) {
	// ...
}
func GetMe(c *gin.Context) {
	// ...
}
func GetUser(c *gin.Context) {
	// ...
}
func UpdateUser(c *gin.Context) {
	// ...
}
func GetTranslationPath(c *gin.Context) {
	// ...
}
func GetTranslation(c *gin.Context) {
	// ...
}
func UpdateTranslation(c *gin.Context) {
	// ...
}

func NewUserRoutes(parent *gin.RouterGroup) {
	r := parent.Group("/user")
	r.GET("/", ListUsers)
	r.GET("/me", GetMe)
	r.GET("/:id", GetUser)
	r.POST("/:id", UpdateUser)
	r.GET("/translation/:language", GetTranslation)
	r.POST("/translation/:lang", UpdateTranslation)
}
