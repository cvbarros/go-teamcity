package teamcity

import (
	"fmt"
)

// BuildType represents a build configuration or a build configuration template
type BuildType struct {

	// agent requirements
	// AgentRequirements *AgentRequirements `json:"agent-requirements,omitempty"`

	// // artifact dependencies
	// ArtifactDependencies *ArtifactDependencies `json:"artifact-dependencies,omitempty"`

	// // branches
	// Branches *Branches `json:"branches,omitempty"`

	// // builds
	// Builds *Builds `json:"builds,omitempty"`

	// // compatible agents
	// CompatibleAgents *Agents `json:"compatibleAgents,omitempty"`

	// description
	Description string `json:"description,omitempty" xml:"description"`

	// features
	// Features *Features `json:"features,omitempty"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// inherited
	// Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// internal Id
	InternalID string `json:"internalId,omitempty" xml:"internalId"`

	// investigations
	// Investigations *Investigations `json:"investigations,omitempty"`

	// links
	// Links *Links `json:"links,omitempty"`

	// locator
	Locator string `json:"locator,omitempty" xml:"locator"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// parameters
	// Parameters *Properties `json:"parameters,omitempty"`

	// paused
	Paused *bool `json:"paused,omitempty" xml:"paused"`

	// project
	// Project *Project `json:"project,omitempty"`

	// project Id
	ProjectID string `json:"projectId,omitempty" xml:"projectId"`

	// project internal Id
	ProjectInternalID string `json:"projectInternalId,omitempty" xml:"projectInternalId"`

	// project name
	ProjectName string `json:"projectName,omitempty" xml:"projectName"`

	// settings
	// Settings *Properties `json:"settings,omitempty"`

	// snapshot dependencies
	// SnapshotDependencies *SnapshotDependencies `json:"snapshot-dependencies,omitempty"`

	// steps
	// Steps *Steps `json:"steps,omitempty"`

	// template
	// Template *BuildType `json:"template,omitempty"`

	// template flag
	TemplateFlag *bool `json:"templateFlag,omitempty" xml:"templateFlag"`

	// templates
	// Templates *BuildTypes `json:"templates,omitempty"`

	// triggers
	// Triggers *Triggers `json:"triggers,omitempty"`

	// type
	Type string `json:"type,omitempty" xml:"type"`

	// uuid
	UUID string `json:"uuid,omitempty" xml:"uuid"`

	// vcs root entries
	// VcsRootEntries *VcsRootEntries `json:"vcs-root-entries,omitempty"`

	// web Url
	WebURL string `json:"webUrl,omitempty" xml:"webUrl"`
}

// BuildTypeReference represents a subset detail of a Build Type
type BuildTypeReference struct {

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// project Id
	ProjectID string `json:"projectId,omitempty" xml:"projectId"`

	// project name
	ProjectName string `json:"projectName,omitempty" xml:"projectName"`
}

// CreateProject Creates a new build type under a project
func (c *Client) CreateBuildType(projectId string, buildType *BuildType) (*BuildTypeReference, error) {
	var created BuildTypeReference

	err := c.doJSONRequest("POST", fmt.Sprintf("projects/id:%s/buildTypes", projectId), buildType, &created)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetBuildType Retrieves a build type resource by ID
func (c *Client) GetBuildType(id string) (*BuildType, error) {
	var out BuildType

	err := c.doJSONRequest("GET", fmt.Sprintf("buildTypes/%s", id), nil, &out)
	if err != nil {
		return nil, err
	}

	return &out, err
}

//DeleteBuildtype Delete a build type resource
func (c *Client) DeleteBuildType(id string) error {
	return c.doJSONRequest("DELETE", fmt.Sprintf("buildTypes/%s", id), nil, nil)
}
