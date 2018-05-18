package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
)

func TestVcsRoot_Create(t *testing.T) {
	client := setup()
	newProject := getTestProjectData(testVcsRootProjectId)
	newVcsRoot := getTestVcsRootData(testVcsRootProjectId)

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

func TestVcsRoot_ValidateRequiredProperties(t *testing.T) {
	client := setup()
	newProject := getTestProjectData(testVcsRootProjectId)

	createdProject, err := client.Projects.Create(newProject)

	if err != nil {
		t.Fatalf("Failed to create project for VCS root: %s", err)
	}

	t.Run("VcsRoot must not be nil", func(t *testing.T) {
		var sut *teamcity.VcsRoot = nil

		_, err := client.VcsRoots.Create(createdProject.ID, sut)

		assert.NotNilf(t, err, "Expected error to be returned when vcsRoot is nil")
	})

	t.Run("Project must be specified", func(t *testing.T) {
		sut := getTestVcsRootData(testVcsRootProjectId)
		sut.Project = nil

		_, err := client.VcsRoots.Create(createdProject.ID, sut)

		assert.NotNilf(t, err, "Expected error to be returned when VcsRoot.Project property is not defined.")
	})

	t.Run("VcsName must be specified", func(t *testing.T) {
		sut := getTestVcsRootData(testVcsRootProjectId)
		sut.VcsName = ""

		_, err := client.VcsRoots.Create(createdProject.ID, sut)

		assert.NotNilf(t, err, "Expected error to be returned when VcsRoot.VcsName property is not defined.")
	})

	t.Run("Properties must have 'url' specified", func(t *testing.T) {
		sut := getTestVcsRootData(testVcsRootProjectId)
		sut.Properties = teamcity.NewProperties(
			&teamcity.Property{
				Name:  "someprop",
				Value: "empty",
			})

		_, err := client.VcsRoots.Create(createdProject.ID, sut)

		assert.EqualError(t, err, "'url' property must be defined in VcsRoot.Properties")
	})

	t.Run("Properties must have 'branch' specified", func(t *testing.T) {
		sut := getTestVcsRootData(testVcsRootProjectId)
		sut.Properties = teamcity.NewProperties(
			&teamcity.Property{
				Name:  "url",
				Value: "anything",
			})

		_, err := client.VcsRoots.Create(createdProject.ID, sut)

		assert.EqualError(t, err, "'branch' property must be defined in VcsRoot.Properties")
	})

	cleanUpProject(t, client, createdProject.ID)
}

func getTestVcsRootData(projectId string) *teamcity.VcsRoot {

	return &teamcity.VcsRoot{
		Name:    "Application",
		VcsName: teamcity.VcsNames.Git,
		Project: &teamcity.ProjectReference{
			ID: projectId,
		},
		Properties: teamcity.NewProperties(
			&teamcity.Property{
				Name:  "url",
				Value: "https://github.com/kelseyhightower/nocode",
			},
			&teamcity.Property{
				Name:  "branch",
				Value: "refs/head/master",
			}),
	}
}

func cleanUpVcsRoot(t *testing.T, c *teamcity.Client, id string) {
	if err := c.VcsRoots.Delete(id); err != nil {
		t.Errorf("Unable to delete vcs root with id = '%s', err: %s", id, err)
		return
	}

	deleted, err := c.VcsRoots.GetById(id)

	if deleted != nil {
		t.Errorf("Vcs root not deleted during cleanup.")
		return
	}

	if err == nil {
		t.Errorf("Expected 404 Not Found error when getting Deleted VCS Root, but no error returned.")
	}
}
