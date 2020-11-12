package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectFeature_GoogleCloudProfile_Create(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureGoogleCloudProfile(project.ID, teamcity.ProjectFeatureGoogleCloudProfileOptions{
		Enabled:           true,
		Name:              "Test",
		CredentialsType:   "key",
		TerminateIdleTime: 20,
		TotalWorkTime:     20,
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)
}

func TestProjectFeature_GoogleCloudProfile_Delete(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureGoogleCloudProfile(project.ID, teamcity.ProjectFeatureGoogleCloudProfileOptions{
		Enabled:           true,
		Name:              "Test",
		CredentialsType:   "key",
		TerminateIdleTime: 20,
		TotalWorkTime:     20,
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

func TestProjectFeature_GoogleCloudProfile_Update(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureGoogleCloudProfile(project.ID, teamcity.ProjectFeatureGoogleCloudProfileOptions{
		Enabled:         true,
		Name:            "Test Profile",
		CredentialsType: "key",
		AccessKey:       "",
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)

	type testData = struct {
		Enabled             bool
		ProfileID           string
		Name                string
		Description         string
		CloudCode           string
		ProfileServerURL    string
		AgentPushPreset     string
		TotalWorkTime       int
		CredentialsType     string
		NextHour            string
		TerminateAfterBuild bool
		TerminateIdleTime   int
		AccessKey           string
	}

	var validate = func(t *testing.T, id string, data testData) {
		retrievedFeature, err := service.GetByID(id)
		require.NoError(t, err)
		cloudProfile, ok := retrievedFeature.(*teamcity.ProjectFeatureGoogleCloudProfile)
		assert.True(t, ok)

		assert.Equal(t, data.Enabled, cloudProfile.Options.Enabled)
		assert.Equal(t, data.Name, cloudProfile.Options.Name)
		assert.Equal(t, data.Description, cloudProfile.Options.Description)
		assert.Equal(t, data.CloudCode, cloudProfile.Options.CloudCode)
		assert.Equal(t, data.ProfileServerURL, cloudProfile.Options.ProfileServerURL)
		assert.Equal(t, data.AgentPushPreset, cloudProfile.Options.AgentPushPreset)
		assert.Equal(t, data.TotalWorkTime, cloudProfile.Options.TotalWorkTime)
		assert.Equal(t, data.CredentialsType, cloudProfile.Options.CredentialsType)
		assert.Equal(t, data.NextHour, cloudProfile.Options.NextHour)
		assert.Equal(t, data.TerminateAfterBuild, cloudProfile.Options.TerminateAfterBuild)
		assert.Equal(t, data.TerminateIdleTime, cloudProfile.Options.TerminateIdleTime)
		assert.Equal(t, data.AccessKey, cloudProfile.Options.AccessKey)
	}
	t.Log("Validating initial creation")
	validate(t, createdFeature.ID(), testData{
		Enabled:         true,
		Name:            "Test Profile",
		CloudCode:       "google",
		CredentialsType: "key",
	})

	updateConfigurations := []testData{
		{
			Enabled:     false,
			Name:        "Test Profile - Updated",
			Description: "Changed Description / Enabled / Name",
		},
		{
			Enabled:       true,
			Name:          "Test Profile",
			Description:   "Updating TotalWorkTime",
			TotalWorkTime: 100,
		},
		{
			Enabled:             true,
			Name:                "Test Profile",
			Description:         "Updating TerminateAfterBuild",
			TerminateAfterBuild: true,
		},
	}
	for _, update := range updateConfigurations {
		t.Logf("Testing %q", update.Description)

		existing, err := service.GetByID(createdFeature.ID())
		require.NoError(t, err)

		settings, ok := existing.(*teamcity.ProjectFeatureGoogleCloudProfile)
		assert.True(t, ok)

		update.ProfileID = settings.Options.ProfileID
		update.CloudCode = settings.Options.CloudCode

		settings.Options.Enabled = update.Enabled
		settings.Options.Name = update.Name
		settings.Options.Description = update.Description
		settings.Options.ProfileServerURL = update.ProfileServerURL
		settings.Options.AgentPushPreset = update.AgentPushPreset
		settings.Options.TotalWorkTime = update.TotalWorkTime
		settings.Options.CredentialsType = update.CredentialsType
		settings.Options.NextHour = update.NextHour
		settings.Options.TerminateAfterBuild = update.TerminateAfterBuild
		settings.Options.TerminateIdleTime = update.TerminateIdleTime
		settings.Options.AccessKey = update.AccessKey

		updatedFeature, err := service.Update(settings)
		require.NoError(t, err)
		assert.NotEmpty(t, updatedFeature.ID)

		// sanity check since we're updating with the same ID
		assert.Equal(t, createdFeature.ID(), updatedFeature.ID())

		validate(t, updatedFeature.ID(), update)
	}
}
