package database

import (
	"database/sql/driver"
	"errors"
	"strings"

	"slices"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type MediaType string
type MediaTypeSet []MediaType

const (
	MediaTypeArticle MediaType = "ARTICLE"
	MediaTypeImage   MediaType = "BITMAP"
	MediaTypeAudio   MediaType = "AUDIO"
	MediaTypeVideo   MediaType = "VIDEO"
	MediaTypePDF     MediaType = "PDF"
)

// Scan implements the sql.Scanner interface
func (m *MediaTypeSet) Scan(value any) error {
	var mediatypesetSlice []MediaType
	if value == nil {
		*m = nil
		return nil
	}
	// check if already a slice of MediaType
	mediatypesetSlice, ok := value.([]MediaType)
	if !ok {
		// check if an array of strings
		mediaTypeStringSet, ok := value.([]string)
		if !ok {
			// check if a single string
			mediaTypeString, ok := value.(string)
			if !ok {
				return errors.New("invalid media type")
			}
			mediaTypeStringSet = strings.Split(mediaTypeString, ",")
		}
		for _, mediaTypeString := range mediaTypeStringSet {
			mediatypesetSlice = append(mediatypesetSlice, MediaType(mediaTypeString))
		}
	}
	*m = mediatypesetSlice
	return nil
}

// Value implements the driver.Valuer interface
func (m MediaTypeSet) Value() (driver.Value, error) {
	var mediatypesetString []string
	for _, mt := range m {
		mediatypesetString = append(mediatypesetString, string(mt))
	}
	return strings.Join(mediatypesetString, ","), nil
}

// GormDataType implements the gorm.Dialector interface
func (m MediaTypeSet) GormDataType() string {
	return "text"
}

// GormDBDataType implements the gorm.Dialector interface
func (m MediaTypeSet) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "text"
}
func (m MediaTypeSet) Contains(mt MediaType) bool {
	return slices.Contains(m, mt)
}
func (m *MediaTypeSet) Add(mt MediaType) {
	if !m.Contains(mt) {
		*m = append(*m, mt)
	}
}
func (m *MediaTypeSet) Remove(mt MediaType) {
	for i, mediaType := range *m {
		if mediaType == mt {
			*m = append((*m)[:i], (*m)[i+1:]...)
			break
		}
	}
}
