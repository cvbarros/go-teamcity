package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	newProject := getTestProjectData()
	client := setup()
	actual, err := client.Projects.CreateProject(newProject)

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

func TestCreateProjectWithNoName(t *testing.T) {
	newProject := teamcity.Project{}
	client := setup()
	_, err := client.Projects.CreateProject(&newProject)

	assert.Equal(t, error.Error(err), "Project must have a name")
}

func getTestProjectData() *teamcity.Project {

	return &teamcity.Project{
		Name:        "Test Project",
		Description: "Test Project Description",
		Archived:    teamcity.NewFalse(),
	}
}

func cleanUpProject(t *testing.T, c *teamcity.Client, id string) {
	if err := c.DeleteProject(id); err != nil {
		t.Fatalf("Unable to delete project with id = '%s', err: %s", id, err)
	}

	deletedProject, err := c.GetProject(id)

	if deletedProject != nil {
		t.Fatalf("Project not deleted during cleanup.")
	}

	if err == nil {
		t.Fatalf("Expected 404 Not Found error when getting Deleted Project, but no error returned.")
	}
}
