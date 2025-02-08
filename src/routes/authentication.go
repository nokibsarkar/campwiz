package routes

import (
	"errors"
	"fmt"
	"nokib/campwiz/consts"
	"nokib/campwiz/database/cache"
	"nokib/campwiz/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

/*
This is the authentication service. It is used to authenticate users.
It would usually access the cache database to check if the user is authenticated.
It would access the JWT
*/
type AuthenticationMiddleWare struct {
	Config *consts.AuthenticationConfiguration
}

const AuthenticationCookieName = "auth"
const RefreshCookieName = "X-Refresh-Token"

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
	return "", errors.New("No Access token found")
}
func (a *AuthenticationMiddleWare) extractRefreshToken(c *gin.Context) (string, error) {
	cookies := c.Request.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == RefreshCookieName {
			return cookie.Value, nil
		}
	}
	return "", errors.New("No Refresh token found")
}
func (a *AuthenticationMiddleWare) decodeToken(tokenString string) (*services.SessionClaims, error) {
	claims := &services.SessionClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		claims := token.Claims.(jwt.Claims)
		iss, ok := claims.GetIssuer()
		if ok != nil {
			return nil, errors.New("Issuer not found")
		}
		if iss != a.Config.Issuer {
			return nil, errors.New("Invalid issuer")
		}

		return []byte(a.Config.Secret), nil
	})
	if err != nil {
		return claims, err
	}
	if !token.Valid {
		return claims, errors.New("Invalid token")
	}
	return claims, nil
}
func (a *AuthenticationMiddleWare) checkIfUnauthenticatedAllowed(c *gin.Context) bool {
	path := c.Request.URL.Path
	return strings.HasPrefix(path, "/user/callback")
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
			var session *cache.Session
			tokenMap, err := a.decodeToken(token)
			cache_db, close := cache.GetCacheDB()
			defer close()
			auth_service := services.NewAuthenticationService()
			if err != nil {
				if strings.Contains(err.Error(), "token is expired") {
					// Token is expired
					newAccessToken, sess, err := auth_service.RefreshSession(cache_db, tokenMap)
					if err != nil {
						c.Set("error", err)
						c.AbortWithStatusJSON(401, ResponseError{Detail: "Unauthorized : Token expired and could not be refreshed"})
						return
					} else {
						c.SetCookie(AuthenticationCookieName, newAccessToken, a.Config.Expiry, "/", "", false, true)
						session = sess
					}
				} else {
					c.Set("error", err)
					c.AbortWithStatusJSON(401, ResponseError{Detail: "Unauthorized : Token could not be decoded"})
					return
				}
			} else {
				session, err = auth_service.VerifyToken(cache_db, tokenMap)
				if err != nil {
					c.Set("error", err)
					c.AbortWithStatusJSON(401, ResponseError{Detail: "Unauthorized : Invalid token"})
					return
				}
			}
			if session == nil {
				c.AbortWithStatusJSON(401, ResponseError{Detail: "Unauthorized : No session found"})
			}
			c.Set("session", session)
		}
	}
	c.Next()
}
