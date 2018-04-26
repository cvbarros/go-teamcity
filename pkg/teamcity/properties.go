package teamcity

// Properties represents a collection of key/value properties for a resource
type Properties struct {

	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// property
	Items []*Property `json:"property"`
}

// NewProperties returns an instance of Properties collection
func NewProperties(items ...*Property) *Properties {
	count := len(items)
	return &Properties{
		Count: int32(count),
		Items: items,
	}
}

// Property represents a key/value/type structure used by several resources to extend their representation
type Property struct {

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// type
	Type *Type `json:"type,omitempty"`

	// value
	Value string `json:"value,omitempty" xml:"value"`
}

// Type represents a parameter type . The rawValue is the parameter specification as defined in the UI.
type Type struct {
	// raw value
	RawValue string `json:"rawValue,omitempty" xml:"rawValue"`
}
