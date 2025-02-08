package services

import (
	"fmt"
	"nokib/campwiz/consts"
	"nokib/campwiz/database/cache"
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
func (a *AuthenticationService) VerifyToken(cacheDB *gorm.DB, tokenMap *SessionClaims) (*cache.Session, error) {
	fmt.Println("Verifying token")
	// Check if the token is in the cache
	sessionIDString := tokenMap.ID
	if sessionIDString == "" {
		return nil, fmt.Errorf("No session ID found")
	}
	session := &cache.Session{
		ID:     sessionIDString,
		UserID: tokenMap.Subject,
	}
	result := cacheDB.First(session)
	if result.Error != nil {
		return nil, result.Error
	}
	return session, nil
}
func (a *AuthenticationService) NewSession(tx *gorm.DB, tokenMap *SessionClaims) (string, *cache.Session, error) {
	session := &cache.Session{
		ID:         GenerateID(),
		UserID:     tokenMap.Subject,
		Username:   tokenMap.Name,
		Permission: tokenMap.Permission,
		ExpiresAt:  tokenMap.ExpiresAt.Time,
	}
	result := tx.Create(session)
	if result.Error != nil {
		return "", nil, result.Error
	}
	tokenMap.ID = session.ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenMap)
	accessToken, err := token.SignedString([]byte(a.Config.Secret))
	if err != nil {
		fmt.Println("Error: ", err)
		return "", nil, err
	}
	return accessToken, session, nil
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
func (a *AuthenticationService) RefreshSession(cacheDB *gorm.DB, tokenMap *SessionClaims) (accessToken string, session *cache.Session, err error) {
	fmt.Println("Refreshing session")
	sessionIDString := tokenMap.ID
	if sessionIDString == "" {
		return "", nil, fmt.Errorf("No session ID found")
	}
	session = &cache.Session{
		ID:         sessionIDString,
		UserID:     tokenMap.Subject,
		Username:   tokenMap.Name,
		Permission: tokenMap.Permission,
		ExpiresAt:  tokenMap.ExpiresAt.Time,
	}
	tx := cacheDB.Begin()
	result := tx.First(session, &cache.Session{ID: sessionIDString})
	if result.Error != nil {
		fmt.Println("Error: ", result.Error)
		tx.Rollback()
		return "", nil, result.Error
	}
	session.ExpiresAt = time.Now().UTC().Add(time.Minute * time.Duration(a.Config.Expiry))
	result = tx.Save(session)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error)
		tx.Rollback()
		return "", nil, result.Error
	}

	accessToken, session, err = a.NewSession(tx, tokenMap)
	if err != nil {
		tx.Rollback()
		return "", nil, err
	}
	tx.Commit()
	return accessToken, session, nil
}
func (a *AuthenticationService) RemoveSession(cacheDB *gorm.DB, ID string) error {
	session := &cache.Session{ID: ID}
	result := cacheDB.Delete(session)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (a *AuthenticationService) Logout(session *cache.Session) error {
	conn, close := cache.GetCacheDB()
	defer close()
	// Remove the session
	return a.RemoveSession(conn, session.ID)
}
