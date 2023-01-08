package nasa

import "url-collector/internal/collector"

// Resource is the data type returned by nasa api
type Resource struct {
	Title       string `json:"title"`
	Copyright   string `json:"copyright"`
	URL         string `json:"url"`
	HDURL       string `json:"hdurl"`
	Type        string `json:"media_type"`
	Date        string `json:"date"`
	Explanation string `json:"explanation"`
}

// ToMetadata converts the data to more generic structure
func (m Resource) ToMetadata() collector.Metadata {
	return collector.Metadata{
		Title: m.Title,
		HDURL: m.HDURL,
		URL:   m.URL,
	}
}
