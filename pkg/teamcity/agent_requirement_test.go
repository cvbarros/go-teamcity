package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentRequirement_Create(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)
	req, _ := teamcity.NewAgentRequirement(teamcity.Conditions.Equals, "param", "value")

	sut := client.AgentRequirementService(buildType.ID)
	sut.Create(req)
	buildType, _ = client.BuildTypes.GetByID(buildType.ID) //refresh

	all, _ := client.AgentRequirementService(buildType.ID).GetAll()
	cleanUpProject(t, client, testBuildTypeProjectId)

	require.NotEmpty(t, all)
	actual := all[0]

	assert.Equal(teamcity.Conditions.Equals, actual.Condition)
	assert.Equal("property-name", actual.Properties.Items[0].Name)
	assert.Equal("param", actual.Properties.Items[0].Value)
	assert.Equal("property-value", actual.Properties.Items[1].Name)
	assert.Equal("value", actual.Properties.Items[1].Value)
}

func TestAgentRequirement_Get(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "AgentRequirementProject", "BuildRelease", true)

	sut := client.AgentRequirementService(bt.ID)
	nt, _ := teamcity.NewAgentRequirement(teamcity.Conditions.Equals, "param", "value")

	created, err := sut.Create(nt)

	require.Nil(err)

	actual, err := sut.GetByID(created.ID)
	cleanUpProject(t, client, bt.ProjectID)

	require.NoError(err)
	assert.Equal(created.ID, actual.ID)
	assert.Equal(created.BuildTypeID, actual.BuildTypeID)
	assert.Equal(created.Name(), actual.Name())
	assert.Equal(created.Value(), actual.Value())
}

func TestAgentRequirement_GetAll(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "AgentRequirementProject", "BuildRelease", true)

	sut := client.AgentRequirementService(bt.ID)
	req1, _ := teamcity.NewAgentRequirement(teamcity.Conditions.Equals, "param", "value")
	req2, _ := teamcity.NewAgentRequirement(teamcity.Conditions.DoesNotEqual, "param2", "value2")

	created1, err := sut.Create(req1)
	require.NoError(err)
	created2, err := sut.Create(req2)
	require.NoError(err)

	actual, err := sut.GetAll()
	cleanUpProject(t, client, bt.ProjectID)

	require.NoError(err)
	require.Equal(2, len(actual))

	assert.Equal(created1.ID, actual[0].ID)
	assert.Equal(created1.BuildTypeID, actual[0].BuildTypeID)
	assert.Equal(created1.Condition, actual[0].Condition)
	assert.Equal(created2.ID, actual[1].ID)
	assert.Equal(created2.BuildTypeID, actual[1].BuildTypeID)
	assert.Equal(created2.Condition, actual[1].Condition)
}

func TestAgentRequirement_Delete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "AgentRequirementProject", "BuildRelease", true)

	sut := client.AgentRequirementService(bt.ID)
	nt, _ := teamcity.NewAgentRequirement(teamcity.Conditions.Equals, "param", "value")

	created, err := sut.Create(nt)
	cleanUpProject(t, client, bt.ProjectID)

	require.Nil(err)

	sut.Delete(created.ID)
	_, err = sut.GetByID(created.ID) // refresh

	require.Error(err)
	assert.Contains(err.Error(), "404")
}
