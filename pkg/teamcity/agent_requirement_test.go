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
	buildType, _ = client.BuildTypes.GetById(buildType.ID) //refresh

	actual := buildType.AgentRequirements

	cleanUpProject(t, client, testBuildTypeProjectId)
	assert.NotEmpty(actual.Items)
	assert.Equal(teamcity.Conditions.Equals, actual.Items[0].Condition)

	assert.Equal("property-name", actual.Items[0].Properties.Items[0].Name)
	assert.Equal("param", actual.Items[0].Properties.Items[0].Value)
	assert.Equal("property-value", actual.Items[0].Properties.Items[1].Name)
	assert.Equal("value", actual.Items[0].Properties.Items[1].Value)
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

	require.NoError(err)
	assert.Equal(created.ID, actual.ID)
	assert.Equal(created.BuildTypeID, actual.BuildTypeID)
	assert.Equal(created.Name(), actual.Name())
	assert.Equal(created.Value(), actual.Value())

	cleanUpProject(t, client, bt.ProjectID)
}

func TestAgentRequirement_Delete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "AgentRequirementProject", "BuildRelease", true)

	sut := client.AgentRequirementService(bt.ID)
	nt, _ := teamcity.NewAgentRequirement(teamcity.Conditions.Equals, "param", "value")

	created, err := sut.Create(nt)

	require.Nil(err)

	sut.Delete(created.ID)
	_, err = sut.GetByID(created.ID) // refresh

	require.Error(err)
	assert.Contains(err.Error(), "404")
	cleanUpProject(t, client, bt.ProjectID)
}
