package services

import (
	"fmt"
	"nokib/campwiz/consts"
	"nokib/campwiz/database/cache"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthenticationService struct {
	Config *consts.AuthenticationConfiguration
}
type SessionClaims struct {
	Permission consts.PermissionGroup `json:"permission,required"`
	Name       string                 `json:"name"`
	jwt.RegisteredClaims
}

func NewAuthenticationService() *AuthenticationService {
	return &AuthenticationService{
		Config: &consts.Config.Auth,
	}
}
func (a *AuthenticationService) VerifyToken(cache *gorm.DB, tokenMap *SessionClaims) error {
	fmt.Println("Verifying token")
	return nil
}
func (a *AuthenticationService) NewSession(tx *gorm.DB, tokenMap *SessionClaims) (string, error) {
	session := &cache.Session{
		UserID:     tokenMap.Subject,
		Username:   tokenMap.Name,
		Permission: tokenMap.Permission,
		ExpiresAt:  tokenMap.ExpiresAt.Time,
	}
	result := tx.Create(session)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error)
		return "", result.Error
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenMap)
	accessToken, err := token.SignedString([]byte(a.Config.Secret))
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	return accessToken, nil
}
func (a *AuthenticationService) NewRefreshToken(tokenMap *SessionClaims) (string, error) {
	refreshClaims := &SessionClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"campwiz"},
			Subject:   tokenMap.Subject,
			Issuer:    a.Config.Issuer,
			ExpiresAt: jwt.NewNumericDate(tokenMap.ExpiresAt.Time.Add(time.Minute * time.Duration(a.Config.Refresh))),
		},
		Permission: tokenMap.Permission,
		Name:       tokenMap.Name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := token.SignedString([]byte(a.Config.Secret))
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	return refreshToken, nil
}
func (a *AuthenticationService) RefreshSession(cacheDB *gorm.DB, tokenMap *SessionClaims) (accessToken string, err error) {
	fmt.Println("Refreshing session")
	sessionIDString := tokenMap.ID
	if sessionIDString == "" {
		return "", fmt.Errorf("No session ID found")
	}
	sessionID, err := strconv.ParseUint(sessionIDString, 10, 64)
	if err != nil {
		return "", err
	}
	session := &cache.Session{
		ID:         sessionID,
		UserID:     tokenMap.Subject,
		Username:   tokenMap.Name,
		Permission: tokenMap.Permission,
		ExpiresAt:  tokenMap.ExpiresAt.Time,
	}
	tx := cacheDB.Begin()
	result := tx.First(session, &cache.Session{ID: sessionID})
	if result.Error != nil {
		fmt.Println("Error: ", result.Error)
		tx.Rollback()
		return "", result.Error
	}
	session.ExpiresAt = time.Now().UTC().Add(time.Minute * time.Duration(a.Config.Expiry))
	result = tx.Save(session)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error)
		tx.Rollback()
		return "", result.Error
	}

	accessToken, err = a.NewSession(tx, tokenMap)
	if err != nil {
		fmt.Println("Error: ", err)
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return accessToken, nil
}
func (a *AuthenticationService) RemoveSession(cache *gorm.DB, tokenMap *SessionClaims) error {
	fmt.Println("Removing session")
	return nil
}
func (a *AuthenticationService) Init(redirectUri string) string {
	fmt.Println("Initializing OAuth 2.0")
	return ""
}
