package teamcity_test

import (
	"testing"

	u "github.com/cvbarros/go-teamcity-sdk/internal/testutil"
	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
)

func init() {
	u.InitTest()
}

func TestCreateProject(t *testing.T) {
	newProject := getTestProjectData()

	actual, err := u.Client.CreateProject(newProject)

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

	_, err := u.Client.CreateProject(&newProject)

	assert.Equal(t, error.Error(err), "Project must have a name")
}

func getTestProjectData() *teamcity.Project {

	return &teamcity.Project{
		Name:        "Test Project",
		Description: "Test Project Description",
		Archived:    teamcity.NewFalse(),
	}
}

func cleanUpProject(t *testing.T, id string) {
	if err := u.Client.DeleteProject(id); err != nil {
		t.Fatalf("Unable to delete project with id = '%s', err: %s", id, err)
	}

	deletedProject, err := u.Client.GetProject(id)

	if deletedProject != nil {
		t.Fatalf("Project not deleted during cleanup.")
	}

	if err == nil {
		t.Fatalf("Expected 404 Not Found error when getting Deleted Project, but no error returned.")
	}
}
