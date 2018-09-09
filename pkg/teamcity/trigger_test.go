package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrigger_Constructor(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	actual, _ := teamcity.NewTriggerVcs("+:*", "")

	require.NotNil(actual)
	assert.Equal("vcsTrigger", actual.Type())

	assert.Equal("+:*", actual.Rules)
	assert.Empty(actual.BranchFilter)
}

func TestTrigger_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerVcs("+:*", "")

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	assert.Equal(created.BuildTypeID(), bt.ID)
	cleanUpProject(t, client, bt.ProjectID)
}

func TestTrigger_Get(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerVcs("+:*", "")

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	actual, err := sut.GetByID(created.ID())

	require.NoError(err)
	assert.Equal(created.ID(), actual.ID())
	assert.Equal(created.BuildTypeID(), actual.BuildTypeID())
	assert.Equal(created.Type(), actual.Type())

	cleanUpProject(t, client, bt.ProjectID)
}

func TestTrigger_Delete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerVcs("+:*", "")

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	sut.Delete(created.ID())
	_, err = sut.GetByID(created.ID()) // refresh

	require.Error(err)
	assert.Contains(err.Error(), "404")
	cleanUpProject(t, client, bt.ProjectID)
}
