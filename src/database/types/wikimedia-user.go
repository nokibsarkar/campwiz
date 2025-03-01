package types

import "strings"

// This represents the username provided by wikimedia, it should be normalized already by wikimedia
type WikimediaUsernameType string

const WIKIMEDIA_USERNAME_PREFIX = "u"

func (username *WikimediaUsernameType) GormDataType() string {
	return "varchar(255)"
}

func NewWikimediaUsernameType(username string) (WikimediaUsernameType, error) {
	if !strings.HasPrefix(username, WIKIMEDIA_USERNAME_PREFIX) {
		return "", ErrorType
	}
	wikimediaUsernameType := WikimediaUsernameType(username)
	return wikimediaUsernameType, nil
}
