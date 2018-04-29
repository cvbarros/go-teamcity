package teamcity_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
)

func TestAddSnapshotDependency(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)
	buildTypeDep := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, "DependencyBuild")

	sut := client.DependencyService(buildType.ID)

	dep := teamcity.NewSnapshotDependency(buildTypeDep.Reference())
	err := sut.AddSnapshotDependency(dep)

	require.Nil(t, err)

	buildType, _ = client.BuildTypes.GetById(buildType.ID) //refresh
	actual := buildType.SnapshotDependencies.Items

	cleanUpProject(t, client, testBuildTypeProjectId)

	assert.Equal(1, len(actual))
	assert.Equal("snapshot_dependency", actual[0].Type)
	assert.NotEmpty(actual[0].Properties)
}
