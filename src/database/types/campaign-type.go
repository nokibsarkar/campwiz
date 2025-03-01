package types

import (
	"strings"
)

type CampaignIDType string

const CAMPAIGN_ID_PREFIX = "c"

func (id *CampaignIDType) GormDataType() string {
	return "varchar(255)"
}

func NewCampaignIDType(id string) (CampaignIDType, error) {
	if !strings.HasPrefix(id, CAMPAIGN_ID_PREFIX) {
		return "", ErrorType
	}
	campaignIDType := CampaignIDType(id)
	return campaignIDType, nil
}
