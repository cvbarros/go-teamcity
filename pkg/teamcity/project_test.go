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

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.NotEmpty(t, actual.ID)

	cleanUpProject(t, client, actual.ID)

	assert.Equal(t, newProject.Name, actual.Name)
	assert.Equal(t, newProject.Description, actual.Description)
}

func TestProject_CreateWithParent(t *testing.T) {
	parent, _ := teamcity.NewProject(testProjectId, "Parent Project", "")
	child, _ := teamcity.NewProject("ChildProject", "Child Project", testProjectId)

	client := setup()

	_, err := client.Projects.Create(parent)
	require.NoError(t, err)
	created, err := client.Projects.Create(child)
	require.NoError(t, err)

	actual, _ := client.Projects.GetByID(created.ID) // Refresh
	cleanUpProject(t, client, testProjectId)

	assert.Equal(t, testProjectId, actual.ParentProjectID)
	require.NotNil(t, actual.ParentProject)
	assert.Equal(t, testProjectId, actual.ParentProject.ID)
}

func TestProject_UpdateWithSameParentDoesNotChangeName(t *testing.T) {
	parent, _ := teamcity.NewProject("ParentProject", "Parent Project", "")
	child, _ := teamcity.NewProject("ChildProject", "Child Project", "ParentProject")

	client := setup()

	_, err := client.Projects.Create(parent)
	require.NoError(t, err)
	created, err := client.Projects.Create(child)
	require.NoError(t, err)

	actual, _ := client.Projects.GetByID(created.ID) // Refresh
	updated, _ := client.Projects.Update(actual)     // Update with no change

	cleanUpProject(t, client, "ParentProject")

	assert.Equal(t, "ChildProject", updated.Name)
}

func TestProject_UpdateParent(t *testing.T) {
	parent, _ := teamcity.NewProject(testProjectId, "Parent Project", "")
	newParent, _ := teamcity.NewProject("NewParent", "NewParent Project", "")
	child, _ := teamcity.NewProject("ChildProject", "Child Project", testProjectId)

	client := setup()

	_, err := client.Projects.Create(parent)
	require.NoError(t, err)
	createdParent, err := client.Projects.Create(newParent)
	require.NoError(t, err)
	created, err := client.Projects.Create(child)
	require.NoError(t, err)

	actual, _ := client.Projects.GetByID(created.ID) // Refresh

	actual.SetParentProject(createdParent.ID)

	actual, err = client.Projects.Update(actual)

	cleanUpProject(t, client, testProjectId)
	cleanUpProject(t, client, createdParent.ID)

	assert.Equal(t, createdParent.ID, actual.ParentProjectID)
	require.NotNil(t, actual.ParentProject)
	assert.Equal(t, createdParent.ID, actual.ParentProject.ID)
}

func TestProject_UpdateParameters(t *testing.T) {
	client := setup()
	pa := newPropertyAssertions(t)
	created := createTestProject(t, client, testProjectId)
	sut := client.Projects

	actual, err := sut.GetByID(created.ID) //Refresh

	//Update some fields
	props := teamcity.NewParametersEmpty()
	props.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "param1", "value1")
	props.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "param2", "value2")
	actual.Parameters = props

	updated, err := sut.Update(actual)
	cleanUpProject(t, client, testProjectId)

	require.NoError(t, err)
	pa.assertPropertyValue(updated.Parameters.Properties(), "param1", "value1")
	pa.assertPropertyValue(updated.Parameters.Properties(), "param2", "value2")
}

func TestProject_UpdateParametersWithRemoval(t *testing.T) {
	client := setup()
	pa := newPropertyAssertions(t)
	created := createTestProject(t, client, testProjectId)
	sut := client.Projects

	actual, err := sut.GetByID(created.ID) //Refresh

	params := teamcity.NewParametersEmpty()
	params.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "param1", "value1")
	params.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "param2", "value2")
	actual.Parameters = params
	actual, err = sut.Update(actual)

	actual.Parameters.Remove(teamcity.ParameterTypes.Configuration, "param2")
	actual, err = sut.Update(actual)
	cleanUpProject(t, client, testProjectId)

	require.NoError(t, err)

	pa.assertPropertyValue(actual.Parameters.Properties(), "param1", "value1")
	pa.assertPropertyDoesNotExist(actual.Parameters.Properties(), "param2")
}

func TestProject_GetByName(t *testing.T) {
	client := setup()
	created := createTestProject(t, client, testProjectId)
	sut := client.Projects

	actual, err := sut.GetByName(created.Name)
	cleanUpProject(t, client, testProjectId)
	require.NoError(t, err)

	assert.Equal(t, created.Name, actual.Name)
}

func TestProject_GetRootByName(t *testing.T) {
	client := setup()
	sut := client.Projects

	actual, err := sut.GetByName("<Root project>") // This is case sensitive!
	require.NoError(t, err)

	assert.Equal(t, "<Root project>", actual.Name)
}

func TestProject_BuildTypes(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	bt1 := createTestBuildTypeWithName(t, client, testProjectId, "Build1", true)
	bt2 := createTestBuildTypeWithName(t, client, testProjectId, "Build2", false)

	actual, _ := client.Projects.GetByID(testProjectId)

	cleanUpProject(t, client, testProjectId)
	assert.Equal(int32(2), actual.BuildTypes.Count)

	projectBuildTypes := make(map[string]string, 2)
	for _, i := range actual.BuildTypes.Items {
		projectBuildTypes[i.ID] = i.Name
	}

	assert.Contains(projectBuildTypes, bt1.ID)
	assert.Contains(projectBuildTypes, bt2.ID)
}

func TestProject_ValidateName(t *testing.T) {
	_, err := teamcity.NewProject("", "", "")

	require.Errorf(t, err, "name is required")
}

func TestProject_GetUnauthorizedHandled(t *testing.T) {
	client, _ := teamcity.New("admin", "error", http.DefaultClient)
	_, err := client.Projects.Create(getTestProjectData(testProjectId))

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "401")
}

func getTestProjectData(name string) *teamcity.Project {
	out, _ := teamcity.NewProject(name, "Test Project Description", "")
	return out
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

func createTestProject(t *testing.T, c *teamcity.Client, name string) *teamcity.Project {
	newProject := getTestProjectData(name)
	createdProject, err := c.Projects.Create(newProject)

	if err != nil {
		t.Fatalf("Failed to create project for VCS root: %s", err)
	}

	return createdProject
}
