package teamcity

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

//ConditionStrings - All possible condition strings. Do not change the values.
var ConditionStrings = []string{
	"exists",
	"not-exists",
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
	NotExists         string
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
	NotExists:         ConditionStrings[1],
	Equals:            ConditionStrings[2],
	DoesNotEqual:      ConditionStrings[3],
	MoreThan:          ConditionStrings[4],
	NoMoreThan:        ConditionStrings[5],
	LessThan:          ConditionStrings[6],
	NoLessThan:        ConditionStrings[7],
	StartsWith:        ConditionStrings[8],
	Contains:          ConditionStrings[9],
	DoesNotContain:    ConditionStrings[10],
	EndsWith:          ConditionStrings[11],
	Matches:           ConditionStrings[12],
	DoesNotMatch:      ConditionStrings[13],
	VersionMoreThan:   ConditionStrings[14],
	VersionNoMoreThan: ConditionStrings[15],
	VersionLessThan:   ConditionStrings[16],
	VersionNoLessThan: ConditionStrings[17],
}

// AgentRequirement is a condition evaluated per agent to see if a build type is compatible or not
type AgentRequirement struct {
	// build type id
	BuildTypeID string `json:"-"`

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

type agentRequirementsJSON struct {
	Count int32               `json:"count,omitempty" xml:"count"`
	Items []*AgentRequirement `json:"agent-requirement"`
}

// AgentRequirementService provides operations for managing agent requirements for a build type
type AgentRequirementService struct {
	BuildTypeID  string
	httpClient   *http.Client
	base         *sling.Sling
	restHelper   *restHelper
	buildLocator Locator
}

func newAgentRequirementService(buildTypeID string, client *http.Client, base *sling.Sling) *AgentRequirementService {
	buildLocator := Locator(buildTypeID)
	sling := base.Path(fmt.Sprintf("buildTypes/%s/agent-requirements/", buildLocator))

	return &AgentRequirementService{
		BuildTypeID:  buildTypeID,
		httpClient:   client,
		base:         sling,
		restHelper:   newRestHelperWithSling(client, sling),
		buildLocator: buildLocator,
	}
}

//Create a new agent requirement for build type
func (s *AgentRequirementService) Create(req *AgentRequirement) (*AgentRequirement, error) {
	var created AgentRequirement
	_, err := s.base.New().Post("").BodyJSON(req).ReceiveSuccess(&created)

	if err != nil {
		return nil, err
	}

	created.BuildTypeID = s.BuildTypeID
	return &created, nil
}

//GetByID returns an agent requirement by its id
func (s *AgentRequirementService) GetByID(id string) (*AgentRequirement, error) {
	var out AgentRequirement
	resp, err := s.base.New().Get(id).ReceiveSuccess(&out)

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("404 Not Found - Trigger (id: %s) for buildTypeId (id: %s) was not found", id, s.BuildTypeID)
	}

	if err != nil {
		return nil, err
	}
	out.BuildTypeID = s.BuildTypeID
	return &out, nil
}

//GetAll returns all agent requirements for a given build configuration
func (s *AgentRequirementService) GetAll() ([]*AgentRequirement, error) {
	var aux agentRequirementsJSON
	err := s.restHelper.get("", &aux, "agent requirements")
	if err != nil {
		return nil, err
	}
	for _, i := range aux.Items {
		i.BuildTypeID = s.BuildTypeID
	}
	return aux.Items, nil
}

//Delete removes an agent requirement from the build configuration by its id
func (s *AgentRequirementService) Delete(id string) error {
	request, _ := s.base.New().Delete(id).Request()
	response, err := s.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode == 204 {
		return nil
	}

	if response.StatusCode != 200 && response.StatusCode != 204 {
		respData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error '%d' when deleting agent requirement: %s", response.StatusCode, string(respData))
	}

	return nil
}
