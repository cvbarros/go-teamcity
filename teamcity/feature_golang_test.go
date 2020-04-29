package teamcity_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeatureGolang_Lifecycle(t *testing.T) {
	name := fmt.Sprintf("Project %d", time.Now().Unix())
	newProject := getTestProjectData(name, "")
	client := setup()

	project, err := client.Projects.Create(newProject)
	require.NoError(t, err)
	defer cleanUpProject(t, client, project.ID)

	buildType, err := teamcity.NewBuildType(project.ID, "Hello")
	require.NoError(t, err)
	buildConfig, err := client.BuildTypes.Create("", buildType)
	require.NoError(t, err)

	service := client.BuildFeatureService(buildConfig.ID)
	feature := teamcity.NewFeatureGolang()
	feature.SetBuildTypeID(buildConfig.ID)
	createdService, err := service.Create(feature)
	require.NoError(t, err)

	retrievedService, err := service.GetByID(createdService.ID())
	require.NoError(t, err)
	assert.False(t, retrievedService.Disabled())
}

func TestFeatureGolang_UnmarshallProperties(t *testing.T) {
	assert := assert.New(t)
	var actual teamcity.FeatureGolangPublisher
	const json = `
	{
		"id": "BUILD_EXT_1",
		"type": "golang",
		"properties": {
			"count": 1,
			"property": [
				{
					"name": "test.format",
					"value": "json"
				}
			]
		}
	}
	`
	actual.UnmarshalJSON([]byte(json))

	assert.Equal("BUILD_EXT_1", actual.ID())
}
