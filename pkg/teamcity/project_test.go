package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, error.Error(err), "Project must have a name")
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

	deletedProject, err := c.Projects.GetById(id)

	if deletedProject != nil {
		t.Fatalf("Project not deleted during cleanup.")
	}

	if err == nil {
		t.Fatalf("Expected 404 Not Found error when getting Deleted Project, but no error returned.")
	}
}
