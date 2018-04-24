package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
)

func TestCreateBuildTypeForProject(t *testing.T) {
	client := setup()
	newProject := getTestProjectData()
	newBuildType := getTestBuildTypeData()

	createdProject, err := client.Projects.CreateProject(newProject)

	if err != nil {
		t.Fatalf("Failed to create project for buildType: %s", err)
	}

	newBuildType.ProjectID = createdProject.ID

	actual, err := client.CreateBuildType(createdProject.ID, newBuildType)

	if err != nil {
		t.Fatalf("Failed to CreateBuildType: %s", err)
	}

	if actual == nil {
		t.Fatalf("CreateBuildType did not return a valid project instance")
	}

	cleanUpBuildType(t, client, actual.ID)
	cleanUpProject(t, client, createdProject.ID)

	assert.NotEmpty(t, actual.ID)
	assert.Equal(t, newBuildType.ProjectID, actual.ProjectID)
	assert.Equal(t, newBuildType.Name, actual.Name)
}

func getTestBuildTypeData() *teamcity.BuildType {

	return &teamcity.BuildType{
		Name:        "Pull Request",
		Description: "Inspection Build",
	}
}

func cleanUpBuildType(t *testing.T, c *teamcity.Client, id string) {
	if err := c.DeleteBuildType(id); err != nil {
		t.Fatalf("Unable to delete build type with id = '%s', err: %s", id, err)
	}

	deletedProject, err := c.GetBuildType(id)

	if deletedProject != nil {
		t.Fatalf("Build type not deleted during cleanup.")
	}

	if err == nil {
		t.Fatalf("Expected 404 Not Found error when getting Deleted Build Type, but no error returned.")
	}
}
