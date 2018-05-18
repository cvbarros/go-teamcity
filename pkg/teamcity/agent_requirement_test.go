package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
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
