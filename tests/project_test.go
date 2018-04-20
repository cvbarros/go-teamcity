package tests

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk"
	"github.com/stretchr/testify/assert"
)

func init() {
	client = initTest()
}

func TestCreateProject(t *testing.T) {
	newProject := getTestProjectData()

	actual, err := client.CreateProject(newProject)

	if err != nil {
		t.Fatalf("Failed to GetServer: %s", err)
	}

	if actual == nil {
		t.Fatalf("CreateProject did not return a valid project instance")
	}

	cleanUpProject(t, actual.ID)

	assert.NotEmpty(t, actual.ID)
	assert.Equal(t, newProject.Name, actual.Name)
}

func TestCreateProjectWithNoName(t *testing.T) {
	newProject := teamcity.Project{}

	_, err := client.CreateProject(&newProject)

	assert.Equal(t, error.Error(err), "Project must have a name")
}

func getTestProjectData() *teamcity.Project {

	return &teamcity.Project{
		Name:        "Test Project",
		Description: "Test Project Description",
		Archived:    newFalse(),
	}
}

func cleanUpProject(t *testing.T, id string) {
	if err := client.DeleteProject(id); err != nil {
		t.Fatalf("Unable to delete project with id = '%s', err: %s", id, err)
	}

	deletedProject, err := client.GetProject(id)

	if deletedProject != nil {
		t.Fatalf("Project not deleted during cleanup.")
	}

	if err == nil {
		t.Fatalf("Expected 404 Not Found error when getting Deleted Project, but no error returned.")
	}
}
