package teamcity_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
)

func TestSnapshotDependency_Create(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)
	buildTypeDep := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, "DependencyBuild", false)

	sut := client.DependencyService(buildType.ID)

	dep := teamcity.NewSnapshotDependency(buildTypeDep.ID)
	created, err := sut.AddSnapshotDependency(dep)

	require.Nil(t, err)

	buildType, _ = client.BuildTypes.GetByID(buildType.ID) //refresh
	actual, _ := sut.GetSnapshotByID(created.ID)

	cleanUpProject(t, client, testBuildTypeProjectId)

	assert.Equal("snapshot_dependency", actual.Type)
	assert.Equal(buildTypeDep.ID, actual.SourceBuildType.ID)
	assert.Equal(buildType.ID, actual.BuildTypeID)
}

func TestSnapshotDependency_Get(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)
	buildTypeDep := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, "DependencyBuild", false)

	sut := client.DependencyService(buildType.ID)

	dep := teamcity.NewSnapshotDependency(buildTypeDep.ID)
	created, err := sut.AddSnapshotDependency(dep)

	require.Nil(t, err)

	actual, err := sut.GetSnapshotByID(created.ID) // refresh

	require.Nil(t, err)
	assert.Equal(created.ID, actual.ID)
	assert.Equal(created.BuildTypeID, actual.BuildTypeID)
	assert.Equal(created.Type, actual.Type)
	assert.Equal(created.SourceBuildType.ID, actual.SourceBuildType.ID)

	cleanUpProject(t, client, testBuildTypeProjectId)
}

func TestSnapshotDependency_Delete(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)
	buildTypeDep := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, "DependencyBuild", false)

	sut := client.DependencyService(buildType.ID)

	dep := teamcity.NewSnapshotDependency(buildTypeDep.ID)
	created, err := sut.AddSnapshotDependency(dep)

	require.Nil(t, err)

	sut.DeleteSnapshot(created.ID)
	_, err = sut.GetSnapshotByID(created.ID) // refresh

	require.Error(t, err)
	assert.Contains(err.Error(), "404")
	cleanUpProject(t, client, testBuildTypeProjectId)
}

func TestArtifactDependency_Create(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)
	buildTypeDep := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, "DependencyBuild", false)

	sut := client.DependencyService(buildType.ID)

	dep, _ := teamcity.NewArtifactDependency(buildTypeDep.ID, createDefaultTestingArtifactDependencyOptions())
	created, err := sut.AddArtifactDependency(dep)

	require.Nil(t, err)

	buildType, _ = client.BuildTypes.GetByID(buildType.ID) //refresh
	actual, _ := sut.GetArtifactByID(created.ID())

	cleanUpProject(t, client, testBuildTypeProjectId)

	assert.Equal("artifact_dependency", actual.Type())
}

func TestArtifactDependency_Get(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)
	buildTypeDep := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, "DependencyBuild", false)

	sut := client.DependencyService(buildType.ID)

	dep, _ := teamcity.NewArtifactDependency(buildTypeDep.ID, createDefaultTestingArtifactDependencyOptions())
	created, err := sut.AddArtifactDependency(dep)

	require.Nil(t, err)

	actual, err := sut.GetArtifactByID(created.ID()) // refresh

	require.Nil(t, err)
	assert.Equal(created.ID(), actual.ID())
	assert.Equal(created.BuildTypeID(), actual.BuildTypeID())
	assert.Equal(created.Type(), actual.Type())
	assert.Equal(created.SourceBuildTypeID, actual.SourceBuildTypeID)

	cleanUpProject(t, client, testBuildTypeProjectId)
}

func TestArtifactDependency_Delete(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)
	buildTypeDep := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, "DependencyBuild", false)

	sut := client.DependencyService(buildType.ID)

	dep, _ := teamcity.NewArtifactDependency(buildTypeDep.ID, createDefaultTestingArtifactDependencyOptions())
	created, err := sut.AddArtifactDependency(dep)

	require.Nil(t, err)

	sut.DeleteArtifact(created.ID())
	_, err = sut.GetArtifactByID(created.ID()) // refresh

	require.Error(t, err)
	assert.Contains(err.Error(), "404")
	cleanUpProject(t, client, testBuildTypeProjectId)
}
