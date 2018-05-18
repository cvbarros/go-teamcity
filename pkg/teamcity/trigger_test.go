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

	actual := teamcity.NewVcsTrigger("+:*", "")

	require.NotNil(actual)
	assert.Equal("vcsTrigger", actual.Type)
	require.NotEmpty(actual.Properties)
	props := actual.Properties.Map()

	assert.Contains(props, "triggerRules")
	assert.NotContains(props, "branchFilter")
	assert.Equal(props["quietPeriodMode"], "DO_NOT_USE")
	assert.Equal(props["enableQueueOptimization"], "true")
}

func TestTrigger_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt := teamcity.NewVcsTrigger("+:*", "")

	_, err := sut.AddTrigger(nt)

	require.Nil(err)

	bt, _ = client.BuildTypes.GetById(bt.ID)

	assert.Equal(int32(1), bt.Triggers.Count)
	actual := bt.Triggers.Items[0]

	assert.NotEmpty(actual.ID)
	assert.Equal("vcsTrigger", actual.Type)
	assert.NotEmpty(actual.Properties)

	cleanUpProject(t, client, bt.ProjectID)
}

func TestTrigger_Get(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt := teamcity.NewVcsTrigger("+:*", "")

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	actual, err := sut.GetById(created.ID)

	require.NoError(err)
	assert.Equal(created.ID, actual.ID)
	assert.Equal(created.BuildTypeID, actual.BuildTypeID)
	assert.Equal(created.Type, actual.Type)

	cleanUpProject(t, client, bt.ProjectID)
}

func TestTrigger_Delete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt := teamcity.NewVcsTrigger("+:*", "")

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	sut.Delete(created.ID)
	_, err = sut.GetById(created.ID) // refresh

	require.Error(err)
	assert.Contains(err.Error(), "404")
	cleanUpProject(t, client, bt.ProjectID)
}
