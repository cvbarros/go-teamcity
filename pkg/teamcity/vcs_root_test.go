package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVcsRoot_Create(t *testing.T) {
	client := setup()
	newProject := getTestProjectData(testVcsRootProjectId)
	newVcsRoot := getTestVcsRootData(testVcsRootProjectId).(*teamcity.GitVcsRoot)

	createdProject, err := client.Projects.Create(newProject)

	if err != nil {
		t.Fatalf("Failed to create project for VCS root: %s", err)
	}

	newVcsRoot.Project.ID = createdProject.ID

	actual, err := client.VcsRoots.Create(createdProject.ID, newVcsRoot)

	if err != nil {
		t.Fatalf("Failed to create VCS Root: %s", err)
	}

	if actual == nil {
		t.Fatalf("Create did not return a valid VCS root instance")
	}

	cleanUpVcsRoot(t, client, actual.ID)
	cleanUpProject(t, client, createdProject.ID)

	assert.NotEmpty(t, actual.ID)
	assert.Equal(t, newVcsRoot.Project.ID, actual.Project.ID)
	assert.Equal(t, newVcsRoot.Name, actual.Name)
}

func TestVcsRoot_CreateWithUsernamePassword(t *testing.T) {
	client := setup()
	newProject := getTestProjectData(testVcsRootProjectId)
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

	cleanUpVcsRoot(t, client, actual.ID)
	cleanUpProject(t, client, createdProject.ID)

	props := created.Properties()
	propAssert := newPropertyAssertions(t)

	propAssert.assertPropertyValue(props, "authMethod", string(teamcity.GitAuthMethodPassword))
	propAssert.assertPropertyValue(props, "username", "admin")
	propAssert.assertPropertyExists(props, "secure:password")
}

func TestVcsRoot_Invariants(t *testing.T) {
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

func cleanUpVcsRoot(t *testing.T, c *teamcity.Client, id string) {
	if err := c.VcsRoots.Delete(id); err != nil {
		t.Errorf("Unable to delete vcs root with id = '%s', err: %s", id, err)
		return
	}

	deleted, err := c.VcsRoots.GetByID(id)

	if deleted != nil {
		t.Errorf("Vcs root not deleted during cleanup.")
		return
	}

	if err == nil {
		t.Errorf("Expected 404 Not Found error when getting Deleted VCS Root, but no error returned.")
	}
}

func ExampleVcsRoot_VcsName() {
	var vcsRoot teamcity.VcsRoot
	//Retrieve vcsRoot from API, for instance
	//Check for its type
	switch vcsRoot.VcsName() {
	case teamcity.VcsNames.Git:
		git := vcsRoot.(*teamcity.GitVcsRoot)
		//Use it strongly-typed
		println(git.Options.FetchURL)
	}
}
