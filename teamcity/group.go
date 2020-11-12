package teamcity

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// Group is the model for group entities in TeamCity
type Group struct {
	Key         string               `json:"key,omitempty" xml:"key"`
	Description string               `json:"description,omitempty" xml:"description"`
	Name        string               `json:"name,omitempty" xml:"name"`
	Users       *UserList            `json:"users,omitempty" xml:"users"`
	Roles       *roleAssignmentsJSON `json:"roles,omitempty" xml:"roles"`
	Properties  *Properties          `json:"properties,omitempty" xml:"properties"`
}

// GroupList is the model for group list in TeamCity
type GroupList struct {
	Count int     `json:"count,omitempty" xml:"count"`
	Items []Group `json:"group, omitempty" xml:"group"`
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
	return s.getByLocator(LocatorKey(key))
}

// GetByName - Get a group by its group name
func (s *GroupService) GetByName(name string) (*Group, error) {
	return s.getByLocator(LocatorName(name))
}

func (s *GroupService) getByLocator(locator Locator) (*Group, error) {
	var out Group
	err := s.restHelper.get(locator.String(), &out, "group")
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

// List - List of all groups
func (s *GroupService) List() (*GroupList, error) {
	var out GroupList
	err := s.restHelper.get("", &out, "group")
	if err != nil {
		return nil, err
	}
	return &out, nil
}
