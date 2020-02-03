package teamcity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

//BuildFeature is an interface representing different types of build features that can be added to a build type.
type BuildFeature interface {
	ID() string
	SetID(value string)
	Type() string
	VcsRootID() string
	SetVcsRootID(value string)
	Properties() *Properties
	BuildTypeID() string
	SetBuildTypeID(value string)
	Disabled() bool
	SetDisabled(value bool)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

type buildFeatureJSON struct {
	Disabled   *bool       `json:"disabled,omitempty" xml:"disabled"`
	Href       string      `json:"href,omitempty" xml:"href"`
	ID         string      `json:"id,omitempty" xml:"id"`
	Inherited  *bool       `json:"inherited,omitempty" xml:"inherited"`
	Properties *Properties `json:"properties,omitempty"`
	Type       string      `json:"type,omitempty" xml:"type"`
	VcsRootID  string      `json:"vcsRootId,omitempty" xml:"vcsRootId"`
}

// Features is a collection of BuildFeature
type Features struct {
	Count int32              `json:"count,omitempty" xml:"count"`
	Href  string             `json:"href,omitempty" xml:"href"`
	Items []buildFeatureJSON `json:"feature"`
}

//BuildFeatureService provides operations for managing build features for a buildType
type BuildFeatureService struct {
	BuildTypeID string
	httpClient  *http.Client
	base        *sling.Sling
	restHelper  *restHelper
}

func newBuildFeatureService(buildTypeID string, c *http.Client, base *sling.Sling) *BuildFeatureService {
	slingName := base.New().Path(fmt.Sprintf("buildTypes/%s/features/", buildTypeID))
	return &BuildFeatureService{
		BuildTypeID: buildTypeID,
		httpClient:  c,
		base:        slingName,
		restHelper:  newRestHelperWithSling(c, slingName),
	}
}

//Create adds a new build feature to build type
func (s *BuildFeatureService) Create(bf BuildFeature) (BuildFeature, error) {
	if bf == nil {
		return nil, errors.New("bf can't be nil")
	}

	req, err := s.base.New().Post("").BodyJSON(bf).Request()

	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unknown error when adding build feature, statusCode: %d", resp.StatusCode)
	}

	return s.readBuildFeatureResponse(resp)
}

//GetByID returns a build feature by its id
func (s *BuildFeatureService) GetByID(id string) (BuildFeature, error) {
	req, err := s.base.New().Get(id).Request()

	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("404 Not Found - Build feature (id: %s) for buildTypeId (id: %s) was not found", id, s.BuildTypeID)
	}

	return s.readBuildFeatureResponse(resp)
}

// GetBuildFeatures gets all the build features of a BuildType
func (s *BuildFeatureService) GetBuildFeatures() ([]BuildFeature, error) {
	var features Features
	err := s.restHelper.get("", &features, "build features")
	if err != nil {
		return nil, err
	}

	buildFeatures := make([]BuildFeature, features.Count)

	for i := range features.Items {
		dt, err := json.Marshal(features.Items[i])
		if err != nil {
			return nil, err
		}

		cbf := GenericBuildFeature{}
		err = cbf.UnmarshalJSON(dt)
		if err != nil {
			return nil, err
		}
		buildFeatures[i] = &cbf
	}

	return buildFeatures, nil
}

//Delete removes a build feature from the build configuration by its id.
func (s *BuildFeatureService) Delete(id string) error {
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
		return fmt.Errorf("Error '%d' when deleting build feature: %s", response.StatusCode, string(respData))
	}

	return nil
}

// DeleteAll removes all build features of a build configuration
func (s *BuildFeatureService) DeleteAll() error {
	emptyStruct := struct{}{}
	req, err := s.base.New().Put("").BodyJSON(emptyStruct).Request()

	if err != nil {
		return err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return nil
	default:
		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error '%d' when deleting build features: %s", resp.StatusCode, string(respData))
	}
}

func (s *BuildFeatureService) readBuildFeatureResponse(resp *http.Response) (BuildFeature, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload buildFeatureJSON
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil, err
	}

	var out BuildFeature
	switch payload.Type {
	case "commit-status-publisher":
		var csp FeatureCommitStatusPublisher
		if err := csp.UnmarshalJSON(bodyBytes); err != nil {
			return nil, err
		}

		out = &csp
	default:
		var cbf GenericBuildFeature
		if err := cbf.UnmarshalJSON(bodyBytes); err != nil {
			return nil, err
		}
		return out, nil
	}

	out.SetBuildTypeID(s.BuildTypeID)
	return out, nil
}
