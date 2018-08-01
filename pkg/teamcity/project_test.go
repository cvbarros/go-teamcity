package teamcity_test

import (
	"net/http"
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProject_Create(t *testing.T) {
	newProject := getTestProjectData(testProjectId)
	client := setup()
	actual, err := client.Projects.Create(newProject)

	if err != nil {
		t.Fatalf("Failed to GetServer: %s", err)
	}

	if actual == nil {
		t.Fatalf("CreateProject did not return a valid project instance")
	}

	assert.NotEmpty(t, actual.ID)

	cleanUpProject(t, client, actual.ID)

	assert.Equal(t, newProject.Name, actual.Name)
}

func TestProject_ValidateName(t *testing.T) {
	newProject := teamcity.Project{}
	client := setup()
	_, err := client.Projects.Create(&newProject)

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "Project name cannot be empty")
}

func TestProject_GetUnauthorizedHandled(t *testing.T) {
	client, _ := teamcity.New("admin", "error", http.DefaultClient)
	_, err := client.Projects.Create(getTestProjectData(testProjectId))

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "401")
}

func getTestProjectData(name string) *teamcity.Project {

	return &teamcity.Project{
		Name:        name,
		Description: "Test Project Description",
		Archived:    teamcity.NewFalse(),
	}
}

func cleanUpProject(t *testing.T, c *teamcity.Client, id string) {
	if err := c.Projects.Delete(id); err != nil {
		t.Fatalf("Unable to delete project with id = '%s', err: %s", id, err)
	}

	deletedProject, err := c.Projects.GetByID(id)

	if deletedProject != nil {
		t.Fatalf("Project not deleted during cleanup.")
	}

	if err == nil {
		t.Fatalf("Expected 404 Not Found error when getting Deleted Project, but no error returned.")
	}
}

func createTestProject(t *testing.T, c *teamcity.Client, name string) *teamcity.ProjectReference {
	newProject := getTestProjectData(name)
	createdProject, err := c.Projects.Create(newProject)

	if err != nil {
		t.Fatalf("Failed to create project for VCS root: %s", err)
	}

	return createdProject
}
