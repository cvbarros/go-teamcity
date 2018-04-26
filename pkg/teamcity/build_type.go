package teamcity

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
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

	// Parameters for the build configuration. Read-only, only useful when retrieving project details
	Parameters *Properties `json:"parameters,omitempty"`

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

// BuildTypeService has operations for handling build configurations and templates
type BuildTypeService struct {
	sling *sling.Sling
}

func newBuildTypeService(base *sling.Sling) *BuildTypeService {
	return &BuildTypeService{
		sling: base.Path("buildTypes/"),
	}
}

// Create Creates a new build type under a project
func (s *BuildTypeService) Create(projectId string, buildType *BuildType) (*BuildTypeReference, error) {
	var created BuildTypeReference

	_, err := s.sling.New().Post("").BodyJSON(buildType).ReceiveSuccess(&created)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetById Retrieves a build type resource by ID
func (s *BuildTypeService) GetById(id string) (*BuildType, error) {
	var out BuildType

	resp, err := s.sling.New().Get(id).ReceiveSuccess(&out)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error when retrieving BuildType id = '%s', status: %d", id, resp.StatusCode)
	}

	return &out, err
}

//Delete a build type resource
func (s *BuildTypeService) Delete(id string) error {
	request, _ := s.sling.New().Delete(id).Request()

	response, err := http.DefaultClient.Do(request)
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
		return fmt.Errorf("Error '%d' when deleting build type: %s", response.StatusCode, string(respData))
	}

	return nil
}
