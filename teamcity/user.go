package teamcity

import (
	"fmt"
	"net/http"

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

// List - Get list of users in range [offset:limit)
func (s *UserService) List(offset, limit int) (*UserList, error) {
	var out UserList
	err := s.restHelper.get("", &out, "Users", buildQueryLocator(
		LocatorStart(offset),
		LocatorCount(limit),
	))
	if err != nil {
		return nil, err
	}
	return &out, err
}

// ListAll returns list of all users
func (s *UserService) ListAll() (*UserList, error) {
	return s.List(0, -1)
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
