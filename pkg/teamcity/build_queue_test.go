package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/require"
)

const testBuildQueueProjectId = "BuildQueueProject"

func Test_EnqueueBuild(t *testing.T) {
	client := setup()
	require := require.New(t)
	buildType := createTestBuildType(t, client, testBuildQueueProjectId)

	sut := client.Queue

	props := teamcity.NewPropertiesEmpty()
	props.AddOrReplaceValue("system.wait_for_upsert", "false")

	actual, err := sut.TriggerBuild(teamcity.NewTriggerBuildRequest(buildType.ID, props))

	cleanUpProject(t, client, testBuildQueueProjectId)
	require.NoError(err)
	require.NotNil(actual)
}
