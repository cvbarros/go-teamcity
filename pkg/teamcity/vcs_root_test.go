package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitVcsRoot_Get(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	client := setup()
	newProject := createTestProject(t, client, testVcsRootProjectId)
	newVcsRoot := getTestVcsRootData(testVcsRootProjectId).(*teamcity.GitVcsRoot)
	sut := client.VcsRoots

	created, err := sut.Create(newProject.ID, newVcsRoot)

	require.NoError(err)
	require.NotNil(created)

	data, err := sut.GetByID(created.ID)
	require.NoError(err)
	require.NotNil(data)
	require.IsType(&teamcity.GitVcsRoot{}, data)

	actual := data.(*teamcity.GitVcsRoot)
	cleanUpProject(t, client, newProject.ID)

	require.NotNil(actual.Project)
	assert.Equal(actual.Project.ID, actual.Project.ID)
	assert.Equal(actual.Name(), actual.Name())
}

func TestGitVcsRoot_Update(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	client := setup()
	newProject := createTestProject(t, client, testVcsRootProjectId)
	updatedProject := createTestProject(t, client, "VcsRootUpdatedProjectTest")
	newVcsRoot := getTestVcsRootData(testVcsRootProjectId).(*teamcity.GitVcsRoot)
	sut := client.VcsRoots

	created, err := sut.Create(newProject.ID, newVcsRoot)

	require.NoError(err)
	require.NotNil(created)

	data, err := sut.GetByID(created.ID) // refresh
	require.NoError(err)
	require.NotNil(data)
	require.IsType(&teamcity.GitVcsRoot{}, data)

	gitVcs := data.(*teamcity.GitVcsRoot)
	as := &teamcity.GitAgentSettings{
		CleanFilesPolicy: teamcity.CleanFilesPolicyIgnoredUntracked,
		CleanPolicy:      teamcity.CleanPolicyNever,
		UseMirrors:       false,
	}
	opt, _ := teamcity.NewGitVcsRootOptionsWithAgentSettings(
		"refs/head/develop",
		"https://github.com/cvbarros/go-teamcity-sdk",
		"",
		teamcity.GitAuthMethodAnonymous,
		"",
		"",
		as)

	gitVcs.Options = opt
	opt.SubModuleCheckout = "IGNORE"
	gitVcs.SetModificationCheckInterval(60)
	gitVcs.SetName("new_name")
	gitVcs.Project.ID = updatedProject.ID

	data, err = sut.Update(gitVcs)
	cleanUpProject(t, client, newProject.ID)
	cleanUpProject(t, client, updatedProject.ID)

	require.NoError(err)
	require.NotNil(data)
	require.IsType(&teamcity.GitVcsRoot{}, data)
	actual := data.(*teamcity.GitVcsRoot)

	assert.Equal("new_name", actual.Name())
	assert.Equal(int32(60), *(actual.ModificationCheckInterval()))
	assert.Equal(updatedProject.ID, actual.Project.ID)

	actualOpt := actual.Options
	assert.Equal(teamcity.GitAuthMethodAnonymous, actualOpt.AuthMethod)
	assert.Equal("https://github.com/cvbarros/go-teamcity-sdk", actualOpt.FetchURL)
	assert.Equal("refs/head/develop", actualOpt.DefaultBranch)
	assert.Equal("IGNORE", actualOpt.SubModuleCheckout)
	assert.Equal(teamcity.CleanFilesPolicyIgnoredUntracked, actualOpt.AgentSettings.CleanFilesPolicy)
	assert.Equal(teamcity.CleanPolicyNever, actualOpt.AgentSettings.CleanPolicy)
	assert.Equal(false, actualOpt.AgentSettings.UseMirrors)
}

func TestGitVcsRoot_CreateWithUsernamePassword(t *testing.T) {
	client := setup()
	newProject := getTestProjectData(testVcsRootProjectId, "")
	opts, _ := teamcity.NewGitVcsRootOptions("refs/head/master", "https://github.com/cvbarros/go-teamcity-sdk/", "", teamcity.GitAuthMethodPassword, "admin", "admin")

	createdProject, err := client.Projects.Create(newProject)

	if err != nil {
		t.Fatalf("Failed to create project for VCS root: %s", err)
	}

	newVcsRoot, _ := teamcity.NewGitVcsRoot(createdProject.ID, "Application", opts)

	actual, err := client.VcsRoots.Create(createdProject.ID, newVcsRoot)

	if err != nil {
		t.Fatalf("Failed to create VCS Root: %s", err)
	}

	if actual == nil {
		t.Fatalf("Create did not return a valid VCS root instance")
	}

	created, err := client.VcsRoots.GetByID(actual.ID)

	require.NoError(t, err)

	cleanUpProject(t, client, createdProject.ID)

	props := created.Properties()
	propAssert := newPropertyAssertions(t)

	propAssert.assertPropertyValue(props, "authMethod", string(teamcity.GitAuthMethodPassword))
	propAssert.assertPropertyValue(props, "username", "admin")
	propAssert.assertPropertyExists(props, "secure:password")
}

func TestGitVcsRoot_Invariants(t *testing.T) {
	gitOpt, _ := teamcity.NewGitVcsRootOptionsDefaults("master", "https://github.com/cvbarros/go-teamcity-sdk/")
	t.Run("projectID is required", func(t *testing.T) {
		_, err := teamcity.NewGitVcsRoot("", "name", gitOpt)
		require.EqualError(t, err, "projectID is required")
	})
	t.Run("name is required", func(t *testing.T) {
		_, err := teamcity.NewGitVcsRoot("project1", "", gitOpt)
		require.EqualError(t, err, "name is required")
	})
	t.Run("opts is required", func(t *testing.T) {
		_, err := teamcity.NewGitVcsRoot("project1", "name", nil)
		require.EqualError(t, err, "opts is required")
	})
}

func getTestVcsRootData(projectId string) teamcity.VcsRoot {
	opts, _ := teamcity.NewGitVcsRootOptionsDefaults("refs/head/master", "https://github.com/cvbarros/go-teamcity-sdk")
	v, _ := teamcity.NewGitVcsRoot(projectId, "Application", opts)
	return v
}
