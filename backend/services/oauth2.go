package services

import (
	"context"
	"encoding/json"
	"nokib/campwiz/consts"
	"nokib/campwiz/database"
	"time"

	"golang.org/x/oauth2"
)

const META_OAUTH_AUTHORIZE_URL = "https://meta.wikimedia.org/w/rest.php/oauth2/authorize"
const META_OAUTH_ACCESS_TOKEN_URL = "https://meta.wikimedia.org/w/rest.php/oauth2/access_token"
const META_PROFILE_URL = "https://meta.wikimedia.org/w/rest.php/oauth2/resource/profile"
const DATETIMEFORMAT = "20060102150405"

var OAuth2Config = oauth2.Config{
	ClientID:     consts.Config.Auth.OAuth2.ClientID,
	ClientSecret: consts.Config.Auth.OAuth2.ClientSecret,
	RedirectURL:  consts.Config.Server.BaseURL + consts.Config.Auth.OAuth2.RedirectPath,
	Endpoint: oauth2.Endpoint{
		AuthURL:  META_OAUTH_AUTHORIZE_URL,
		TokenURL: META_OAUTH_ACCESS_TOKEN_URL,
	},
}

type OAuth2Service struct {
	Config *consts.OAuth2Configuration
}

func NewOAuth2Service() *OAuth2Service {
	return &OAuth2Service{
		Config: &consts.Config.Auth.OAuth2,
	}
}
func (o *OAuth2Service) Init(callback string) string {
	// state := url.QueryEscape(callback)
	return OAuth2Config.AuthCodeURL(callback)
}
func (o *OAuth2Service) GetToken(code string) (*oauth2.Token, error) {
	token, err := OAuth2Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

/*
sub (central user id)
username
editcount
confirmed_email
blocked
registered
groups
rights
*/
type WikipediaProfileBasic struct {
	CentralID string            `json:"sub"`
	Name      database.UserName `json:"username"`
	Rights    []string          `json:"rights"`
	Blocked   bool              `json:"blocked"`
	Groups    []string          `json:"groups"`
}
type WikipediaProfile struct {
	WikipediaProfileBasic
	Registered string `json:"registered"`
}
type WikipediaProfileFull struct {
	WikipediaProfileBasic
	Registered time.Time `json:"registered"`
}

func (o *OAuth2Service) GetUser(token *oauth2.Token) (*WikipediaProfileFull, error) {
	client := OAuth2Config.Client(context.Background(), token)
	resp, err := client.Get(META_PROFILE_URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	user := &WikipediaProfile{}
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, err
	}
	registered, err := time.Parse(DATETIMEFORMAT, user.Registered)
	if err != nil {
		return nil, err
	}

	return &WikipediaProfileFull{
		WikipediaProfileBasic: user.WikipediaProfileBasic,
		Registered:            registered,
	}, nil
}
