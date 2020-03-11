package teamcity_test

import (
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

func cleanUpGroup(t *testing.T, client *teamcity.Client, key string) {
	client.Groups.Delete(key)
}
