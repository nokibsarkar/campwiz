package database

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"nokib/campwiz/consts"
	"strings"
	"time"
)

const COMMONS_API = "http://commons.wikimedia.org/w/api.php"

// This Repository would be used to communicate with wikimedia commons
type CommonsRepository struct {
	endpoint    string
	accessToken string
	cl          *http.Client
}
type ImageResult struct {
	ID               uint64 `json:"pageid"`
	Name             string `json:"title"`
	URL              string
	SubmittedAt      time.Time
	UploaderUsername string
	Height           uint64
	Width            uint64
	Size             uint64
	MediaType        string
	Duration         uint64
	License          string
	Description      string
	CreditHTML       string
}
type GContinue struct {
	Gcmcontinue string `json:"gcmcontinue"`
}
type Paginator[PageType any] struct {
	repo *CommonsRepository
}
type WikimediaUser struct {
	Name       string    `json:"name"`
	Registered time.Time `json:"registration"`
	CentralIds *struct {
		CentralAuth uint64 `json:"CentralAuth"`
	} `json:"centralids"`
}
type UserList struct {
	Users map[string]WikimediaUser `json:"users"`
}

type KeyValue struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type ImageInfo struct {
	Info []struct {
		Timestamp      time.Time    `json:"timestamp"`
		User           string       `json:"user"`
		Size           uint64       `json:"size"`
		Width          uint64       `json:"width"`
		Height         uint64       `json:"height"`
		Title          string       `json:"canonicaltitle"`
		URL            string       `json:"url"`
		DescriptionURL string       `json:"descriptionurl"`
		MediaType      string       `json:"mediatype"`
		Duration       float64      `json:"duration"`
		ExtMetadata    *ExtMetadata `json:"extmetadata"`
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	resp, err := c.cl.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// returns images from commons categories
func (c *CommonsRepository) GetImagesFromCommonsCategories(category string) ([]ImageResult, []string) {
	// Get images from commons category
	// Create batch from commons category
	log.Println("Getting images from commons category: ", category)
	paginator := NewPaginator[ImageInfoPage](c)
	params := url.Values{
		"action":    {"query"},
		"format":    {"json"},
		"prop":      {"imageinfo"},
		"generator": {"categorymembers"},
		"gcmtitle":  {category},
		"gcmtype":   {"file"},
		"iiprop":    {"timestamp|user|url|size|mediatype|dimensions|extmetadata|canonicaltitle"},
		"limit":     {"max"},
	}
	images, err := paginator.Query(params)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, nil
	}
	result := []ImageResult{}
	for image := range images {
		// Append images to result
		if image == nil {
			break
		}
		if len(image.Info) == 0 {
			fmt.Println("No image info found. Skipping")
			continue
		}
		info := image.Info[0]
		img := ImageResult{
			ID:               uint64(image.Pageid),
			Name:             image.Title,
			URL:              info.URL,
			UploaderUsername: info.User,
			SubmittedAt:      info.Timestamp,
			Height:           info.Height,
			Width:            info.Width,
			Size:             info.Size,
			MediaType:        info.MediaType,
			Duration:         uint64(info.Duration * 1e3), // Convert to milliseconds

		}
		if info.ExtMetadata != nil {
			img.License = info.ExtMetadata.GetLicense()
			img.Description = info.ExtMetadata.GetImageDescription()
			img.CreditHTML = info.ExtMetadata.GetCredit()
		}
		result = append(result, img)
	}
	log.Println("Found images: ", len(result))
	return result, []string{}
}

// returns images from commons categories
func (c *CommonsRepository) GeUsersFromUsernames(usernames []string) ([]WikimediaUser, error) {
	// Get images from commons category
	// Create batch from commons category
	paginator := NewPaginator[WikimediaUser](c)
	batchSize := 40
	batchCount := len(usernames) / batchSize
	if len(usernames)%batchSize != 0 {
		batchCount++
	}
	result := []WikimediaUser{}
	for i := range batchCount {
		start := i * batchSize
		end := min((i+1)*batchSize, len(usernames))
		batch := usernames[start:end]
		params := url.Values{
			"action":  {"query"},
			"format":  {"json"},
			"list":    {"users"},
			"ususers": {strings.Join(batch, "|")},
			"usprop":  {"centralids|registration"},
			"limit":   {"max"},
		}
		users, err := paginator.UserList(params)
		if err != nil {
			fmt.Println("Error: ", err)
			return nil, nil
		}
		for user := range users {
			// Append images to result
			if user == nil {
				break
			}
			result = append(result, *user)
		}
	}
	return result, nil
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

type BaseQueryResponse[QueryType any, ContinueType map[string]string] struct {
	BatchComplete string        `json:"batchcomplete"`
	Next          *ContinueType `json:"continue"`
	Query         QueryType     `json:"query"`
	Error         *struct {
		Code    string `json:"code"`
		Info    string `json:"info"`
		Details string `json:"details"`
	} `json:"error"`
}
type PageQueryResponse[PageType any] = BaseQueryResponse[struct {
	Normalized []struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"normalized"`
	Pages map[string]PageType `json:"pages"`
}, map[string]string]

type UserListQueryResponse = BaseQueryResponse[struct {
	Users []WikimediaUser `json:"users"`
}, map[string]string]

// NewCommonsRepository returns a new instance of CommonsRepository
func NewCommonsRepository() *CommonsRepository {
	return &CommonsRepository{
		endpoint:    COMMONS_API,
		accessToken: consts.Config.Auth.AccessToken,
		cl:          &http.Client{},
	}
}
