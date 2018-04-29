package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
)

func TestCreateProjectParameters(t *testing.T) {
	newProject := getTestProjectData(testParameterProjectId)
	client := setup()

	parameters := getTestParameters()

	expected := parameters.Map()
	_, err := client.Projects.Create(newProject)

	if err != nil {
		t.Fatalf("Error when creating project: %s", err)
	}

	err = client.ProjectParameterService(testParameterProjectId).Add(parameters.Items...)

	if err != nil {
		cleanUpProject(t, client, testParameterProjectId)
		t.Fatalf("Failed to create parameters: %s", err)
	}

	actual, err := client.Projects.GetById(testParameterProjectId)
	cleanUpProject(t, client, testParameterProjectId)

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

func TestCreateBuildTypeParameters(t *testing.T) {
	newProject := getTestProjectData(testParameterProjectId)
	newBuildType := getTestBuildTypeData("PullRequest", "Description", testParameterProjectId)
	client := setup()

	parameters := getTestParameters()
	expected := parameters.Map()
	_, err := client.Projects.Create(newProject)

	if err != nil {
		t.Fatalf("Error when creating project: %s", err)
	}

	created, err := client.BuildTypes.Create(testParameterProjectId, newBuildType)

	if err != nil {
		t.Fatalf("Error when creating build type: %s", err)
	}

	err = client.BuildTypeParameterService(created.ID).Add(parameters.Items...)

	if err != nil {
		cleanUpProject(t, client, testParameterProjectId)
		t.Fatalf("Failed to create parameters: %s", err)
	}

	actual, err := client.BuildTypes.GetById(created.ID)
	cleanUpProject(t, client, testParameterProjectId)

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

func getTestParameters() *teamcity.Properties {
	return teamcity.NewProperties(
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
}
