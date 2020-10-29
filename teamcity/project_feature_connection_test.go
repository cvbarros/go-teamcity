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
