package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/require"
)

func TestUserGroupMemberShip_IsGroupMember(t *testing.T) {
	client := setup()
	admin, err := client.Users.GetByUsername("admin")
	require.NoError(t, err)
	groupKey := "ALL_USERS_GROUP"

	isMember, err := client.UserGroupMemberShip.IsGroupMemberByID(admin.ID, groupKey)
	require.NoError(t, err)
	require.True(t, isMember)

	isMember, err = client.UserGroupMemberShip.IsGroupMemberByID(admin.ID, "INVALID_GROUP")
	require.NoError(t, err)
	require.False(t, isMember)

	isMember, err = client.UserGroupMemberShip.IsGroupMemberByUsername("INVALIDUSERNAME", groupKey)
	require.NoError(t, err)
	require.False(t, isMember)

}

func TestUserGropMemberShip_GetGroupMembers(t *testing.T) {
	client := setup()

	newGroup, _ := teamcity.NewGroup("TESTGROUPMEMBER", "Test Group Member", "")
	actualGroup, err := client.Groups.Create(newGroup)
	require.NoError(t, err)

	newUser, _ := teamcity.NewUser("testusermember", "Test User Member", "test@member.com")
	actualUser, err := client.Users.Create(newUser)
	require.NoError(t, err)

	actualGroup, err = client.Users.GroupAddByID(actualUser.ID, actualGroup.Key)
	require.NoError(t, err)
	require.NotNil(t, actualGroup)

	memberList, err := client.UserGroupMemberShip.GetGroupMembersListAllByKey(actualGroup.Key)
	require.NoError(t, err)
	require.NotEmpty(t, memberList.Items)
	require.Equal(t, memberList.Count, 1)
	require.Equal(t, actualUser.ID, memberList.Items[0].ID)

	require.NoError(t, client.UserGroupMemberShip.GroupDeleteMemberByID(actualUser.ID, actualGroup.Key))

	memberList, err = client.UserGroupMemberShip.GetGroupMembersListAllByKey(actualGroup.Key)
	require.NoError(t, err)
	require.Empty(t, memberList.Items)
	require.Equal(t, memberList.Count, 0)

	cleanUpGroup(t, client, newGroup.Key)
	cleanUpUser(t, client, actualUser.ID)

	require.NoError(t, err)
	require.NotNil(t, actualGroup)
}
func TestUserGropMemberShip_GetUserGroups(t *testing.T) {
	client := setup()

	newGroup, _ := teamcity.NewGroup("TESTGUSERGROUPS", "Test User groups", "")
	actualGroup, err := client.Groups.Create(newGroup)
	require.NoError(t, err)

	newUser, _ := teamcity.NewUser("testusergroups", "Test User Groups", "testgroups@member.com")
	actualUser, err := client.Users.Create(newUser)
	require.NoError(t, err)

	actualGroup, err = client.Users.GroupAddByID(actualUser.ID, actualGroup.Key)
	require.NoError(t, err)
	require.NotNil(t, actualGroup)

	groupsList, err := client.UserGroupMemberShip.GetUserGroupsListAllByID(actualUser.ID)
	require.NoError(t, err)
	require.NotEmpty(t, groupsList.Items)
	// first is ALL_USERS_GROUP
	require.Equal(t, groupsList.Count, 2)
	require.Equal(t, actualGroup.Key, groupsList.Items[1].Key)

	require.NoError(t, client.UserGroupMemberShip.GroupDeleteMemberByID(actualUser.ID, actualGroup.Key))

	groupsList, err = client.UserGroupMemberShip.GetUserGroupsListAllByID(actualUser.ID)
	require.NoError(t, err)
	require.Equal(t, groupsList.Count, 1)

	cleanUpGroup(t, client, newGroup.Key)
	cleanUpUser(t, client, actualUser.ID)

	require.NoError(t, err)
	require.NotNil(t, actualGroup)
}
