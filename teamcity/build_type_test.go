package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildType_Create(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	actual := createTestBuildTypeWithName(t, client, "BuildTypeProject", "BuildRelease", true)

	cleanUpProject(t, client, "BuildTypeProject")

	assert.NotEmpty(t, actual.ID)
	assert.Equal("BuildTypeProject_BuildRelease", actual.ID)
	assert.Equal("BuildRelease", actual.Name)
	assert.Equal("BuildTypeProject", actual.ProjectID)
	assert.Equal(false, actual.IsTemplate)

	//Verify some default properties
	optExpected := teamcity.NewBuildTypeOptionsWithDefaults()
	optActual := actual.Options
	require.NotNil(t, optActual)
	assert.Equal(optExpected.ArtifactRules, optActual.ArtifactRules)
	assert.Equal(optExpected.BuildCounter, optActual.BuildCounter)
	assert.Equal(optExpected.BuildNumberFormat, optActual.BuildNumberFormat)
	assert.Equal(optExpected.AllowPersonalBuildTriggering, optActual.AllowPersonalBuildTriggering)
}

func TestBuildType_CreateWithChildProject(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	parent := createTestProject(t, client, testBuildTypeProjectId)
	child := createTestProjectWithParent(t, client, "BuildTypeProjectChild", parent.ID)

	actual := createTestBuildTypeWithName(t, client, child.ID, "BuildRelease", false)

	cleanUpProject(t, client, testBuildTypeProjectId)

	assert.NotEmpty(t, actual.ID)
	assert.Equal("BuildTypeProjectTest_BuildTypeProjectChild_BuildRelease", actual.ID)
	assert.Equal("BuildRelease", actual.Name)
	assert.Equal("BuildTypeProjectTest_BuildTypeProjectChild", actual.ProjectID)
	assert.Equal(false, actual.IsTemplate)

	//Verify some default properties
	optExpected := teamcity.NewBuildTypeOptionsWithDefaults()
	optActual := actual.Options
	require.NotNil(t, optActual)
	assert.Equal(optExpected.ArtifactRules, optActual.ArtifactRules)
	assert.Equal(optExpected.BuildCounter, optActual.BuildCounter)
	assert.Equal(optExpected.BuildNumberFormat, optActual.BuildNumberFormat)
	assert.Equal(optExpected.AllowPersonalBuildTriggering, optActual.AllowPersonalBuildTriggering)
}

func TestBuildTypeTemplate_Create(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	actual := createTestBuildTypeTemplateWithName(t, client, "BuildTypeProject", "BuildRelease", true)

	cleanUpProject(t, client, "BuildTypeProject")

	assert.NotEmpty(t, actual.ID)
	assert.Equal("BuildRelease", actual.Name)
	assert.Equal("BuildTypeProject", actual.ProjectID)
	assert.Equal(true, actual.IsTemplate)

	//Verify some default properties
	optExpected := teamcity.NewBuildTypeOptionsWithDefaults()
	optActual := actual.Options
	require.NotNil(t, optActual)
	assert.Equal(optExpected.ArtifactRules, optActual.ArtifactRules)
	assert.Equal(optExpected.BuildNumberFormat, optActual.BuildNumberFormat)
	assert.Equal(optExpected.AllowPersonalBuildTriggering, optActual.AllowPersonalBuildTriggering)
}

func TestBuildType_Get(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	created := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, testBuildTypeId, true)

	actual, err := client.BuildTypes.GetByID(created.ID)
	cleanUpProject(t, client, testBuildTypeProjectId)

	require.NoError(t, err)

	assert.Equal(created.ProjectID, actual.ProjectID)
	assert.Equal(created.Name, actual.Name)
}

func TestBuildType_Update(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	created := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, testBuildTypeId, true)
	sut := client.BuildTypes

	actual, err := sut.GetByID(created.ID) //Refresh

	//Update some fields
	actual.Description = "Updated description"
	actual.Options.ArtifactRules = []string{"rule1", "rule2"}
	actual.Options.BuildCounter = 10
	actual.Options.AllowPersonalBuildTriggering = false

	updated, err := sut.Update(actual)
	cleanUpProject(t, client, testBuildTypeProjectId)

	require.NoError(t, err)

	assert.Equal("Updated description", updated.Description)
	assert.Equal([]string{"rule1", "rule2"}, updated.Options.ArtifactRules)
	assert.Equal(10, updated.Options.BuildCounter)
	assert.Equal(false, updated.Options.AllowPersonalBuildTriggering)
}

func TestBuildType_CreateWithParameters(t *testing.T) {
	client := setup()
	pa := newPropertyAssertions(t)
	project := createTestProject(t, client, testBuildTypeProjectId)
	bt, _ := teamcity.NewBuildType(project.ID, "testbuild")

	bt.Parameters.AddOrReplaceValue(teamcity.ParameterTypes.EnvironmentVariable, "param1", "value1")
	bt.Parameters.AddOrReplaceValue(teamcity.ParameterTypes.System, "param2", "value2")

	sut := client.BuildTypes
	created, err := sut.Create(project.ID, bt)
	require.NoError(t, err)

	actual, _ := sut.GetByID(created.ID) //Refresh
	props := actual.Parameters.Properties()
	pa.assertPropertyValue(props, "env.param1", "value1")
	pa.assertPropertyValue(props, "system.param2", "value2")

	cleanUpProject(t, client, project.ID)
}

func TestBuildType_UpdateParameters(t *testing.T) {
	client := setup()
	pa := newPropertyAssertions(t)
	created := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, testBuildTypeId, true)
	sut := client.BuildTypes

	actual, err := sut.GetByID(created.ID) //Refresh

	//Update some fields
	props := teamcity.NewParametersEmpty()
	props.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "param1", "value1")
	props.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "param2", "value2")
	actual.Parameters = props

	updated, err := sut.Update(actual)
	cleanUpProject(t, client, testBuildTypeProjectId)

	require.NoError(t, err)
	pa.assertPropertyValue(updated.Parameters.Properties(), "param1", "value1")
	pa.assertPropertyValue(updated.Parameters.Properties(), "param2", "value2")
}

func TestBuildType_UpdateParametersWithRemoval(t *testing.T) {
	client := setup()
	pa := newPropertyAssertions(t)
	created := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, testBuildTypeId, true)
	sut := client.BuildTypes

	actual, err := sut.GetByID(created.ID) //Refresh

	props := teamcity.NewParametersEmpty()
	props.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "param1", "value1")
	props.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "param2", "value2")
	actual.Parameters = props
	actual, err = sut.Update(actual)

	actual.Parameters.Remove(teamcity.ParameterTypes.Configuration, "param2")
	actual, err = sut.Update(actual)
	cleanUpProject(t, client, testBuildTypeProjectId)

	require.NoError(t, err)

	pa.assertPropertyValue(actual.Parameters.Properties(), "param1", "value1")
	pa.assertPropertyDoesNotExist(actual.Parameters.Properties(), "param2")
}

func TestBuildType_GetParametersExcludeInherited(t *testing.T) {
	client := setup()
	pa := newPropertyAssertions(t)
	require := require.New(t)
	created := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, testBuildTypeId, true)

	//Add parameters to parent project
	proj, _ := client.Projects.GetByID(testBuildTypeProjectId)
	proj.Parameters.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "project_inherited", "value")
	proj, _ = client.Projects.Update(proj)

	pa.assertPropertyExists(proj.Parameters.Properties(), "project_inherited")

	sut := client.BuildTypes
	actual, err := sut.GetByID(created.ID) //Refresh
	props := teamcity.NewParametersEmpty()
	props.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "param1", "value1")
	props.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "param2", "value2")
	actual.Parameters = props
	actual, err = sut.Update(actual)

	cleanUpProject(t, client, testBuildTypeProjectId)

	require.NoError(err)

	pa.assertPropertyValue(actual.Parameters.Properties(), "param1", "value1")
	pa.assertPropertyDoesNotExist(actual.Parameters.Properties(), "project_inherited")
}

func TestBuildType_AttachVcsRoot(t *testing.T) {
	client := setup()
	assert := assert.New(t)

	createdBuildType := createTestBuildTypeWithName(t, client, testBuildTypeProjectId, "BuildRelease", true)

	newVcsRoot := getTestVcsRootData(testBuildTypeProjectId)

	vcsRootCreated, err := client.VcsRoots.Create(testBuildTypeProjectId, newVcsRoot)

	if err != nil {
		t.Fatalf("Failed to create vcs root: %s", err)
	}

	err = client.BuildTypes.AttachVcsRoot(createdBuildType.ID, vcsRootCreated)
	if err != nil {
		t.Fatalf("Failed to attach vcsRoot '%s' to buildType '%s': %s", createdBuildType.ID, vcsRootCreated.ID, err)
	}

	actual, err := client.BuildTypes.GetByID(createdBuildType.ID)
	if err != nil {
		t.Fatalf("Failed to get buildType '%s' for asserting: %s", createdBuildType.ID, err)
	}

	assert.NotEmpty(actual.VcsRootEntries, "Expected VcsRootEntries to contain at least one element")
	vcsEntries := idMapVcsRootEntries(actual.VcsRootEntries)
	assert.Containsf(vcsEntries, vcsRootCreated.ID, "Expected VcsRootEntries to contain the VcsRoot with id = %s, but it did not", vcsRootCreated.ID)

	cleanUpProject(t, client, testBuildTypeProjectId)
}

func idMapVcsRootEntries(v []*teamcity.VcsRootEntry) map[string]string {
	out := make(map[string]string)
	for _, item := range v {
		out[item.VcsRoot.ID] = item.ID
	}

	return out
}

func createTestBuildType(t *testing.T, client *teamcity.Client, buildTypeProjectId string) *teamcity.BuildType {
	return createTestBuildTypeWithName(t, client, buildTypeProjectId, "PullRequest", true)
}

func createTestBuildTypeWithName(t *testing.T, client *teamcity.Client, buildTypeProjectId string, name string, createProject bool) *teamcity.BuildType {
	return createTestBuildTypeInternal(t, client, buildTypeProjectId, name, createProject, false)
}

func createTestBuildTypeTemplateWithName(t *testing.T, client *teamcity.Client, buildTypeProjectId string, name string, createProject bool) *teamcity.BuildType {
	return createTestBuildTypeInternal(t, client, buildTypeProjectId, name, createProject, true)
}

func createTestBuildTypeInternal(t *testing.T, client *teamcity.Client, buildTypeProjectId string, name string, createProject bool, template bool) *teamcity.BuildType {
	if createProject {
		newProject := getTestProjectData(buildTypeProjectId, "")

		if _, err := client.Projects.Create(newProject); err != nil {
			t.Fatalf("Failed to create project for buildType: %s", err)
		}
	}

	newBuildType := getTestBuildTypeData(name, "Inspection", buildTypeProjectId, template)

	createdBuildType, err := client.BuildTypes.Create(buildTypeProjectId, newBuildType)
	if err != nil {
		t.Fatalf("Failed to CreateBuildType: %s", err)
	}

	detailed, _ := client.BuildTypes.GetByID(createdBuildType.ID)
	return detailed
}

func getTestBuildTypeData(name string, description string, projectId string, template bool) (out *teamcity.BuildType) {
	if template {
		out, _ = teamcity.NewBuildTypeTemplate(projectId, name)
	} else {
		out, _ = teamcity.NewBuildType(projectId, name)
	}
	out.Description = description
	return
}

func cleanUpBuildType(t *testing.T, c *teamcity.Client, id string) {
	if err := c.BuildTypes.Delete(id); err != nil {
		t.Errorf("Unable to delete build type with id = '%s', err: %s", id, err)
		return
	}

	deleted, err := c.BuildTypes.GetByID(id)

	if deleted != nil {
		t.Errorf("Build type not deleted during cleanup.")
		return
	}

	if err == nil {
		t.Errorf("Expected 404 Not Found error when getting Deleted Build Type, but no error returned.")
	}
}
