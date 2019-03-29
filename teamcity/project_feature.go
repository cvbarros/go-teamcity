package teamcity

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// The ProjectFeature interface represents the different types of features that can be added to a project.
type ProjectFeature interface {
	ID() string
	SetID(value string)

	Type() string

	ProjectID() string
	SetProjectID(value string)

	Disabled() bool
	SetDisabled(value bool)

	Properties() *Properties
}

type projectFeatureJSON struct {
	ID         string      `json:"id,omitempty" xml:"id"`
	Type       string      `json:"type,omitempty" xml:"type"`
	Disabled   *bool       `json:"disabled,omitempty" xml:"disabled"`
	Inherited  *bool       `json:"inherited,omitempty" xml:"inherited"`
	Href       string      `json:"href,omitempty" xml:"href"`
	Properties *Properties `json:"properties,omitempty"`
}

type projectFeatures struct {
	Count int32                `json:"count,omitempty" xml:"count"`
	Href  string               `json:"href,omitempty" xml:"href"`
	Items []projectFeatureJSON `json:"projectFeature"`
}

// ProjectFeatureService provides operations for managing project features.
type ProjectFeatureService struct {
	ProjectID  string
	restHelper *restHelper
}

func newProjectFeatureService(projectID string, c *http.Client, sling *sling.Sling) *ProjectFeatureService {
	return &ProjectFeatureService{
		ProjectID:  projectID,
		restHelper: newRestHelperWithSling(c, sling),
	}
}

// Create a new ProjectFeature under the current project.
func (s *ProjectFeatureService) Create(feature ProjectFeature) (ProjectFeature, error) {
	if feature == nil {
		return nil, errors.New("feature is nil")
	}
	if feature.ProjectID() != s.ProjectID {
		return nil, errors.Errorf("given ProjectFeature for project %q to ProjectFeatureService for project %q.", feature.ProjectID(), s.ProjectID)
	}

	requestBody := &projectFeatureJSON{
		Type:       feature.Type(),
		Disabled:   NewBool(feature.Disabled()),
		Properties: feature.Properties(),
	}
	createdProjectFeature := &projectFeatureJSON{}

	url := fmt.Sprintf("projects/%s/projectFeatures", s.ProjectID)
	if err := s.restHelper.post(url, &requestBody, createdProjectFeature, "projectFeature"); err != nil {
		return nil, err
	}

	return s.parseProjectFeatureJSONResponse(*createdProjectFeature)
}

// Get all project features for the current project.
func (s *ProjectFeatureService) Get() ([]ProjectFeature, error) {
	var out projectFeatures

	url := fmt.Sprintf("projects/%s/projectFeatures", s.ProjectID)
	if err := s.restHelper.get(url, &out, "projectFeature"); err != nil {
		return nil, err
	}

	result := make([]ProjectFeature, len(out.Items))
	for i, featureJSON := range out.Items {
		feature, err := s.parseProjectFeatureJSONResponse(featureJSON)
		if err != nil {
			return result[:i], err
		}

		result[i] = feature
	}

	return result, nil
}

// GetByID returns a single ProjectFeature for the current project by it's id.
func (s *ProjectFeatureService) GetByID(id string) (ProjectFeature, error) {
	var out projectFeatureJSON

	url := fmt.Sprintf("projects/%s/projectFeatures/%s", s.ProjectID, id)
	if err := s.restHelper.get(url, &out, "projectFeature"); err != nil {
		return nil, err
	}

	return s.parseProjectFeatureJSONResponse(out)
}

func (s *ProjectFeatureService) parseProjectFeatureJSONResponse(feature projectFeatureJSON) (ProjectFeature, error) {
	switch feature.Type {
	case "versionedSettings":
		return loadProjectFeatureVersionedSettings(s.ProjectID, feature)
	default:
		return nil, errors.Errorf("Unknown project feature type %q", feature.Type)
	}
}
