package database

type ExtMetadataValue[Value any] struct {
	Value  Value  `json:"value"`
	Source string `json:"source"`
	Hidden string `json:"hidden"`
}
type ExtMetadata struct {
	ImageDescription    ExtMetadataValue[string] `json:"ImageDescription"`
	Credit              ExtMetadataValue[string] `json:"Credit"`
	Artist              ExtMetadataValue[string] `json:"Artist"`
	LicenseShortName    ExtMetadataValue[string] `json:"LicenseShortName"`
	UsageTerms          ExtMetadataValue[string] `json:"UsageTerms"`
	AttributionRequired ExtMetadataValue[string] `json:"AttributionRequired"`
	Copyrighted         ExtMetadataValue[string] `json:"Copyrighted"`
	License             ExtMetadataValue[string] `json:"License"`
}

func (e *ExtMetadata) GetImageDescription() string {
	return e.ImageDescription.Value
}
func (e *ExtMetadata) GetCredit() string {
	return e.Credit.Value
}
func (e *ExtMetadata) GetArtist() string {
	return e.Artist.Value
}
func (e *ExtMetadata) GetLicense() string {
	license := e.License.Value
	if license == "" {
		license = e.LicenseShortName.Value
	}
	if license == "" {
		license = e.UsageTerms.Value
	}
	return license
}
