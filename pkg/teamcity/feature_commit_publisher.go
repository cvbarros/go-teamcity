package teamcity

type FeatureCommitStatusPublisherOptions interface {
	Properties() *Properties
}

type FeatureCommitStatusPublisher struct {
	ID        string
	Type      string
	VcsRootID string
	Options   FeatureCommitStatusPublisherOptions

	properties *Properties
}

type buildFeatureJson struct {
	// disabled
	Disabled *bool `json:"disabled,omitempty" xml:"disabled"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// properties
	Properties *Properties `json:"properties,omitempty"`

	// type
	Type string `json:"type,omitempty" xml:"type"`
}
