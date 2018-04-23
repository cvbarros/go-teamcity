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

func TestCreateBuildTypeForProject(t *testing.T) {
	newProject := getTestProjectData()
	newBuildType := getTestBuildTypeData()

	createdProject, err := u.Client.CreateProject(newProject)

	if err != nil {
		t.Fatalf("Failed to create project for buildType: %s", err)
	}

	newBuildType.ProjectID = createdProject.ID

	actual, err := u.Client.CreateBuildType(createdProject.ID, newBuildType)

	if err != nil {
		t.Fatalf("Failed to CreateBuildType: %s", err)
	}

	if actual == nil {
		t.Fatalf("CreateBuildType did not return a valid project instance")
	}

	cleanUpBuildType(t, actual.ID)
	cleanUpProject(t, createdProject.ID)

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

func cleanUpBuildType(t *testing.T, id string) {
	if err := u.Client.DeleteBuildType(id); err != nil {
		t.Fatalf("Unable to delete build type with id = '%s', err: %s", id, err)
	}

	deletedProject, err := u.Client.GetBuildType(id)

	if deletedProject != nil {
		t.Fatalf("Build type not deleted during cleanup.")
	}

	if err == nil {
		t.Fatalf("Expected 404 Not Found error when getting Deleted Build Type, but no error returned.")
	}
}
