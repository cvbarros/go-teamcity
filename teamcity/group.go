package teamcity

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// Group is the model for group entities in TeamCity
type Group struct {
	Key         string `json:"key,omitempty" xml:"key"`
	Description string `json:"description,omitempty" xml:"description"`
	Name        string `json:"name,omitempty" xml:"name"`
}

// NewGroup returns an instance of a Group. A non-empty Key and Name is required.
// Description can be an empty string and will be omitted.
func NewGroup(key string, name string, description string) (*Group, error) {
	if key == "" {
		return nil, fmt.Errorf("Key is required")
	}

	if name == "" {
		return nil, fmt.Errorf("Name is required")
	}

	return &Group{
		Key:         key,
		Name:        name,
		Description: description,
	}, nil
}

// GroupService has operations for handling groups
type GroupService struct {
	sling      *sling.Sling
	httpClient *http.Client
	restHelper *restHelper
}

func newGroupService(base *sling.Sling, httpClient *http.Client) *GroupService {
	sling := base.Path("userGroups/")
	return &GroupService{
		httpClient: httpClient,
		sling:      sling,
		restHelper: newRestHelperWithSling(httpClient, sling),
	}
}

// Create - Creates a new group
func (s *GroupService) Create(group *Group) (*Group, error) {
	var created Group
	err := s.restHelper.post("", group, &created, "group")

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetByKey - Get a group by its group key
func (s *GroupService) GetByKey(key string) (*Group, error) {
	var out Group
	locator := LocatorKey(key).String()
	err := s.restHelper.get(locator, &out, "group")
	if err != nil {
		return nil, err
	}

	return &out, err
}

// Delete - Deletes a group by its group key
func (s *GroupService) Delete(key string) error {
	locator := LocatorKey(key).String()
	err := s.restHelper.delete(locator, "group")
	return err
}
