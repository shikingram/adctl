package chart

type Metadata struct {
	// The name of the chart. Required.
	Name string `json:"name,omitempty"`
	// A SemVer 2 conformant version string of the chart. Required.
	Version string `json:"version,omitempty"`
	// A one-sentence description of the chart
	Description string `json:"description,omitempty"`
}

func (md *Metadata) Validate() error {
	if md == nil {
		return ValidationError("chart.metadata is required")
	}
	return nil
}
