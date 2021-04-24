package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectFeatureConnection_CreateVault(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectConnectionVault(project.ID, teamcity.ConnectionProviderVaultOptions{
		AuthMethod:  "approle",
		DisplayName: "Hashicorp Vault",
		Endpoint:    "approle",
		RoleID:      "123456",
		SecretID:    "abcdefg",
		URL:         "http://vault.service:8200",
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID())
}

func TestProjectFeatureConnection_DeleteVault(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectConnectionVault(project.ID, teamcity.ConnectionProviderVaultOptions{
		AuthMethod:  "approle",
		DisplayName: "Hashicorp Vault",
		Endpoint:    "approle",
		RoleID:      "123456",
		SecretID:    "abcdefg",
		URL:         "http://vault.service:8200",
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID())

	err = service.Delete(createdFeature.ID())
	require.NoError(t, err)

	deletedFeature, err := service.GetByID(createdFeature.ID())
	assert.NotNil(t, err)
	assert.Nil(t, deletedFeature)
}

func TestProjectFeatureConnection_UpdateVault(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectConnectionVault(project.ID, teamcity.ConnectionProviderVaultOptions{
		AuthMethod:  "approle",
		DisplayName: "Hashicorp Vault",
		Endpoint:    "approle",
		RoleID:      "123456",
		SecretID:    "abcdefg",
		URL:         "http://vault.service:8200",
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID())

	type testData = struct {
		description string

		authMethod     teamcity.ConnectionProviderVaultAuthMethod
		displayName    string
		failOnError    bool
		namespace      string
		roleID         string
		secretID       string
		URL            string
		vaultNamespace string
	}

	var validate = func(t *testing.T, id string, data testData) {
		retrievedFeature, err := service.GetByID(id)
		require.NoError(t, err)

		vaultConnectionSettings, ok := retrievedFeature.(*teamcity.ConnectionProviderVault)
		assert.True(t, ok)

		assert.Equal(t, data.authMethod, vaultConnectionSettings.Options.AuthMethod)
		assert.Equal(t, data.displayName, vaultConnectionSettings.Options.DisplayName)
		assert.Equal(t, data.failOnError, vaultConnectionSettings.Options.FailOnError)
		assert.Equal(t, data.namespace, vaultConnectionSettings.Options.Namespace)
		assert.Equal(t, data.roleID, vaultConnectionSettings.Options.RoleID)
		assert.Equal(t, data.URL, vaultConnectionSettings.Options.URL)
		assert.Equal(t, data.vaultNamespace, vaultConnectionSettings.Options.VaultNamespace)
	}

	t.Log("Validating initial connection creation")
	validate(t, createdFeature.ID(), testData{
		authMethod:  teamcity.ConnectionProviderVaultAuthMethodApprole,
		displayName: "Hashicorp Vault",
		failOnError: false,
		roleID:      "123456",
		URL:         "http://vault.service:8200",
	})

	updateConfigurations := []testData{
		{
			description: "Update displayName",
			authMethod:  teamcity.ConnectionProviderVaultAuthMethodApprole,
			displayName: "Hashicorp Vault test",
			failOnError: false,
			roleID:      "123456",
			URL:         "http://vault.service:8200",
		},
		{
			description: "Update failOnError",
			authMethod:  teamcity.ConnectionProviderVaultAuthMethodApprole,
			displayName: "Hashicorp Vault",
			failOnError: true,
			roleID:      "123456",
			URL:         "http://vault.service:8200",
		},
		{
			description: "Update roleID",
			authMethod:  teamcity.ConnectionProviderVaultAuthMethodApprole,
			displayName: "Hashicorp Vault",
			failOnError: true,
			roleID:      "567890",
			URL:         "http://vault.service:8200",
		},
		{
			description: "Update URL",
			authMethod:  teamcity.ConnectionProviderVaultAuthMethodApprole,
			displayName: "Hashicorp Vault",
			failOnError: true,
			roleID:      "567890",
			URL:         "http://vault.differet-service:8200",
		},
		{
			description: "Update authMethod",
			authMethod:  teamcity.ConnectionProviderVaultAuthMethodIAM,
			displayName: "Hashicorp Vault",
			failOnError: true,
			roleID:      "567890",
			URL:         "http://vault.differet-service:8200",
		},
		{
			description: "Update namespace",
			authMethod:  teamcity.ConnectionProviderVaultAuthMethodIAM,
			displayName: "Hashicorp Vault",
			failOnError: true,
			namespace:   "test",
			roleID:      "567890",
			URL:         "http://vault.differet-service:8200",
		},
	}

	for _, update := range updateConfigurations {
		t.Logf("Testing %q", update.description)

		existing, err := service.GetByID(createdFeature.ID())
		require.NoError(t, err)

		settings, ok := existing.(*teamcity.ConnectionProviderVault)
		assert.True(t, ok)

		settings.Options.AuthMethod = update.authMethod
		settings.Options.DisplayName = update.displayName
		// settings.Options.Endpoint = update.endpoint
		settings.Options.FailOnError = update.failOnError
		settings.Options.Namespace = update.namespace
		settings.Options.RoleID = update.roleID
		settings.Options.URL = update.URL
		settings.Options.VaultNamespace = update.vaultNamespace
		// settings.Options.CredentialsStorageType = update.credentialsType

		updatedFeature, err := service.Update(settings)
		require.NoError(t, err)
		assert.NotEmpty(t, updatedFeature.ID)

		// sanity check since we're updating with the same ID
		assert.Equal(t, createdFeature.ID(), updatedFeature.ID())

		validate(t, updatedFeature.ID(), update)
	}
}
