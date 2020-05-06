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

// GroupRoleAssignmentReference represents a response of a request to assign role to a group
type GroupRoleAssignmentReference struct {
	RoleID string `json:"roleId,omitempty" xml:"roleId"`
	Scope  string `json:"scope,omitempty" xml:"scope"`
	Href   string `json:"href,omitempty" xml:"href"`
}

type groupRoleAssignmentsJSON struct {
	Items []GroupRoleAssignmentReference `json:"role"`
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
		return nil, fmt.Errorf("scope is required. Use the Project ID or use \"_Root\" for the Root project")
	}

	return &GroupRoleAssignment{
		GroupKey: groupKey,
		RoleID:   roleID,
		Scope:    scope,
	}, nil
}

// GroupRoleAssignmentService has operations for handling role assignments for groups
type GroupRoleAssignmentService struct {
	sling      *sling.Sling
	httpClient *http.Client
	restHelper *restHelper
}

func newGroupRoleAssignmentService(base *sling.Sling, httpClient *http.Client) *GroupRoleAssignmentService {
	sling := base.Path(fmt.Sprintf("userGroups/"))
	return &GroupRoleAssignmentService{
		httpClient: httpClient,
		sling:      sling,
		restHelper: newRestHelperWithSling(httpClient, sling),
	}
}

// Assign adds a role assignment to a group
func (s *GroupRoleAssignmentService) Assign(assignment *GroupRoleAssignment) (*GroupRoleAssignmentReference, error) {
	var out GroupRoleAssignmentReference

	// URL for assigning role is /app/rest/userGroups/{groupLocator}/roles/{roleId}/p:{scope}
	err := s.restHelper.post(fmt.Sprintf("%s/roles/%s/p:%s", assignment.GroupKey, assignment.RoleID, assignment.Scope), nil, &out, "Assign role to group")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Get get a specific role assignment for a group
func (s *GroupRoleAssignmentService) Get(assignment *GroupRoleAssignment) (*GroupRoleAssignmentReference, error) {
	var out GroupRoleAssignmentReference

	// URL for getting a specific role assignments is /app/rest/userGroups/{groupLocator}/roles/{roleId}/p:{scope}
	err := s.restHelper.get(fmt.Sprintf("%s/roles/%s/p:%s", assignment.GroupKey, assignment.RoleID, assignment.Scope), &out, "Get role assignmens for group")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAll gets all the role assignments for a group
func (s *GroupRoleAssignmentService) GetAll(group *Group) ([]GroupRoleAssignmentReference, error) {
	var aux groupRoleAssignmentsJSON

	// URL for getting role assignments is /app/rest/userGroups/{groupLocator}/roles
	err := s.restHelper.get(fmt.Sprintf("%s/roles", group.Key), &aux, "Get role assignments for group")
	if err != nil {
		return nil, err
	}
	return aux.Items, nil
}

// Unassign removes the role assignment from a group
func (s *GroupRoleAssignmentService) Unassign(assignment *GroupRoleAssignment) error {
	// URL for unassigning role is /app/rest/userGroups/{groupLocator}/roles/{roleId}/p:{scope}
	return s.restHelper.delete(fmt.Sprintf("%s/roles/%s/p:%s", assignment.GroupKey, assignment.RoleID, assignment.Scope), "Unassign role from group")
}
