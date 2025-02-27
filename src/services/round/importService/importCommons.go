package importservice

import (
	"maps"
	"nokib/campwiz/database"
)

type CommonsCategoryListSource struct {
	Categories   []string
	currentIndex int
	commons_repo *database.CommonsRepository
}

// ImportImageResults imports images from commons categories
// For Each invocation it will import images from a single category
// If all categories are imported it will return nil
// If there are no images in the category it will return nil
// If there are images in the category it will return the images
// If there are failed images in the category it will return the reason as value of the map
func (c *CommonsCategoryListSource) ImportImageResults(failedImageReason *map[string]string) ([]database.ImageResult, *map[string]string) {
	if c.currentIndex < len(c.Categories) {
		category := c.Categories[c.currentIndex]
		c.currentIndex++
		successMedia, currentfailedImages := c.commons_repo.GetImagesFromCommonsCategories(category)
		maps.Copy(*failedImageReason, currentfailedImages)
		return successMedia, failedImageReason
	}
	return nil, failedImageReason
}
func NewCommonsCategoryListSource(categories []string) *CommonsCategoryListSource {
	return &CommonsCategoryListSource{
		Categories:   categories,
		currentIndex: 0,
		commons_repo: database.NewCommonsRepository(),
	}
}
