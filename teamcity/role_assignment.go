package teamcity

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// GroupRoleAssignment is the model for role assignment for groups in TeamCity
type GroupRoleAssignment struct {
	GroupKey string //`json:"groupkey,omitempty" xml:"groupkey"`
	RoleID   string //`json:"roleid,omitempty" xml:"roleid"`
	Scope    string //`json:"scope,omitempty" xml:"scope"`
}

// RoleAssignmentReference represents a response of a request to assign role to a group or a user
type RoleAssignmentReference struct {
	RoleID string `json:"roleId,omitempty" xml:"roleId"`
	Scope  string `json:"scope,omitempty" xml:"scope"`
	Href   string `json:"href,omitempty" xml:"href"`
}

type roleAssignmentsJSON struct {
	Items []RoleAssignmentReference `json:"role"`
}

// NewGroupRoleAssignment returns an instance of a GroupRoleAssignment. A non-empty groupKey, roleId, and scope is required.
func NewGroupRoleAssignment(groupKey string, roleID string, scope string) (*GroupRoleAssignment, error) {
	if groupKey == "" {
		return nil, fmt.Errorf("GroupKey is required")
	}

	if roleID == "" {
		return nil, fmt.Errorf("RoleId is required")
	}

	if scope == "" {
		return nil, fmt.Errorf("scope is required. Use \"g\" at the global level for System Administrators, otherwise for other roles, use \"p:_Root\" for the root project, or \"p:<project_id>\" for other projects")
	}

	return &GroupRoleAssignment{
		GroupKey: groupKey,
		RoleID:   roleID,
		Scope:    scope,
	}, nil
}

// RoleAssignmentService has operations for handling role assignments for groups or users
type RoleAssignmentService struct {
	groupSling  *sling.Sling
	httpClient  *http.Client
	groupHelper *restHelper
}

func newRoleAssignmentService(base *sling.Sling, httpClient *http.Client) *RoleAssignmentService {
	groupSling := base.New().Path(fmt.Sprintf("userGroups/"))
	return &RoleAssignmentService{
		httpClient:  httpClient,
		groupSling:  groupSling,
		groupHelper: newRestHelperWithSling(httpClient, groupSling),
	}
}

// AssignToGroup adds a role assignment to a group
func (s *RoleAssignmentService) AssignToGroup(assignment *GroupRoleAssignment) (*RoleAssignmentReference, error) {
	var out RoleAssignmentReference

	// URL for assigning role is /app/rest/userGroups/{groupLocator}/roles/{roleId}/{scope}
	err := s.groupHelper.post(fmt.Sprintf("%s/roles/%s/%s", assignment.GroupKey, assignment.RoleID, assignment.Scope), nil, &out, "AssignToGroup role to group")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetForGroup get a specific role assignment for a group
func (s *RoleAssignmentService) GetForGroup(assignment *GroupRoleAssignment) (*RoleAssignmentReference, error) {
	var out RoleAssignmentReference

	// URL for getting a specific role assignments is /app/rest/userGroups/{groupLocator}/roles/{roleId}/{scope}
	err := s.groupHelper.get(fmt.Sprintf("%s/roles/%s/%s", assignment.GroupKey, assignment.RoleID, assignment.Scope), &out, "GetForGroup role assignmens for group")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAllForGroup gets all the role assignments for a group
func (s *RoleAssignmentService) GetAllForGroup(group *Group) ([]RoleAssignmentReference, error) {
	var aux roleAssignmentsJSON

	// URL for getting role assignments is /app/rest/userGroups/{groupLocator}/roles
	err := s.groupHelper.get(fmt.Sprintf("%s/roles", group.Key), &aux, "GetForGroup role assignments for group")
	if err != nil {
		return nil, err
	}
	return aux.Items, nil
}

// UnassignFromGroup removes the role assignment from a group
func (s *RoleAssignmentService) UnassignFromGroup(assignment *GroupRoleAssignment) error {
	// URL for unassigning role is /app/rest/userGroups/{groupLocator}/roles/{roleId}/{scope}
	return s.groupHelper.delete(fmt.Sprintf("%s/roles/%s/%s", assignment.GroupKey, assignment.RoleID, assignment.Scope), "UnassignFromGroup role from group")
}
