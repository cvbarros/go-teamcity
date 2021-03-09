package teamcity

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
)

// UserGroupMemberShipService has operations for handling UserGroupMemberships
type UserGroupMemberShipService struct {
	// sling           *sling.Sling
	httpClient      *http.Client
	userRestHelper  *restHelper
	groupRestHelper *restHelper
}

func newUserGroupMembershipService(base *sling.Sling, httpClient *http.Client) *UserGroupMemberShipService {
	return &UserGroupMemberShipService{
		httpClient: httpClient,
		// sling:           sling,
		userRestHelper:  newRestHelperWithSling(httpClient, base.New().Path("users/")),
		groupRestHelper: newRestHelperWithSling(httpClient, base.New().Path("userGroups/")),
	}
}

// GroupAddByID - Add User with UserID to Group with groupKey
func (s *UserGroupMemberShipService) GroupAddByID(userID int, groupKey string) (*Group, error) {
	return s.groupAddByKey(LocatorID(fmt.Sprint(userID)), groupKey)
}

// GroupAddByUsername - Add User with Username to Group with groupKey
func (s *UserGroupMemberShipService) GroupAddByUsername(username, groupKey string) (*Group, error) {
	return s.groupAddByKey(LocatorUsername(username), groupKey)
}

// GroupAddByName - Add User with name to Group with groupKey
func (s *UserGroupMemberShipService) GroupAddByName(name, groupKey string) (*Group, error) {
	return s.groupAddByKey(LocatorName(name), groupKey)
}

func (s *UserGroupMemberShipService) groupAddByKey(locator Locator, groupKey string) (*Group, error) {
	var out Group

	err := s.userRestHelper.post(fmt.Sprintf("%s/groups", locator), Group{Key: groupKey}, &out, "UserGroupMemberShip")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GroupDeleteMemberByID - Delete User with UserID from the Group with groupKey
func (s *UserGroupMemberShipService) GroupDeleteMemberByID(userID int, groupKey string) error {
	return s.groupDeleteMemberByKey(LocatorID(fmt.Sprint(userID)), groupKey)
}

// GroupDeleteMemberByUsername - Delete User with username from the Group with groupKey
func (s *UserGroupMemberShipService) GroupDeleteMemberByUsername(username, groupKey string) error {
	return s.groupDeleteMemberByKey(LocatorUsername(username), groupKey)
}

// GroupDeleteMemberByName - Delete User with name from the Group with groupKey
func (s *UserGroupMemberShipService) GroupDeleteMemberByName(name, groupKey string) error {
	return s.groupDeleteMemberByKey(LocatorName(name), groupKey)
}

func (s *UserGroupMemberShipService) groupDeleteMemberByKey(locator Locator, groupKey string) error {
	err := s.userRestHelper.delete(fmt.Sprintf("%s/groups/%s", locator, groupKey), "UserGroupMemberShip")
	if err != nil {
		return err
	}
	return nil
}

// IsGroupMemberByID - checks the UserGroupMembership's group membership by ID
func (s *UserGroupMemberShipService) IsGroupMemberByID(id int, groupKey string) (bool, error) {
	return s.isGroupMemberByLocator(LocatorIDInt(id), groupKey)
}

// IsGroupMemberByUsername - checks the User's group membership by Username
func (s *UserGroupMemberShipService) IsGroupMemberByUsername(username, groupKey string) (bool, error) {
	return s.isGroupMemberByLocator(LocatorUsername(username), groupKey)
}

// IsGroupMemberByName - checks the User's group membership by Name
func (s *UserGroupMemberShipService) IsGroupMemberByName(name, key string) (bool, error) {
	return s.isGroupMemberByLocator(LocatorName(name), key)
}

func (s *UserGroupMemberShipService) isGroupMemberByLocator(locator Locator, groupKey string) (bool, error) {
	var out Group
	err := s.userRestHelper.get(fmt.Sprintf("%s/groups/%s", locator, LocatorKey(groupKey)), &out, "UserGroupMemberShip")
	if err != nil {
		strErr := err.Error()
		if strings.Contains(strErr, "status code: 404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *UserGroupMemberShipService) GetUserGroupsListByUsername(username string, offset, limit int) (*GroupList, error) {
	return s.getUserGroupsList(LocatorUsername(username), offset, limit)
}

func (s *UserGroupMemberShipService) GetUserGroupsListAllByUsername(username string) (*GroupList, error) {
	return s.getUserGroupsList(LocatorUsername(username), 0, -1)
}

func (s *UserGroupMemberShipService) GetUserGroupsListByID(id, offset, limit int) (*GroupList, error) {
	return s.getUserGroupsList(LocatorIDInt(id), offset, limit)
}
func (s *UserGroupMemberShipService) GetUserGroupsListAllByID(id int) (*GroupList, error) {
	return s.getUserGroupsList(LocatorIDInt(id), 0, -1)
}
func (s *UserGroupMemberShipService) getUserGroupsList(locator Locator, offset, limit int) (*GroupList, error) {
	var out GroupList
	err := s.userRestHelper.get(fmt.Sprintf("%s/groups/", locator), &out, "UserGroupMemberShip",
		buildQueryLocator(
			LocatorStart(offset),
			LocatorCount(limit),
		))
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *UserGroupMemberShipService) GetGroupMembersListByName(groupName string, offset, limit int) (*UserList, error) {
	return s.getGroupMembersList(LocatorName(groupName), offset, limit)
}
func (s *UserGroupMemberShipService) GetGroupMembersListAllByName(groupName string) (*UserList, error) {
	return s.getGroupMembersList(LocatorName(groupName), 0, -1)
}
func (s *UserGroupMemberShipService) GetGroupMembersListByKey(groupKey string, offset, limit int) (*UserList, error) {
	return s.getGroupMembersList(LocatorKey(groupKey), offset, limit)
}
func (s *UserGroupMemberShipService) GetGroupMembersListAllByKey(groupKey string) (*UserList, error) {
	return s.getGroupMembersList(LocatorKey(groupKey), 0, -1)
}
func (s *UserGroupMemberShipService) getGroupMembersList(locator Locator, offset, limit int) (*UserList, error) {
	var out struct {
		Users UserList
	}
	err := s.groupRestHelper.get(locator.String(), &out, "UserGroupMemberShip list members",
		buildQueryLocator(
			LocatorStart(offset),
			LocatorCount(limit),
		))
	if err != nil {
		return nil, err
	}
	return &out.Users, nil
}
