package teamcity

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
)

// User is the model for User entities in TeamCity
type User struct {
	Username   string               `json:"username,omitempty" xml:"username"`
	Name       string               `json:"name,omitempty" xml:"name"`
	ID         int                  `json:"id,omitempty" xml:"id"`
	Email      string               `json:"email,omitempty" xml:"email"`
	Properties *Properties          `json:"properties,omitempty" xml:"properties"`
	Roles      *roleAssignmentsJSON `json:"roles,omitempty" xml:"roles"`
	Groups     *groupAssignments    `json:"groups,omitempty" xml:"groups"`
}

type groupAssignments struct {
	Count int     `json:"count,omitempty" xml:"count"`
	Items []Group `json:"group,omitempty" xml:"groups"`
}

// UserList contains list of users
type UserList struct {
	Count int    `json:"count,omitempty" xml:"count"`
	Items []User `json:"user,omitempty" xml:"user"`
}

// NewUser returns an instance of a User. A non-empty Username, Name and Email is required.
func NewUser(username string, name string, email string) (*User, error) {
	if username == "" {
		return nil, fmt.Errorf("Key is required")
	}

	if name == "" {
		return nil, fmt.Errorf("Name is required")
	}

	if email == "" {
		return nil, fmt.Errorf("Email is required")
	}

	return &User{
		Username: username,
		Name:     name,
		Email:    email,
	}, nil
}

// UserService has operations for handling Users
type UserService struct {
	sling      *sling.Sling
	httpClient *http.Client
	restHelper *restHelper
}

func newUserService(base *sling.Sling, httpClient *http.Client) *UserService {
	sling := base.Path("users/")
	return &UserService{
		httpClient: httpClient,
		sling:      sling,
		restHelper: newRestHelperWithSling(httpClient, sling),
	}
}

// Create - Creates a new User
func (s *UserService) Create(user *User) (*User, error) {
	var created User
	err := s.restHelper.post("", user, &created, "User")

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetByID - Get a User by its User ID
func (s *UserService) GetByID(ID int) (*User, error) {
	return s.getByLocator(LocatorID(fmt.Sprint(ID)))
}

// GetByUsername - Get a User by its User Username
func (s *UserService) GetByUsername(username string) (*User, error) {
	return s.getByLocator(LocatorUsername(username))
}

// GetByName - Get a User by its User Name
func (s *UserService) GetByName(name string) (*User, error) {
	return s.getByLocator(LocatorName(name))
}

func (s *UserService) getByLocator(locator Locator) (*User, error) {
	var out User
	err := s.restHelper.get(locator.String(), &out, "User")
	if err != nil {
		return nil, err
	}

	return &out, err
}

// DeleteByID - Deletes a User by its User ID
func (s *UserService) DeleteByID(id int) error {
	return s.deleteByLocator(LocatorID(fmt.Sprint(id)))
}

// DeleteByName - Deletes a User by its User Name
func (s *UserService) DeleteByName(name string) error {
	return s.deleteByLocator(LocatorName(name))
}

// DeleteByUsername - Deletes a User by its User Username
func (s *UserService) DeleteByUsername(username string) error {
	return s.deleteByLocator(LocatorUsername(username))
}

func (s *UserService) deleteByLocator(locator Locator) error {
	err := s.restHelper.delete(locator.String(), "User")
	return err
}

// List - Get list of all User
func (s *UserService) List() (*UserList, error) {
	var out UserList
	err := s.restHelper.get("", &out, "Users")
	if err != nil {
		return nil, err
	}
	return &out, err
}

// GroupAddByID - Add User with userID to Group with groupKey
func (s *UserService) GroupAddByID(userID int, groupKey string) (*Group, error) {
	return s.groupAddByKey(LocatorID(fmt.Sprint(userID)), groupKey)
}

// GroupAddByUsername - Add User with username to Group with groupKey
func (s *UserService) GroupAddByUsername(username, groupKey string) (*Group, error) {
	return s.groupAddByKey(LocatorUsername(username), groupKey)
}

// GroupAddByName - Add User with name to Group with groupKey
func (s *UserService) GroupAddByName(userName, groupKey string) (*Group, error) {
	return s.groupAddByKey(LocatorName(userName), groupKey)
}

func (s *UserService) groupAddByKey(locator Locator, groupKey string) (*Group, error) {
	var out Group

	err := s.restHelper.post(fmt.Sprintf("%s/groups", locator), Group{Key: groupKey}, &out, "User")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GroupDeleteByID - Add User with userID to Group with groupKey
func (s *UserService) GroupDeleteByID(userID int, groupKey string) (*Group, error) {
	return s.groupDeleteByKey(LocatorID(fmt.Sprint(userID)), groupKey)
}

// GroupDeleteByUsername - Add User with username to Group with groupKey
func (s *UserService) GroupDeleteByUsername(username, groupKey string) (*Group, error) {
	return s.groupDeleteByKey(LocatorUsername(username), groupKey)
}

// GroupDeleteByName - Add User with name to Group with groupKey
func (s *UserService) GroupDeleteByName(userName, groupKey string) (*Group, error) {
	return s.groupDeleteByKey(LocatorName(userName), groupKey)
}

func (s *UserService) groupDeleteByKey(locator Locator, groupKey string) (*Group, error) {
	var out Group
	err := s.restHelper.delete(fmt.Sprintf("%s/groups/%s", locator, groupKey), "User")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// IsGroupMemberByID - checks the user's group membership by ID
func (s *UserService) IsGroupMemberByID(id int, key string) (bool, error) {
	return s.isGroupMemberByLocator(LocatorID(fmt.Sprint(id)), key)
}

// IsGroupMemberByUsername - checks the user's group membership by Username
func (s *UserService) IsGroupMemberByUsername(username, key string) (bool, error) {
	return s.isGroupMemberByLocator(LocatorUsername(username), key)
}

// IsGroupMemberByName - checks the user's group membership by Name
func (s *UserService) IsGroupMemberByName(name, key string) (bool, error) {
	return s.isGroupMemberByLocator(LocatorName(name), key)
}

func (s *UserService) isGroupMemberByLocator(locator Locator, key string) (bool, error) {
	var out Group
	err := s.restHelper.get(fmt.Sprintf("%s/groups/%s", locator, LocatorKey(key)), &out, "User")
	if err != nil {
		strErr := err.Error()
		if strings.Contains(strErr, "status code: 404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
