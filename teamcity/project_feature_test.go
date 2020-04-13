package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: Delete

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

	createdFeature, err := service.Create(feature)
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

	createdFeature, err := service.Create(feature)
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
		Enabled:       true,
		BuildSettings: teamcity.VersionedSettingsBuildSettingsPreferCurrent,
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)

	var validate = func(t *testing.T, id string, enabled bool, buildSettings teamcity.VersionedSettingsBuildSettings, format teamcity.VersionedSettingsFormat) {
		retrievedFeature, err := service.GetByID(id)
		require.NoError(t, err)
		versionedSettings, ok := retrievedFeature.(*teamcity.ProjectFeatureVersionedSettings)
		assert.True(t, ok)

		assert.Equal(t, enabled, versionedSettings.Options.Enabled)
		assert.Equal(t, buildSettings, versionedSettings.Options.BuildSettings)
		assert.Equal(t, format, versionedSettings.Options.Format)
	}
	t.Log("Validating initial creation")
	validate(t, createdFeature.ID(), true, teamcity.VersionedSettingsBuildSettingsPreferCurrent, teamcity.VersionedSettingsFormatXML)

	// then let's toggle some things
	updateConfigurations := []struct {
		description   string
		enabled       bool
		buildSettings teamcity.VersionedSettingsBuildSettings
		format        teamcity.VersionedSettingsFormat
	}{
		{
			description:   "Switch to Kotlin",
			enabled:       true,
			buildSettings: teamcity.VersionedSettingsBuildSettingsPreferVcs,
			format:        teamcity.VersionedSettingsFormatKotlin,
		},
		{
			description:   "Switch back to XML",
			enabled:       true,
			buildSettings: teamcity.VersionedSettingsBuildSettingsPreferVcs,
			format:        teamcity.VersionedSettingsFormatXML,
		},
		{
			description:   "Disabled",
			enabled:       false,
			buildSettings: teamcity.VersionedSettingsBuildSettingsPreferVcs,
			format:        teamcity.VersionedSettingsFormatXML,
		},
		{
			description:   "Enabled & Prefer Current",
			enabled:       true,
			buildSettings: teamcity.VersionedSettingsBuildSettingsPreferCurrent,
			format:        teamcity.VersionedSettingsFormatXML,
		},
		{
			description:   "Always Use Current",
			enabled:       true,
			buildSettings: teamcity.VersionedSettingsBuildSettingsAlwaysUseCurrent,
			format:        teamcity.VersionedSettingsFormatXML,
		},
		{
			description:   "Kotlin",
			enabled:       true,
			buildSettings: teamcity.VersionedSettingsBuildSettingsAlwaysUseCurrent,
			format:        teamcity.VersionedSettingsFormatXML,
		},
	}
	for _, update := range updateConfigurations {
		t.Logf("Testing %q", update.description)

		existing, err := service.GetByID(createdFeature.ID())
		require.NoError(t, err)

		settings, ok := existing.(*teamcity.ProjectFeatureVersionedSettings)
		assert.True(t, ok)

		settings.Options.BuildSettings = update.buildSettings
		settings.Options.Enabled = update.enabled
		settings.Options.Format = update.format

		updatedFeature, err := service.Update(settings)
		require.NoError(t, err)
		assert.NotEmpty(t, updatedFeature.ID)

		// sanity check since we're updating with the same ID
		assert.Equal(t, createdFeature.ID(), updatedFeature.ID())

		validate(t, updatedFeature.ID(), update.enabled, update.buildSettings, update.format)
	}
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

	createdFeature, err := service.Create(feature)
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

	createdFeature, err := service.Create(feature)
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
