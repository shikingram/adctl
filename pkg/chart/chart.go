package chart

type Chart struct {
	Raw []*File `json:"-"`

	// Metadata is the contents of the Chartfile.
	Metadata *Metadata `json:"metadata"`

	Templates []*File `json:"templates"`
	// Values are default config for this chart.
	Values map[string]interface{} `json:"values"`
	// Schema is an optional JSON schema for imposing structure on Values
	Schema []byte `json:"schema"`
	// Files are miscellaneous files in a chart archive,
	// e.g. README, LICENSE, etc.
	Files []*File `json:"files"`

	InitTemplates []*File `json:"initTemplates"`
}

// Name returns the name of the chart.
func (ch *Chart) Name() string {
	if ch.Metadata == nil {
		return ""
	}
	return ch.Metadata.Name
}

func (ch *Chart) ChartPath() string {
	return ch.Name()
}

// Validate validates the metadata.
func (ch *Chart) Validate() error {
	return ch.Metadata.Validate()
}
