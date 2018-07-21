package teamcity

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

// VcsRoot represents a detailed VCS Root entity
type VcsRoot struct {
	// id
	ID string `json:"id,omitempty" xml:"id"`

	// internal Id
	InternalID string `json:"internalId,omitempty" xml:"internalId"`

	// ModificationCheckInterval value in seconds to override the global server setting.
	ModificationCheckInterval int32 `json:"modificationCheckInterval,omitempty" xml:"modificationCheckInterval"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// project
	Project *ProjectReference `json:"project,omitempty"`

	// Properties for the VCS Root. Do not set directly, instead use NewVcsRoot... constructors.
	Properties *Properties `json:"properties,omitempty"`

	// VcsName is the VCS Type used for this VCS Root. See VcsNames for allowed values.
	// Use NewVcsRoot... constructors to avoid setting this directly.
	VcsName string `json:"vcsName,omitempty" xml:"vcsName"`
}

// VcsRootReference represents a subset detail of a VCS Root
type VcsRootReference struct {

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// project Id
	Project *ProjectReference `json:"project,omitempty" xml:"project"`
}

// VcsRootService has operations for handling vcs roots
type VcsRootService struct {
	sling      *sling.Sling
	httpClient *http.Client
}

//NewGitVcsRoot returns a VCS Root instance that connects to Git VCS.
func NewGitVcsRoot(projectID string, name string, opts *GitVcsRootOptions) (*VcsRoot, error) {
	if projectID == "" {
		return nil, errors.New("projectID is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if opts == nil {
		return nil, errors.New("opts is required")
	}
	return &VcsRoot{
		Name: name,
		Project: &ProjectReference{
			ID: projectID,
		},
		VcsName:    string(VcsNames.Git),
		Properties: opts.gitVcsRootProperties(),
	}, nil
}

func newVcsRootService(base *sling.Sling, httpClient *http.Client) *VcsRootService {
	return &VcsRootService{
		sling:      base.Path("vcs-roots/"),
		httpClient: httpClient,
	}
}

// Create creates a new vcs root
func (s *VcsRootService) Create(projectID string, vcsRoot *VcsRoot) (*VcsRootReference, error) {
	var created VcsRootReference

	_, err := s.sling.New().Post("").BodyJSON(vcsRoot).ReceiveSuccess(&created)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetByID Retrieves a vcs root by id using the id: locator
func (s *VcsRootService) GetByID(id string) (*VcsRoot, error) {
	var out VcsRoot

	resp, err := s.sling.New().Get(id).ReceiveSuccess(&out)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error when retrieving VcsRoot id = '%s', status: %d", id, resp.StatusCode)
	}

	return &out, err
}

//Delete a VCS Root resource using id: locator
func (s *VcsRootService) Delete(id string) error {
	request, _ := s.sling.New().Delete(id).Request()

	//TODO: Expose the same httpClient used by sling
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
		return fmt.Errorf("Error '%d' when deleting vcsRoot: %s", response.StatusCode, string(respData))
	}

	return nil
}
