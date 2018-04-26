package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	newProject := getTestProjectData("Project_Test")
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

func TestCreateProjectWithNoName(t *testing.T) {
	newProject := teamcity.Project{}
	client := setup()
	_, err := client.Projects.Create(&newProject)

	assert.Equal(t, error.Error(err), "Project must have a name")
}

func TestCreateProjectWithParameters(t *testing.T) {
	newProject := getTestProjectData("Project_Test")
	client := setup()

	parameters := teamcity.NewProperties(
		&teamcity.Property{
			Name:  "env.ENV_PARAMETER",
			Value: "env parameter value",
		},
		&teamcity.Property{
			Name:  "system.system_parameter",
			Value: "system parameter value",
		},
		&teamcity.Property{
			Name:  "configuration_parameter",
			Value: "configuration parameter value",
		})

	expected := parameters.Map()
	created, err := client.Projects.Create(newProject)

	if err != nil {
		t.Fatalf("Error when creating project: %s", err)
	}

	err = client.Projects.AddParameters(created.ID, parameters.Items...)

	if err != nil {
		cleanUpProject(t, client, created.ID)
		t.Fatalf("Failed to create parameters: %s", err)
	}

	actual, err := client.Projects.GetById(created.ID)
	cleanUpProject(t, client, created.ID)

	if err != nil || actual == nil {
		t.Fatalf("Error when getting created project: %s", err)
	}

	assert.NotNilf(t, actual.Parameters, "Expected parameters for project, but got nil")
	assert.NotEmpty(t, actual.Parameters.Items, "Expected parameters for project, but its empty")

	params := actual.Parameters.Map()
	for k, v := range expected {
		if value, ok := params[k]; !ok || value != v {
			t.Errorf("parameter '%s' was expected but was not defined or had incorrect value (actual: '%s', expected: %s)", k, value, v)
		}
	}
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
