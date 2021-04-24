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

	createdRoot := setupFakeRoot(t, client, project, "Test Root")
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

func TestProjectFeature_CreateWithContextParameters(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project, "Test Root")
	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureVersionedSettings(project.ID, teamcity.ProjectFeatureVersionedSettingsOptions{
		Format:        teamcity.VersionedSettingsFormatKotlin,
		VcsRootID:     createdRoot.ID,
		BuildSettings: teamcity.VersionedSettingsBuildSettingsPreferVcs,
		ContextParameters: map[string]string{
			"Hello": "World",
		},
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)
	assert.Equal(t, "World", createdFeature.Properties().Map()["context.Hello"])
}

func TestProjectFeature_CreateXML(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project, "Test Root")
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

func TestProjectFeature_Delete(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project, "Test Root")
	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureVersionedSettings(project.ID, teamcity.ProjectFeatureVersionedSettingsOptions{
		Format:        teamcity.VersionedSettingsFormatXML,
		VcsRootID:     createdRoot.ID,
		BuildSettings: teamcity.VersionedSettingsBuildSettingsPreferVcs,
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)

	err = service.Delete(createdFeature.ID())
	require.NoError(t, err)

	deletedFeature, err := service.GetByID(createdFeature.ID())
	assert.NotNil(t, err)
	assert.Nil(t, deletedFeature)
}

func TestProjectFeature_Update(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project, "Test Root")
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

	type testData = struct {
		description     string
		enabled         bool
		buildSettings   teamcity.VersionedSettingsBuildSettings
		format          teamcity.VersionedSettingsFormat
		credentialsType teamcity.CredentialsStorageType
	}

	var validate = func(t *testing.T, id string, data testData) {
		retrievedFeature, err := service.GetByID(id)
		require.NoError(t, err)
		versionedSettings, ok := retrievedFeature.(*teamcity.ProjectFeatureVersionedSettings)
		assert.True(t, ok)

		assert.Equal(t, data.enabled, versionedSettings.Options.Enabled)
		assert.Equal(t, data.buildSettings, versionedSettings.Options.BuildSettings)
		assert.Equal(t, data.format, versionedSettings.Options.Format)
		assert.Equal(t, data.credentialsType, versionedSettings.Options.CredentialsStorageType)
	}
	t.Log("Validating initial creation")
	validate(t, createdFeature.ID(), testData{
		enabled:         true,
		credentialsType: teamcity.CredentialsStorageTypeScrambledInVcs,
		format:          teamcity.VersionedSettingsFormatXML,
		buildSettings:   teamcity.VersionedSettingsBuildSettingsPreferCurrent,
	})

	// then let's toggle some things
	updateConfigurations := []testData{
		{
			description:     "Switch to Kotlin",
			enabled:         true,
			buildSettings:   teamcity.VersionedSettingsBuildSettingsPreferVcs,
			format:          teamcity.VersionedSettingsFormatKotlin,
			credentialsType: teamcity.CredentialsStorageTypeCredentialsJSON,
		},
		{
			description:     "Switch back to XML",
			enabled:         true,
			buildSettings:   teamcity.VersionedSettingsBuildSettingsPreferVcs,
			format:          teamcity.VersionedSettingsFormatXML,
			credentialsType: teamcity.CredentialsStorageTypeScrambledInVcs,
		},
		{
			description:     "Disabled",
			enabled:         false,
			buildSettings:   teamcity.VersionedSettingsBuildSettingsPreferVcs,
			format:          teamcity.VersionedSettingsFormatXML,
			credentialsType: teamcity.CredentialsStorageTypeScrambledInVcs,
		},
		{
			description:     "Enabled & Prefer Current",
			enabled:         true,
			buildSettings:   teamcity.VersionedSettingsBuildSettingsPreferCurrent,
			format:          teamcity.VersionedSettingsFormatXML,
			credentialsType: teamcity.CredentialsStorageTypeScrambledInVcs,
		},
		{
			description:     "Always Use Current",
			enabled:         true,
			buildSettings:   teamcity.VersionedSettingsBuildSettingsAlwaysUseCurrent,
			format:          teamcity.VersionedSettingsFormatXML,
			credentialsType: teamcity.CredentialsStorageTypeScrambledInVcs,
		},
		{
			description:     "Kotlin with Scrambled",
			enabled:         true,
			buildSettings:   teamcity.VersionedSettingsBuildSettingsAlwaysUseCurrent,
			format:          teamcity.VersionedSettingsFormatXML,
			credentialsType: teamcity.CredentialsStorageTypeScrambledInVcs,
		},
		{
			description:     "Kotlin with CredentialsJSON",
			enabled:         true,
			buildSettings:   teamcity.VersionedSettingsBuildSettingsAlwaysUseCurrent,
			format:          teamcity.VersionedSettingsFormatXML,
			credentialsType: teamcity.CredentialsStorageTypeCredentialsJSON,
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
		settings.Options.CredentialsStorageType = update.credentialsType

		updatedFeature, err := service.Update(settings)
		require.NoError(t, err)
		assert.NotEmpty(t, updatedFeature.ID)

		// sanity check since we're updating with the same ID
		assert.Equal(t, createdFeature.ID(), updatedFeature.ID())

		validate(t, updatedFeature.ID(), update)
	}
}

func TestProjectFeature_UpdateContextParameters(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project, "Test Root")
	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureVersionedSettings(project.ID, teamcity.ProjectFeatureVersionedSettingsOptions{
		Format:        teamcity.VersionedSettingsFormatXML,
		VcsRootID:     createdRoot.ID,
		Enabled:       true,
		BuildSettings: teamcity.VersionedSettingsBuildSettingsPreferCurrent,
		ContextParameters: map[string]string{
			"Hello": "World",
		},
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)
	assert.Equal(t, "World", createdFeature.Properties().Map()["context.Hello"])

	retrieved, err := service.GetByID(createdFeature.ID())
	require.NoError(t, err)
	retrievedFeature := retrieved.(*teamcity.ProjectFeatureVersionedSettings)
	assert.Equal(t, "World", retrievedFeature.Options.ContextParameters["Hello"])

	feature.SetID(createdFeature.ID())
	feature.Options.ContextParameters["Hello"] = "London"
	_, err = service.Update(feature)

	require.NoError(t, err)
	retrieved, err = service.GetByID(createdFeature.ID())
	require.NoError(t, err)
	retrievedFeature = retrieved.(*teamcity.ProjectFeatureVersionedSettings)
	assert.Equal(t, "London", retrievedFeature.Options.ContextParameters["Hello"])

	feature.Options.ContextParameters["Hello"] = "World"
	feature.Options.ContextParameters["Germany"] = "Deutschland"
	_, err = service.Update(feature)

	require.NoError(t, err)
	retrieved, err = service.GetByID(createdFeature.ID())
	require.NoError(t, err)
	retrievedFeature = retrieved.(*teamcity.ProjectFeatureVersionedSettings)
	assert.Equal(t, "World", retrievedFeature.Options.ContextParameters["Hello"])
	assert.Equal(t, "Deutschland", retrievedFeature.Options.ContextParameters["Germany"])

	feature.Options.ContextParameters = make(map[string]string)
	_, err = service.Update(feature)
	require.NoError(t, err)

	retrieved, err = service.GetByID(createdFeature.ID())
	require.NoError(t, err)
	retrievedFeature = retrieved.(*teamcity.ProjectFeatureVersionedSettings)
	assert.Equal(t, 0, len(retrievedFeature.Options.ContextParameters))
}

func TestProjectFeature_UpdateVCSRoot(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project, "First Root")
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

	updatedRoot := setupFakeRoot(t, client, project, "Second Root")
	existing, err := service.GetByID(createdFeature.ID())
	require.NoError(t, err)

	existingFeatures, ok := existing.(*teamcity.ProjectFeatureVersionedSettings)
	assert.True(t, ok)
	assert.Equal(t, createdRoot.ID, existingFeatures.Options.VcsRootID)
	existingFeatures.Options.VcsRootID = updatedRoot.ID

	_, err = service.Update(existingFeatures)
	require.NoError(t, err)

	existing, err = service.GetByID(createdFeature.ID())
	require.NoError(t, err)

	existingFeatures, ok = existing.(*teamcity.ProjectFeatureVersionedSettings)
	assert.True(t, ok)
	assert.Equal(t, updatedRoot.ID, existingFeatures.Options.VcsRootID)
}
