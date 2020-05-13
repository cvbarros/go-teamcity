package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupRoleAssignment_GlobalSysAdminAssign(t *testing.T) {
	client := setup()

	newGroup, _ := teamcity.NewGroup("TESTGROUPKEY", "Test Group Name", "Test Group Description")
	newGroupRoleAssignment, _ := teamcity.NewGroupRoleAssignment("TESTGROUPKEY", "SYSTEM_ADMIN", "g") // "g" is for sys admins at the global level

	actualGroup, err := client.Groups.Create(newGroup)
	require.NoError(t, err)
	require.NotNil(t, actualGroup)

	createdGroupRoleAssignment, err := client.GroupRoleAssignments.Assign(newGroupRoleAssignment)
	require.NoError(t, err)
	require.NotNil(t, createdGroupRoleAssignment)
	assert.NotEmpty(t, createdGroupRoleAssignment.RoleID)
	assert.NotEmpty(t, createdGroupRoleAssignment.Scope)
	assert.NotEmpty(t, createdGroupRoleAssignment.Href)

	groupRoleAssignments, err := client.GroupRoleAssignments.GetAll(newGroup)
	require.NoError(t, err)
	assert.Equal(t, 1, len(groupRoleAssignments))
	actualGroupRoleAssignmentReference := groupRoleAssignments[0]
	assert.Equal(t, createdGroupRoleAssignment.RoleID, actualGroupRoleAssignmentReference.RoleID)
	assert.Equal(t, createdGroupRoleAssignment.Scope, actualGroupRoleAssignmentReference.Scope)
	assert.Equal(t, createdGroupRoleAssignment.Href, actualGroupRoleAssignmentReference.Href)

	actualGroupRoleAssignmentReference2, err := client.GroupRoleAssignments.Get(newGroupRoleAssignment)
	require.NoError(t, err)
	assert.Equal(t, createdGroupRoleAssignment.RoleID, actualGroupRoleAssignmentReference2.RoleID)
	assert.Equal(t, createdGroupRoleAssignment.Scope, actualGroupRoleAssignmentReference2.Scope)
	assert.Equal(t, createdGroupRoleAssignment.Href, actualGroupRoleAssignmentReference2.Href)

	// Clean up group after test
	cleanUpGroup(t, client, actualGroup.Key)
}

func TestGroupRoleAssignment_Assign(t *testing.T) {
	client := setup()

	parent, _ := teamcity.NewProject("ParentProject", "Parent Project", "")
	child, _ := teamcity.NewProject("ChildProject", "Child Project", "ParentProject")

	_, err := client.Projects.Create(parent)
	require.NoError(t, err)
	created, err := client.Projects.Create(child)
	require.NoError(t, err)

	newGroup, _ := teamcity.NewGroup("TESTGROUPKEY", "Test Group Name", "Test Group Description")
	newGroupRoleAssignment, _ := teamcity.NewGroupRoleAssignment("TESTGROUPKEY", "PROJECT_DEVELOPER", "p:"+created.ID)

	actualGroup, err := client.Groups.Create(newGroup)
	require.NoError(t, err)
	require.NotNil(t, actualGroup)

	createdGroupRoleAssignment, err := client.GroupRoleAssignments.Assign(newGroupRoleAssignment)
	require.NoError(t, err)
	require.NotNil(t, createdGroupRoleAssignment)
	assert.NotEmpty(t, createdGroupRoleAssignment.RoleID)
	assert.NotEmpty(t, createdGroupRoleAssignment.Scope)
	assert.NotEmpty(t, createdGroupRoleAssignment.Href)

	groupRoleAssignments, err := client.GroupRoleAssignments.GetAll(newGroup)
	require.NoError(t, err)
	assert.Equal(t, 1, len(groupRoleAssignments))
	actualGroupRoleAssignmentReference := groupRoleAssignments[0]
	assert.Equal(t, createdGroupRoleAssignment.RoleID, actualGroupRoleAssignmentReference.RoleID)
	assert.Equal(t, createdGroupRoleAssignment.Scope, actualGroupRoleAssignmentReference.Scope)
	assert.Equal(t, createdGroupRoleAssignment.Href, actualGroupRoleAssignmentReference.Href)

	actualGroupRoleAssignmentReference2, err := client.GroupRoleAssignments.Get(newGroupRoleAssignment)
	require.NoError(t, err)
	assert.Equal(t, createdGroupRoleAssignment.RoleID, actualGroupRoleAssignmentReference2.RoleID)
	assert.Equal(t, createdGroupRoleAssignment.Scope, actualGroupRoleAssignmentReference2.Scope)
	assert.Equal(t, createdGroupRoleAssignment.Href, actualGroupRoleAssignmentReference2.Href)

	// Clean up group and projects after test
	cleanUpGroup(t, client, actualGroup.Key)
	cleanUpProject(t, client, "ParentProject")
}

func TestGroupRoleAssignment_Unassign(t *testing.T) {
	client := setup()

	parent, _ := teamcity.NewProject("ParentProject", "Parent Project", "")
	child, _ := teamcity.NewProject("ChildProject", "Child Project", "ParentProject")

	_, err := client.Projects.Create(parent)
	require.NoError(t, err)
	created, err := client.Projects.Create(child)
	require.NoError(t, err)

	newGroup, _ := teamcity.NewGroup("TESTGROUPKEY", "Test Group Name", "Test Group Description")
	newGroupRoleAssignment, _ := teamcity.NewGroupRoleAssignment("TESTGROUPKEY", "PROJECT_DEVELOPER", "p:"+created.ID)

	actualGroup, err := client.Groups.Create(newGroup)
	require.NoError(t, err)
	require.NotNil(t, actualGroup)

	createdGroupRoleAssignment, err := client.GroupRoleAssignments.Assign(newGroupRoleAssignment)
	require.NoError(t, err)
	require.NotNil(t, createdGroupRoleAssignment)
	assert.NotEmpty(t, createdGroupRoleAssignment.RoleID)
	assert.NotEmpty(t, createdGroupRoleAssignment.Scope)
	assert.NotEmpty(t, createdGroupRoleAssignment.Href)

	err = client.GroupRoleAssignments.Unassign(newGroupRoleAssignment)
	require.NoError(t, err)

	groupRoleAssignments, err := client.GroupRoleAssignments.GetAll(newGroup)
	require.NoError(t, err)
	assert.Equal(t, 0, len(groupRoleAssignments))

	// The Role has been unassigneded, so expect error, and message to contain 404 (NOT FOUND)
	_, err = client.GroupRoleAssignments.Get(newGroupRoleAssignment)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")

	// Clean up group and projects after test
	cleanUpGroup(t, client, actualGroup.Key)
	cleanUpProject(t, client, "ParentProject")
}
