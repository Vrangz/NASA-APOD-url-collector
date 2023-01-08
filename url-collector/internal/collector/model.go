package collector

// MetadataSet is alias for a slice of Metadata type
type MetadataSet []Metadata

// Metadata is a structure to store just necessary resource data
type Metadata struct {
	Title string
	URL   string
	HDURL string
}

// ToURLList converts the metadata to a list of urls ex. []string{"url1", "url2"}
func (set MetadataSet) ToURLList() []string {
	var list = make([]string, 0, len(set))
	for _, m := range set {
		list = append(list, m.URL)
	}
	return list
}
