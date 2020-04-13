package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectFeature_CreateKotlin(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project)
	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureVersionedSettings(project.ID, teamcity.ProjectFeatureVersionedSettingsOptions{
		Format:        teamcity.VersionedSettingsFormatKotlin,
		VcsRootID:     createdRoot.ID,
		BuildSettings: teamcity.VersionedSettingsBuildSettingsPreferVcs,
	})

	createdFeature, err := service.Put(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)
}

func TestProjectFeature_CreateXML(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project)
	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureVersionedSettings(project.ID, teamcity.ProjectFeatureVersionedSettingsOptions{
		Format:        teamcity.VersionedSettingsFormatXML,
		VcsRootID:     createdRoot.ID,
		BuildSettings: teamcity.VersionedSettingsBuildSettingsPreferVcs,
	})

	createdFeature, err := service.Put(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)
}

func TestProjectFeature_Update(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project)
	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureVersionedSettings(project.ID, teamcity.ProjectFeatureVersionedSettingsOptions{
		Format:        teamcity.VersionedSettingsFormatXML,
		VcsRootID:     createdRoot.ID,
		BuildSettings: teamcity.VersionedSettingsBuildSettingsPreferCurrent,
	})

	createdFeature, err := service.Put(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)

	feature.Options.Format = teamcity.VersionedSettingsFormatKotlin
	feature.Options.BuildSettings = teamcity.VersionedSettingsBuildSettingsPreferVcs

	updatedFeature, err := service.Put(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, updatedFeature.ID)
}

func TestProjectFeature_GetEmptyFeatures(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)
	features, err := service.Get()
	require.NoError(t, err)
	require.Empty(t, features)
}

func TestProjectFeature_GetWithOneCreatedFeature(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project)
	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureVersionedSettings(project.ID, teamcity.ProjectFeatureVersionedSettingsOptions{
		Format:        teamcity.VersionedSettingsFormatKotlin,
		VcsRootID:     createdRoot.ID,
		BuildSettings: teamcity.VersionedSettingsBuildSettingsPreferVcs,
	})

	createdFeature, err := service.Put(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)

	features, err := service.Get()
	require.NoError(t, err)
	require.Len(t, features, 1)
}

func TestProjectFeature_GetByIdWithCreatedFeature(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project)
	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureVersionedSettings(project.ID, teamcity.ProjectFeatureVersionedSettingsOptions{
		Format:        teamcity.VersionedSettingsFormatKotlin,
		VcsRootID:     createdRoot.ID,
		BuildSettings: teamcity.VersionedSettingsBuildSettingsPreferVcs,
	})

	createdFeature, err := service.Put(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)

	retrievedFeature, err := service.GetByID(createdFeature.ID())
	require.NoError(t, err)
	assert.Equal(t, createdFeature.ID(), retrievedFeature.ID())
	assert.Equal(t, createdFeature.Type(), retrievedFeature.Type())
}

func TestProjectFeature_GetByIdWithFeatureNotExisting(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)
	retrievedFeature, err := service.GetByID("random_id")
	require.Nil(t, retrievedFeature)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func setupFakeRoot(t *testing.T, client *teamcity.Client, project *teamcity.Project) *teamcity.VcsRootReference {
	rootOptions, err := teamcity.NewGitVcsRootOptionsDefaults("master", "git@test.com")
	require.NoError(t, err)

	root, err := teamcity.NewGitVcsRoot(project.ID, "Test Root", rootOptions)
	require.NoError(t, err)

	createdRoot, err := client.VcsRoots.Create(project.ID, root)
	require.NoError(t, err)

	return createdRoot
}
