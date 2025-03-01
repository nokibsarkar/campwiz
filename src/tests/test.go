package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const URL = "http://localhost:8080/api/v2"

type Client struct {
	cl          *http.Client
	AccessToken string
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.AddCookie(&http.Cookie{
		Name:  "c-auth",
		Value: c.AccessToken,
	})

	req.Header.Set("Content-Type", "application/json")
	return c.cl.Do(req)
}
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
func (c *Client) Post(url string, body map[string]interface{}) (*http.Response, error) {
	b := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(b)
	err := encoder.Encode(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

//go:embed .env
var accessToken string

var cl = &Client{
	cl:          &http.Client{},
	AccessToken: accessToken,
}

func GetHomePage() string {
	resp, err := cl.Get(URL)
	if err != nil {
		log.Println("Error: ", err)
		return ""
	}
	defer resp.Body.Close()
	log.Println("Response: ", resp)
	return URL
}

func CreateCampaign() {
	campaign := map[string]interface{}{
		"allowCreations":         true,
		"allowExpansions":        true,
		"allowJuryToParticipate": true,
		"allowMultipleJudgement": true,
		"allowedMediaTypes": []string{
			"ARTICLE",
		},
		"blacklist": "string",
		"coordinators": []string{
			"Nokib Sarkar",
		},
		"description":                    "string",
		"endDate":                        "2006-01-02T15:04:05Z",
		"image":                          "string",
		"language":                       "am",
		"maximumSubmissionOfSameArticle": 0,
		"minimumAddedBytes":              0,
		"minimumAddedWords":              0,
		"minimumDurationMilliseconds":    0,
		"minimumHeight":                  0,
		"minimumResolution":              0,
		"minimumTotalBytes":              0,
		"minimumTotalWords":              0,
		"minimumWidth":                   0,
		"name":                           "string",
		"organizers": []string{
			"Nokib Sarkar",
		},
		"rules":        "string",
		"secretBallot": true,
		"startDate":    "2006-01-02T15:04:05Z",
	}

	url := URL + "/campaign/"
	resp, err := cl.Post(url, campaign)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	defer resp.Body.Close()
	responseJson := make([]byte, 1024)
	resp.Body.Read(responseJson)
	log.Println("Response: ", string(responseJson))
}
func main() {
	log.Println(GetHomePage())

	for range 1000 {
		go CreateCampaign()
	}
	fmt.Scanln()
}
