package database

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const COMMONS_API = "http://commons.wikimedia.org/w/api.php"

// This Repository would be used to communicate with wikimedia commons
type CommonsRepository struct {
	endpoint    string
	accessToken string
	cl          *http.Client
}
type Image struct {
	ID   uint64 `json:"pageid"`
	Name string `json:"title"`
}
type GContinue struct {
	Gcmcontinue string `json:"gcmcontinue"`
}
type Paginator[PageType any] struct {
	repo *CommonsRepository
}

/*
"timestamp": "2024-02-18T18:54:46Z",

	"user": "Sakil302",
	"userid": 7229062,
	"size": 6294651,
	"width": 3120,
	"height": 3900,
	"canonicaltitle": "File:Farmer's joy.jpg",
	"url": "https://upload.wikimedia.org/wikipedia/commons/8/87/Farmer%27s_joy.jpg",
	"descriptionurl": "https://commons.wikimedia.org/wiki/File:Farmer%27s_joy.jpg",
	"descriptionshorturl": "https://commons.wikimedia.org/w/index.php?curid=145519602",
	"metadata": [
	    {
	        "name": "ImageDescription",
	        "value": null
	    },
*/
type KeyValue struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}
type ImageInfo struct {
	Info []struct {
		Timestamp           time.Time  `json:"timestamp"`
		User                string     `json:"user"`
		UserID              uint64     `json:"userid"`
		Size                uint64     `json:"size"`
		Width               uint64     `json:"width"`
		Height              uint64     `json:"height"`
		Title               string     `json:"canonicaltitle"`
		URL                 string     `json:"url"`
		DescriptionURL      string     `json:"descriptionurl"`
		DescriptionShortURL string     `json:"descriptionshorturl"`
		Metadata            []KeyValue `json:"metadata"`
	} `json:"imageinfo"`
}
type Page struct {
	Pageid int    `json:"pageid"`
	Ns     int    `json:"ns"`
	Title  string `json:"title"`
}
type ImageInfoPage struct {
	Page
	ImageInfo
}

func (c *CommonsRepository) Get(values url.Values) (_ io.ReadCloser, err error) {
	// Get values from commons
	url := fmt.Sprintf("%s?%s", c.endpoint, values.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.cl.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// returns images from commons categories
func (c *CommonsRepository) GetImagesFromCommonsCategories(categories []string) []Image {
	// Get images from commons category
	// Create batch from commons category
	paginator := NewPaginator[ImageInfoPage](c)
	params := url.Values{
		"action":    {"query"},
		"format":    {"json"},
		"prop":      {"imageinfo"},
		"generator": {"categorymembers"},
		"gcmtitle":  {"Category:Bangladesh"},
		"gcmtype":   {"file"},
		"iiprop":    {"timestamp|user|url|size|userid|mediatype|metadata|extmetadata|dimensions|commonmetadata|canonicaltitle"},
		"limit":     {"max"},
	}
	images, err := paginator.Query(params)
	if err != nil {
		fmt.Println("Error: ", err)
		return []Image{}
	}
	result := []Image{}
	for image := range images {
		// Append images to result
		result = append(result, Image{
			ID:   uint64(image.Pageid),
			Name: image.Title,
		})
	}
	return result
}
func (c *CommonsRepository) GetImageDetails() {
	// Get image details
}
func (c *CommonsRepository) GetImageMetadata() {
	// Get image metadata
}
func (c *CommonsRepository) GetImageCategories() {
	// Get image categories
}
func (c *CommonsRepository) GetCsrfToken() {
	// Get csrf token
}
func (c *CommonsRepository) GetEditToken() {
	// Get edit token
}
func (c *CommonsRepository) GetUserInfo() {
	// Get user info
}

type QueryResponse[PageType any, ContinueType map[string]string] struct {
	BatchComplete string        `json:"batchcomplete"`
	Next          *ContinueType `json:"continue"`
	Query         struct {
		Normalized []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"normalized"`
		Pages map[string]PageType `json:"pages"`
	} `json:"query"`
}

// NewCommonsRepository returns a new instance of CommonsRepository
func NewCommonsRepository(accessToken string) *CommonsRepository {
	return &CommonsRepository{
		endpoint:    COMMONS_API,
		accessToken: accessToken,
		cl:          &http.Client{},
	}
}
