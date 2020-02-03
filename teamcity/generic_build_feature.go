package teamcity

import (
	"encoding/json"
)

type GenericBuildFeature struct {
	id          string
	featureType string
	buildTypeID string
	disabled    bool
	properties  *Properties
}

func (bf *GenericBuildFeature) ID() string {
	return bf.id
}

func (bf *GenericBuildFeature) SetID(value string) {
	bf.id = value
}

func (bf *GenericBuildFeature) Type() string {
	return bf.featureType
}

func (bf *GenericBuildFeature) VcsRootID() string {
	return ""
}

func (bf *GenericBuildFeature) SetVcsRootID(value string) {
}

func (bf *GenericBuildFeature) Properties() *Properties {
	return bf.properties
}

func (bf *GenericBuildFeature) BuildTypeID() string {
	return bf.buildTypeID
}

func (bf *GenericBuildFeature) SetBuildTypeID(value string) {
	bf.buildTypeID = value
}

func (bf *GenericBuildFeature) Disabled() bool {
	return bf.disabled
}

func (bf *GenericBuildFeature) SetDisabled(value bool) {
	bf.disabled = value
}

func (bf *GenericBuildFeature) MarshalJSON() ([]byte, error) {
	out := &buildFeatureJSON{
		ID:         bf.id,
		Disabled:   NewBool(bf.disabled),
		Properties: bf.properties,
		Inherited:  NewFalse(),
		Type:       bf.Type(),
	}

	return json.Marshal(out)
}

func (bf *GenericBuildFeature) UnmarshalJSON(data []byte) error {
	var aux buildFeatureJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	bf.id = aux.ID
	bf.featureType = aux.Type

	disabled := aux.Disabled
	if disabled == nil {
		disabled = NewFalse()
	}
	bf.disabled = *disabled

	if aux.Properties != nil {
		bf.properties = NewProperties(aux.Properties.Items...)
	}

	return nil
}

func NewGenericBuildFeature(featureType string, propertiesRaw map[string]interface{}) (*GenericBuildFeature, error) {
	properties := NewPropertiesEmpty()
	for name, value := range propertiesRaw {
		value := value.(string)
		properties.Add(&Property{
			Name:  name,
			Value: value,
		})
	}

	return &GenericBuildFeature{
		featureType: featureType,
		properties:  properties,
	}, nil
}
