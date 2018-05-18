package teamcity

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

//ConditionStrings - All possible condition strings. Do not change the values.
var ConditionStrings = []string{
	"exists",
	"equals",
	"does-not-equal",
	"more-than",
	"no-more-than",
	"less-than",
	"no-less-than",
	"starts-with",
	"contains",
	"does-not-contain",
	"ends-with",
	"matches",
	"does-not-match",
	"ver-more-than",
	"ver-no-more-than",
	"ver-less-than",
	"ver-no-less-than",
}

//Conditions - Possible conditions for requirements. Do not change the values.
var Conditions = struct {
	Exists            string
	Equals            string
	DoesNotEqual      string
	MoreThan          string
	NoMoreThan        string
	LessThan          string
	NoLessThan        string
	StartsWith        string
	Contains          string
	DoesNotContain    string
	EndsWith          string
	Matches           string
	DoesNotMatch      string
	VersionMoreThan   string
	VersionNoMoreThan string
	VersionLessThan   string
	VersionNoLessThan string
}{
	Exists:            ConditionStrings[0],
	Equals:            ConditionStrings[1],
	DoesNotEqual:      ConditionStrings[2],
	MoreThan:          ConditionStrings[3],
	NoMoreThan:        ConditionStrings[4],
	LessThan:          ConditionStrings[5],
	NoLessThan:        ConditionStrings[6],
	StartsWith:        ConditionStrings[7],
	Contains:          ConditionStrings[8],
	DoesNotContain:    ConditionStrings[9],
	EndsWith:          ConditionStrings[10],
	Matches:           ConditionStrings[11],
	DoesNotMatch:      ConditionStrings[12],
	VersionMoreThan:   ConditionStrings[13],
	VersionNoMoreThan: ConditionStrings[14],
	VersionLessThan:   ConditionStrings[15],
	VersionNoLessThan: ConditionStrings[16],
}

// AgentRequirement is a condition evaluated per agent to see if a build type is compatible or not
type AgentRequirement struct {

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// inherited
	Disabled *bool `json:"disabled,omitempty" xml:"disabled"`

	// type
	Condition string `json:"type,omitempty"`

	// Do not use this directly, build this struct via NewAgentRequirement
	Properties *Properties `json:"properties,omitempty"`
}

//Name - Getter for "property-name" field of the requirement
func (a *AgentRequirement) Name() string {
	v, _ := a.Properties.GetOk("property-name")
	return v
}

//Value - Getter for "property-value" field of the requirement
func (a *AgentRequirement) Value() string {
	v, _ := a.Properties.GetOk("property-value")
	return v
}

// NewAgentRequirement creates AgentRequirement structure with correct representation. Use this instead of creating the struct manually.
func NewAgentRequirement(condition string, paramName string, paramValue string) (*AgentRequirement, error) {

	// Sample structure for a requirement
	// The "property-name" and "property-value" properties nested are used as operands for the condition
	// {
	// 	"id": "RQ_17",
	// 	"type": "ver-no-more-than",
	// 	"properties": {
	// 		"count": 2,
	// 		"property": [
	// 			{
	// 				"name": "property-name",
	// 				"value": "r"
	// 			},
	// 			{
	// 				"name": "property-value",
	// 				"value": "a"
	// 			}
	// 		]
	// 	}
	// },

	if condition != Conditions.Exists && paramValue == "" {
		return nil, errors.New("paramValue is required except for 'exists' condition")
	}

	propertyNameProp := &Property{Name: "property-name", Value: paramName}
	props := NewProperties(propertyNameProp)

	// 'exists' uses only "property-name" operand
	if condition != Conditions.Exists {
		propertyValueProp := &Property{Name: "property-value", Value: paramValue}
		props.Add(propertyValueProp)
	}

	return &AgentRequirement{
		Condition:  condition,
		Properties: props,
	}, nil
}

// AgentRequirements is a collection of AgentRequirement
type AgentRequirements struct {

	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// property
	Items []*AgentRequirement `json:"agent-requirement"`
}

// AgentRequirementService provides operations for managing agent requirements for a build type
type AgentRequirementService struct {
	BuildTypeID string
	httpClient  *http.Client
	base        *sling.Sling
}

func newAgentRequirementService(buildTypeID string, c *http.Client, base *sling.Sling) *AgentRequirementService {
	return &AgentRequirementService{
		BuildTypeID: buildTypeID,
		httpClient:  c,
		base:        base.Path(fmt.Sprintf("buildTypes/%s/agent-requirements/", Locator(buildTypeID).String())),
	}
}

//Create a new agent requirement for build type
func (s *AgentRequirementService) Create(req *AgentRequirement) error {
	var created AgentRequirement
	_, err := s.base.New().Post("").BodyJSON(req).ReceiveSuccess(&created)

	if err != nil {
		return err
	}

	return nil
}
