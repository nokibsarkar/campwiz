package routes

import (
	"errors"
	"fmt"
	"log"
	"nokib/campwiz/consts"
	"nokib/campwiz/services"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
This is the authentication service. It is used to authenticate users.
It would usually access the cache database to check if the user is authenticated.
It would access the JWT
*/
type AuthenticationMiddleWare struct {
	Config *consts.AuthenticationConfiguration
}

const AuthenticationCookieName = "c-auth"
const RefreshCookieName = "X-Refresh-Token"
const SESSION_KEY = "session"

func NewAuthenticationService() *AuthenticationMiddleWare {
	return &AuthenticationMiddleWare{
		Config: &consts.Config.Auth,
	}
}

// This function extracts the access token from the cookies or headers
func (a *AuthenticationMiddleWare) extractAccessToken(c *gin.Context) (string, error) {
	token, _ := c.Cookie(AuthenticationCookieName)
	if token != "" {
		return token, nil
	}
	// Check if the token is in the headers
	token = c.GetHeader("Authorization")
	if token != "" {
		if token[:7] == "Bearer " {
			return token[7:], nil
		}
	}
	return "", errors.New("no-access-token")
}

func (a *AuthenticationMiddleWare) checkIfUnauthenticatedAllowed(c *gin.Context) bool {
	path := c.Request.URL.Path
	// return true
	var UnRestrictedPaths = []string{
		"/user/login",
		"/user/callback",
		// "/api/v2/campaign/",
	}
	for _, p := range UnRestrictedPaths {
		if strings.HasSuffix(p, "*") {
			rest := p[:len(p)-1]
			if strings.HasPrefix(path, rest) {
				return true
			}
		} else {
			if p == path {
				fmt.Println("Unrestricted path")
				return true
			}
		}
	}
	return false
}

/*
This is the authenticator middleware. It is used to authenticate users.
It would usually access the cache database to check if the user is authenticated.
It would access the JWT
*/
func (a *AuthenticationMiddleWare) Authenticate(c *gin.Context) {
	if !a.checkIfUnauthenticatedAllowed(c) {
		token, err := a.extractAccessToken(c)
		if err != nil {
			fmt.Println("Error", err)
			c.Set("error", err)
			c.AbortWithStatusJSON(401, ResponseError{Detail: "Unauthorized : No token found"})
			return
		} else {
			auth_service := services.NewAuthenticationService()
			accessToken, session, err, setCookie := auth_service.Authenticate(token)
			if err != nil {
				fmt.Println("Error", err)
				c.Set("error", err)
				c.AbortWithStatusJSON(401, ResponseError{Detail: "Unauthorized : Invalid token"})
				return
			} else {
				if setCookie {
					log.Println("Setting Authentication Cookie")
					c.SetCookie(AuthenticationCookieName, accessToken, consts.Config.Auth.Expiry, "/", "", false, true)
				}
				c.Set(SESSION_KEY, session)
			}
		}
		c.Next()
	}
}
