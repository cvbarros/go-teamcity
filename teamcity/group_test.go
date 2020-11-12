package teamcity_test

import (
	"fmt"
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroup_Create(t *testing.T) {
	newGroup, _ := teamcity.NewGroup("TESTGROUPKEY", "Test Group Name", "Test Group Description")
	client := setup()
	actual, err := client.Groups.Create(newGroup)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.NotEmpty(t, actual.Key)

	cleanUpGroup(t, client, actual.Key)

	assert.Equal(t, newGroup.Key, actual.Key)
	assert.Equal(t, newGroup.Name, actual.Name)
	assert.Equal(t, newGroup.Description, actual.Description)
}

func TestGroup_GetByKey(t *testing.T) {
	newGroup, _ := teamcity.NewGroup("TESTGROUPKEY", "Test Group Name", "Test Group Description")
	client := setup()
	client.Groups.Create(newGroup)

	actual, err := client.Groups.GetByKey(newGroup.Key)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.NotEmpty(t, actual.Key)

	cleanUpGroup(t, client, actual.Key)

	assert.Equal(t, newGroup.Key, actual.Key)
	assert.Equal(t, newGroup.Name, actual.Name)
	assert.Equal(t, newGroup.Description, actual.Description)
}

func TestGroup_GetByName(t *testing.T) {
	newGroup, _ := teamcity.NewGroup("TESTGROUPKEY2", "Test Group Name 2", "Test Group Description 2")
	client := setup()
	client.Groups.Create(newGroup)

	actual, err := client.Groups.GetByName(newGroup.Name)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.NotEmpty(t, actual.Key)

	cleanUpGroup(t, client, actual.Key)

	assert.Equal(t, newGroup.Key, actual.Key)
	assert.Equal(t, newGroup.Name, actual.Name)
	assert.Equal(t, newGroup.Description, actual.Description)
}

func TestGroup_Delete(t *testing.T) {
	newGroup, _ := teamcity.NewGroup("TESTGROUPKEY", "Test Group Name", "Test Group Description")
	client := setup()
	client.Groups.Create(newGroup)

	err := client.Groups.Delete(newGroup.Key)

	require.NoError(t, err)

	_, err = client.Groups.GetByKey(newGroup.Key)

	// Group is deleted, so expect error, and message to contain 404 (NOT FOUND)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestGroup_List(t *testing.T) {
	groups := []*teamcity.Group{}
	client := setup()

	groupListBefore, err := client.Groups.List()
	require.NoError(t, err)
	for i := 0; i < 5; i++ {
		group, err := teamcity.NewGroup(
			fmt.Sprint("TESTGROUPLIST", i),
			fmt.Sprint("Test Group List ", i),
			fmt.Sprint("Test Group Description List ", i),
		)
		require.NoError(t, err)
		groups = append(groups, group)
		client.Groups.Create(group)
	}
	groupList, err := client.Groups.List()

	require.NoError(t, err)
	assert.Equal(t, groupListBefore.Count+5, groupList.Count)

	for _, group := range groupList.Items {
		if group.Key == "ALL_USERS_GROUP" {
			continue
		}
		_, err := client.Groups.GetByKey(group.Key)
		require.NoError(t, err)

		cleanUpGroup(t, client, group.Key)

		_, err = client.Groups.GetByKey(group.Key)
		require.Error(t, err)
	}
}

func cleanUpGroup(t *testing.T, client *teamcity.Client, key string) {
	client.Groups.Delete(key)
}
