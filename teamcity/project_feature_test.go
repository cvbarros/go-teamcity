package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	features, err := service.Get()
	require.NoError(t, err)
	require.Len(t, features, 1)
}

func TestProjectFeature_GetByIdWithCreatedFeature(t *testing.T) {
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

	retrieved, err := service.GetByID(createdFeature.ID())
	require.NoError(t, err)
	assert.Equal(t, createdFeature.ID(), retrieved.ID())
	assert.Equal(t, createdFeature.Type(), retrieved.Type())

	// sanity check - should be none by default
	retrievedFeature := retrieved.(*teamcity.ProjectFeatureVersionedSettings)
	assert.Equal(t, 0, len(retrievedFeature.Options.ContextParameters))
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

func TestProjectFeature_GetByType(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	createdRoot := setupFakeRoot(t, client, project, "Test Root")
	service := client.ProjectFeatureService(project.ID)

	featureVersionedSettings := teamcity.NewProjectFeatureVersionedSettings(project.ID, teamcity.ProjectFeatureVersionedSettingsOptions{
		Format:        teamcity.VersionedSettingsFormatKotlin,
		VcsRootID:     createdRoot.ID,
		BuildSettings: teamcity.VersionedSettingsBuildSettingsPreferVcs,
	})

	featureVaultConnection := teamcity.NewProjectConnectionVault(project.ID, teamcity.ConnectionProviderVaultOptions{
		DisplayName: "Hashicorp Vault",
		URL:         "http://vault.service:8200",
	})

	createdFeatureVersionedSettings, err := service.Create(featureVersionedSettings)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeatureVersionedSettings.ID)

	createdFeatureVaultConnection, err := service.Create(featureVaultConnection)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeatureVaultConnection.ID)

	retrievedVersionedSettings, err := service.GetByType("versionedSettings")
	require.NoError(t, err)
	// we can't compare the ID, since they change after creation
	assert.Equal(t, createdFeatureVersionedSettings.Type(), retrievedVersionedSettings.Type())

	retrievedVaultConnection, err := service.GetByType("OAuthProvider")
	require.NoError(t, err)
	// we can't compare the ID, since they change after creation
	assert.Equal(t, createdFeatureVaultConnection.Type(), retrievedVaultConnection.Type())

	// sanity check - should be none by default
	retrievedFeature := retrievedVersionedSettings.(*teamcity.ProjectFeatureVersionedSettings)
	assert.Equal(t, 0, len(retrievedFeature.Options.ContextParameters))
}

func TestProjectFeature_GetByTypeAndProvider(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectConnectionVault(project.ID, teamcity.ConnectionProviderVaultOptions{
		DisplayName: "Hashicorp Vault",
		URL:         "http://vault.service:8200",
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)

	retrieved, err := service.GetByTypeAndProvider("OAuthProvider", "teamcity-vault")
	require.NoError(t, err)
	// we can't compare the ID, since they change after creation
	assert.Equal(t, createdFeature.Type(), retrieved.Type())

	// sanity check - should be none by default
	retrievedFeature := retrieved.(*teamcity.ConnectionProviderVault)
	assert.Equal(t, feature.Options.DisplayName, retrievedFeature.Options.DisplayName)
	assert.Equal(t, feature.Options.URL, retrievedFeature.Options.URL)
}

func setupFakeRoot(t *testing.T, client *teamcity.Client, project *teamcity.Project, name string) *teamcity.VcsRootReference {
	rootOptions, err := teamcity.NewGitVcsRootOptionsDefaults("master", "git@test.com")
	require.NoError(t, err)

	root, err := teamcity.NewGitVcsRoot(project.ID, name, rootOptions)
	require.NoError(t, err)

	createdRoot, err := client.VcsRoots.Create(project.ID, root)
	require.NoError(t, err)

	return createdRoot
}
