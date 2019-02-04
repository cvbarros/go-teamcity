package teamcity

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type paramType = string

const (
	configParamType = "configuration"
	systemParamType = "system"
	envVarParamType = "env"
)

// ParameterTypes represent the possible parameter types
var ParameterTypes = struct {
	Configuration       paramType
	System              paramType
	EnvironmentVariable paramType
}{
	Configuration:       configParamType,
	System:              systemParamType,
	EnvironmentVariable: envVarParamType,
}

//Parameters is a strongly-typed collection of "Parameter" suitable for serialization
type Parameters struct {
	Count int32        `json:"count,omitempty" xml:"count"`
	Href  string       `json:"href,omitempty" xml:"href"`
	Items []*Parameter `json:"property,omitempty"`
}

//Parameter represents a project or build configuration parameter that may be defined as "configuration", "system" or "environment variable"
type Parameter struct {
	Inherited bool `json:"inherited,omitempty" xml:"inherited"`

	Name string `json:"name,omitempty" xml:"name"`

	Value string `json:"value" xml:"value"`

	Type string `json:"-"`

	Label string

	Description string

	Display string //normal, hidden, prompt

	ReadOnly string

	ControlType string //checkbox, password, text, select
}

//NewParametersEmpty returns an empty collection of Parameters
func NewParametersEmpty() *Parameters {
	return &Parameters{
		Count: 0,
		Items: make([]*Parameter, 0),
	}
}

// NewParameters returns an instance of Parameters collection with the given parameters slice
func NewParameters(items ...*Parameter) *Parameters {
	count := len(items)
	return &Parameters{
		Count: int32(count),
		Items: items,
	}
}

//NewParameter creates a new instance of a parameter with the given type
func NewParameter(t string, name string, value string) (*Parameter, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if t != ParameterTypes.Configuration && t != ParameterTypes.EnvironmentVariable && t != ParameterTypes.System {
		return nil, fmt.Errorf("invalid parameter type, use one of the values defined in ParameterTypes")
	}

	return &Parameter{
		Type:  string(t),
		Name:  name,
		Value: value,
	}, nil
}

//MarshalJSON implements JSON serialization for Parameter
func (p *Parameter) MarshalJSON() ([]byte, error) {
	out := p.Property()

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for Parameter
func (p *Parameter) UnmarshalJSON(data []byte) error {
	var aux Property
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var name, paramType string
	if strings.HasPrefix(aux.Name, "system.") {
		name = strings.TrimPrefix(aux.Name, "system.")
		paramType = ParameterTypes.System
	} else if strings.HasPrefix(aux.Name, "env.") {
		name = strings.TrimPrefix(aux.Name, "env.")
		paramType = ParameterTypes.EnvironmentVariable
	} else {
		name = aux.Name
		paramType = ParameterTypes.Configuration
	}
	p.Name = name
	if aux.Inherited != nil {
		p.Inherited = *aux.Inherited
	}
	p.Value = aux.Value

	if aux.Type != nil {
		p.UnmarshalType(aux.Type.RawValue)
	}

	p.Type = paramType
	return nil
}

// MarshalType performs the special marshalling of the parameter.type.rawValue field
func (p *Parameter) MarshalType() (string, error) {
	s := []string{}

	if p.ControlType != "" {
		s = append(s, p.ControlType)
	}
	if p.Display != "" {
		s = append(s, fmt.Sprintf("display='%s'", p.Display))
	}
	if p.Description != "" {
		s = append(s, fmt.Sprintf("description='%s'", p.Description))
	}
	if p.ReadOnly != "" {
		s = append(s, fmt.Sprintf("readOnly='%s'", p.ReadOnly))
	}
	if p.Label != "" {
		s = append(s, fmt.Sprintf("label='%s'", p.Label))
	}

	return strings.Join(s, " "), nil
}

// UnmarshalType performs the special unmarshalling of the parameter.type.rawValue field
func (p *Parameter) UnmarshalType(t string) error {
	e := regexp.MustCompile("^(checkbox|password|text|select)")
	match := e.FindStringSubmatch(t)
	if len(match) > 0 {
		p.ControlType = match[1]
	}

	e = regexp.MustCompile("display='(nomal|hidden|prompt)'")
	match = e.FindStringSubmatch(t)
	if len(match) > 0 {
		p.Display = match[1]
	}

	e = regexp.MustCompile("readOnly='(true|false)'")
	match = e.FindStringSubmatch(t)
	if len(match) > 0 {
		p.ReadOnly = match[1]
	}

	e = regexp.MustCompile("label='([^']+)'")
	match = e.FindStringSubmatch(t)
	if len(match) > 0 {
		p.Label = match[1]
	}

	e = regexp.MustCompile("description='([^']+)'")
	match = e.FindStringSubmatch(t)
	if len(match) > 0 {
		p.Description = match[1]
	}

	return nil
}

//Properties convert a Parameters collection to a Properties collection
func (p *Parameters) Properties() *Properties {
	out := NewPropertiesEmpty()
	for _, i := range p.Items {
		out.AddOrReplaceProperty(i.Property())
	}
	return out
}

//Property converts a Parameter instance to a Property
func (p *Parameter) Property() *Property {
	out := &Property{
		Name:  fmt.Sprintf("%s%s", paramPrefixByType[p.Type], p.Name),
		Value: p.Value,
	}
	//Omit default inherited value
	if p.Inherited {
		out.Inherited = NewBool(p.Inherited)
	}

	rawValue, err := p.MarshalType()
	if err == nil && rawValue != "" {
		out.Type = &Type{
			RawValue: rawValue,
		}
	}
	return out
}

// AddOrReplaceValue will update a parameter value if it exists, or add if it doesnt
func (p *Parameters) AddOrReplaceValue(t string, n string, v string) {
	for _, elem := range p.Items {
		if elem == nil {
			continue
		}

		if elem.Name == n {
			elem.Value = v
			return
		}
	}
	param, _ := NewParameter(t, n, v)
	p.Add(param)
}

// AddOrReplaceParameter will update a parameter value if another parameter with the same name exists. It won't replace the Parameter struct within the Parameters collection.
func (p *Parameters) AddOrReplaceParameter(param *Parameter) {
	p.AddOrReplaceValue(param.Type, param.Name, param.Value)
}

// Add a new parameter to this collection
func (p *Parameters) Add(param *Parameter) {
	p.Count++
	p.Items = append(p.Items, param)
}

// Concat appends the source Parameters collection to this collection and returns the appended collection
func (p *Parameters) Concat(source *Parameters) *Parameters {
	for _, item := range source.Items {
		p.AddOrReplaceParameter(item)
	}
	return p
}

//Remove a parameter if it exists in the collection
func (p *Parameters) Remove(t string, n string) {
	removed := -1
	for i := range p.Items {
		if p.Items[i].Name == n && p.Items[i].Type == t {
			removed = i
			break
		}
	}
	if removed >= 0 {
		p.Count--
		p.Items = append(p.Items[:removed], p.Items[removed+1:]...)
	}
}

//NonInherited returns a new Parameters collection filtering out all inherited parameters
func (p *Parameters) NonInherited() (po *Parameters) {
	po = NewParametersEmpty()
	for _, c := range p.Items {
		if !c.Inherited {
			po.AddOrReplaceParameter(c)
		}
	}
	return po
}

//GetOk returns a Parameter by it's type/name combination
func (p *Parameters) GetOk(t string, n string) (out *Parameter, ok bool) {
	for i := range p.Items {
		if p.Items[i].Name == n && p.Items[i].Type == t {
			out, ok = p.Items[i], true
			return
		}
	}
	return nil, false
}

var paramPrefixByType = map[string]string{
	string(ParameterTypes.Configuration):       "",
	string(ParameterTypes.System):              "system.",
	string(ParameterTypes.EnvironmentVariable): "env.",
}
