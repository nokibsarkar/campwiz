package routes

import (
	"github.com/gin-gonic/gin"
)

func NewRoutes(parent *gin.RouterGroup) {

	NewUserAuthenticationRoutes(parent)
	r := parent.Group("/api/v2")
	authenticatorService := NewAuthenticationService()
	r.Use(authenticatorService.Authenticate)
	NewCampaignRoutes(r)
	NewSubmissionRoutes(r)
}
