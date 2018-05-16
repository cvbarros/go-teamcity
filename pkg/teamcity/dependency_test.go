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
	buildTypeDep := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, "DependencyBuild", false)

	sut := client.DependencyService(buildType.ID)

	dep := teamcity.NewSnapshotDependency(buildTypeDep.ID)
	_, err := sut.AddSnapshotDependency(dep)

	require.Nil(t, err)

	buildType, _ = client.BuildTypes.GetById(buildType.ID) //refresh
	actual := buildType.SnapshotDependencies.Items

	cleanUpProject(t, client, testBuildTypeProjectId)

	assert.Equal(1, len(actual))
	assert.Equal("snapshot_dependency", actual[0].Type)
	assert.NotEmpty(actual[0].Properties)
}

func TestGetSnapshotDependency(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)
	buildTypeDep := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, "DependencyBuild", false)

	sut := client.DependencyService(buildType.ID)

	dep := teamcity.NewSnapshotDependency(buildTypeDep.ID)
	created, err := sut.AddSnapshotDependency(dep)

	require.Nil(t, err)

	actual, err := sut.GetById(created.ID) // refresh

	require.Nil(t, err)
	assert.Equal(created.ID, actual.ID)
	assert.Equal(created.BuildTypeID, actual.BuildTypeID)
	assert.Equal(created.Type, actual.Type)
	assert.Equal(created.SourceBuildType.ID, actual.SourceBuildType.ID)

	cleanUpProject(t, client, testBuildTypeProjectId)
}
