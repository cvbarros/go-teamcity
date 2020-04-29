package teamcity

import (
	"encoding/json"
)

//FeatureCommitStatusPublisher represents a golang build feature. Implements BuildFeature interface
type FeatureGolangPublisher struct {
	id          string
	disabled    bool
	buildTypeID string

	properties *Properties
}

// NewFeatureGolang returns a new instance of the FeatureGolangPublisher struct
func NewFeatureGolang() *FeatureGolangPublisher {
	return &FeatureGolangPublisher{
		properties: NewProperties(),
	}
}

//ID returns the ID for this instance.
func (f *FeatureGolangPublisher) ID() string {
	return f.id
}

//SetID sets the ID for this instance.
func (f *FeatureGolangPublisher) SetID(value string) {
	f.id = value
}

//Type returns the "commit-status-publisher", the keyed-type for this build feature instance
func (f *FeatureGolangPublisher) Type() string {
	return "golang"
}

//Disabled returns whether this build feature is disabled or not.
func (f *FeatureGolangPublisher) Disabled() bool {
	return f.disabled
}

//SetDisabled sets whether this build feature is disabled or not.
func (f *FeatureGolangPublisher) SetDisabled(value bool) {
	f.disabled = value
}

//BuildTypeID is a getter for the Build Type ID associated with this build feature.
func (f *FeatureGolangPublisher) BuildTypeID() string {
	return f.buildTypeID
}

//SetBuildTypeID is a setter for the Build Type ID associated with this build feature.
func (f *FeatureGolangPublisher) SetBuildTypeID(value string) {
	f.buildTypeID = value
}

//Properties returns a *Properties instance representing a serializable collection to be used.
func (f *FeatureGolangPublisher) Properties() *Properties {
	return f.properties
}

//MarshalJSON implements JSON serialization for FeatureCommitStatusPublisher
func (f *FeatureGolangPublisher) MarshalJSON() ([]byte, error) {
	out := &buildFeatureJSON{
		ID:         f.id,
		Disabled:   NewBool(f.disabled),
		Properties: f.properties,
		Inherited:  NewFalse(),
		Type:       f.Type(),
	}

	// this is the only value and has to be set to this - no no point making it user configurable
	out.Properties.AddOrReplaceValue("test.format", "json")
	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for FeatureCommitStatusPublisher
func (f *FeatureGolangPublisher) UnmarshalJSON(data []byte) error {
	var aux buildFeatureJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	f.id = aux.ID

	disabled := aux.Disabled
	if disabled == nil {
		disabled = NewFalse()
	}
	f.disabled = *disabled
	f.properties = NewProperties(aux.Properties.Items...)

	return nil
}
