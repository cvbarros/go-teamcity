package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildType_Create(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	actual := createTestBuildTypeWithName(t, client, "BuildTypeProject", "BuildRelease", true)

	cleanUpProject(t, client, "BuildTypeProject")

	assert.NotEmpty(t, actual.ID)
	assert.Equal("BuildRelease", actual.Name)
	assert.Equal("BuildTypeProject", actual.ProjectID)

	//Verify some default properties
	optExpected := teamcity.NewBuildTypeOptionsWithDefaults()
	optActual := actual.Options
	require.NotNil(t, optActual)
	assert.Equal(optExpected.ArtifactRules, optActual.ArtifactRules)
	assert.Equal(optExpected.BuildCounter, optActual.BuildCounter)
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

	updated, err := sut.Update(actual)
	cleanUpProject(t, client, testBuildTypeProjectId)

	require.NoError(t, err)

	assert.Equal("Updated description", updated.Description)
	assert.Equal([]string{"rule1", "rule2"}, updated.Options.ArtifactRules)
	assert.Equal(10, updated.Options.BuildCounter)
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

func TestBuildType_AddStepPowerShell(t *testing.T) {
	client := setup()
	step, _ := teamcity.NewStepPowershellScriptFile("step1", "build.ps1", "")
	_, created := createTestBuildStep(t, client, step, testBuildTypeProjectId)

	cleanUpProject(t, client, testBuildTypeProjectId)

	assert.NotNil(t, created)
}

func TestBuildType_AddStepCommandLineExecutable(t *testing.T) {
	assert := assert.New(t)
	client := setup()
	step, _ := teamcity.NewStepCommandLineExecutable("step_exe", "./script.sh", "hello")
	_, actual := createTestBuildStep(t, client, step, testBuildTypeProjectId)

	cleanUpProject(t, client, testBuildTypeProjectId)

	require.NotNil(t, actual)
	assert.Equal(teamcity.StepTypeCommandLine, actual.Type())
}

func TestBuildType_AddStepCommandLineScript(t *testing.T) {
	assert := assert.New(t)
	client := setup()
	script := `echo "Hello World
	echo "World, Hello!
	export HELLO_WORLD=1
	`
	step, _ := teamcity.NewStepCommandLineScript("step_exe", script)
	_, actual := createTestBuildStep(t, client, step, testBuildTypeProjectId)

	cleanUpProject(t, client, testBuildTypeProjectId)

	require.NotNil(t, actual)
	assert.Equal(teamcity.StepTypeCommandLine, actual.Type())
}

func TestBuildType_GetSteps(t *testing.T) {
	client := setup()
	step1, _ := teamcity.NewStepCommandLineExecutable("step1", "./script.sh", "hello")
	step2, _ := teamcity.NewStepCommandLineExecutable("step2", "./script.sh", "hello")

	buildType := createTestBuildType(t, client, testBuildTypeProjectId)

	created1, _ := client.BuildTypes.AddStep(buildType.ID, step1)
	created2, _ := client.BuildTypes.AddStep(buildType.ID, step2)

	steps, _ := client.BuildTypes.GetSteps(buildType.ID)

	cleanUpProject(t, client, testBuildTypeProjectId)

	assert.Contains(t, steps, created1)
	assert.Contains(t, steps, created2)
}

func TestBuildType_DeleteStep(t *testing.T) {
	client := setup()
	step, _ := teamcity.NewStepCommandLineExecutable("step_exe", "./script.sh", "hello")
	buildType, s := createTestBuildStep(t, client, step, testBuildTypeProjectId)
	created := s.(*teamcity.StepCommandLine)
	client.BuildTypes.DeleteStep(buildType.ID, created.ID)

	steps, _ := client.BuildTypes.GetSteps(buildType.ID)

	cleanUpProject(t, client, testBuildTypeProjectId)

	assert.Empty(t, steps)
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
	if createProject {
		newProject := getTestProjectData(buildTypeProjectId)

		if _, err := client.Projects.Create(newProject); err != nil {
			t.Fatalf("Failed to create project for buildType: %s", err)
		}
	}

	newBuildType := getTestBuildTypeData(name, "Inspection", buildTypeProjectId)

	createdBuildType, err := client.BuildTypes.Create(buildTypeProjectId, newBuildType)
	if err != nil {
		t.Fatalf("Failed to CreateBuildType: %s", err)
	}

	detailed, _ := client.BuildTypes.GetByID(createdBuildType.ID)
	return detailed
}

func createTestBuildStep(t *testing.T, client *teamcity.Client, step teamcity.Step, buildTypeProjectId string) (*teamcity.BuildType, teamcity.Step) {
	createdBuildType := createTestBuildType(t, client, buildTypeProjectId)

	created, err := client.BuildTypes.AddStep(createdBuildType.ID, step)
	if err != nil {
		t.Fatalf("Failed to add step to buildType '%s'", createdBuildType.ID)
	}

	updated, _ := client.BuildTypes.GetByID(createdBuildType.ID)
	return updated, created
}

func getTestBuildTypeData(name string, description string, projectId string) *teamcity.BuildType {
	out, _ := teamcity.NewBuildType(projectId, name)
	out.Description = description
	return out
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
