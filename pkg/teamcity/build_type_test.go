package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
)

func TestCreateBuildTypeForProject(t *testing.T) {
	client := setup()
	newProject := getTestProjectData(testBuildTypeProjectId)
	_, err := client.Projects.Create(newProject)

	if err != nil {
		t.Fatalf("Failed to create project for buildType: %s", err)
	}
	newBuildType := getTestBuildTypeData("PullRequest", "Description", testBuildTypeProjectId)

	actual, err := client.BuildTypes.Create(testBuildTypeProjectId, newBuildType)

	if err != nil {
		t.Fatalf("Failed to CreateBuildType: %s", err)
	}

	if actual == nil {
		t.Fatalf("CreateBuildType did not return a valid instance")
	}

	cleanUpProject(t, client, testBuildTypeProjectId)

	assert.NotEmpty(t, actual.ID)
	assert.Equal(t, newBuildType.ProjectID, actual.ProjectID)
	assert.Equal(t, newBuildType.Name, actual.Name)
}

func TestAttachVcsRoot(t *testing.T) {
	client := setup()
	newProject := getTestProjectData(testBuildTypeProjectId)

	if _, err := client.Projects.Create(newProject); err != nil {
		t.Fatalf("Failed to create project for buildType: %s", err)
	}

	newBuildType := getTestBuildTypeData("PullRequest", "Description", testBuildTypeProjectId)

	createdBuildType, err := client.BuildTypes.Create(testBuildTypeProjectId, newBuildType)
	if err != nil {
		t.Fatalf("Failed to CreateBuildType: %s", err)
	}

	newVcsRoot := getTestVcsRootData(testBuildTypeProjectId)

	vcsRootCreated, err := client.VcsRoots.Create(testBuildTypeProjectId, newVcsRoot)

	if err != nil {
		t.Fatalf("Failed to create vcs root: %s", err)
	}

	err = client.BuildTypes.AttachVcsRoot(createdBuildType.ID, vcsRootCreated)
	if err != nil {
		t.Fatalf("Failed to attach vcsRoot '%s' to buildType '%s': %s", createdBuildType.ID, vcsRootCreated.ID, err)
	}

	actual, err := client.BuildTypes.GetById(createdBuildType.ID)
	if err != nil {
		t.Fatalf("Failed to get buildType '%s' for asserting: %s", createdBuildType.ID, err)
	}

	assert.NotEqualf(t, actual.VcsRootEntries.Count, 0, "Expected VcsRootEntries to contain at least one element")
	vcsEntries := idMapVcsRootEntries(actual.VcsRootEntries)
	assert.Containsf(t, vcsEntries, vcsRootCreated.ID, "Expected VcsRootEntries to contain the VcsRoot with id = %s, but it did not", vcsRootCreated.ID)

	cleanUpProject(t, client, testBuildTypeProjectId)
}

func TestAddBuildStep(t *testing.T) {
	client := setup()
	updatedBuildType := createTestBuildStep(t, client, testBuildTypeProjectId)

	cleanUpProject(t, client, testBuildTypeProjectId)

	actual := updatedBuildType.Steps.Items

	assert.NotEmpty(t, actual)
}

func TestDeleteBuildStep(t *testing.T) {
	client := setup()
	updatedBuildType := createTestBuildStep(t, client, testBuildTypeProjectId)

	deleteStep := updatedBuildType.Steps.Items[0]

	client.BuildTypes.DeleteStep(updatedBuildType.ID, deleteStep.ID)

	updatedBuildType, _ = client.BuildTypes.GetById(updatedBuildType.ID)

	cleanUpProject(t, client, testBuildTypeProjectId)

	actual := updatedBuildType.Steps.Items

	assert.Empty(t, actual)
}

func TestAddAgentRequirement(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)
	req, _ := teamcity.NewAgentRequirement(teamcity.Conditions.Equals, "param", "value")

	client.BuildTypes.AddAgentRequirement(buildType.ID, req)
	buildType, _ = client.BuildTypes.GetById(buildType.ID) //refresh

	actual := buildType.AgentRequirements

	cleanUpProject(t, client, testBuildTypeProjectId)
	assert.NotEmpty(actual.Items)
	assert.Equal(teamcity.Conditions.Equals, actual.Items[0].Condition)

	assert.Equal("property-name", actual.Items[0].Properties.Items[0].Name)
	assert.Equal("param", actual.Items[0].Properties.Items[0].Value)
	assert.Equal("property-value", actual.Items[0].Properties.Items[1].Name)
	assert.Equal("value", actual.Items[0].Properties.Items[1].Value)
}

func TestSettings(t *testing.T) {
	client := setup()
	assert := assert.New(t)

	buildType := createTestBuildType(t, client, testBuildTypeProjectId)
	builder := teamcity.BuildTypeSettingsBuilder

	settings := builder.ConfigurationType("composite").
		PersonalBuildTrigger(false).
		ArtifactRules("**/*.zip").
		Build()

	err := client.BuildTypes.UpdateSettings(buildType.ID, settings)
	assert.Nil(err)

	buildType, _ = client.BuildTypes.GetById(buildType.ID) //refresh
	cleanUpProject(t, client, testBuildTypeProjectId)

	actual := buildType.Settings.Map()

	assert.Equal("COMPOSITE", actual["buildConfigurationType"])
	assert.Equal("false", actual["allowPersonalBuildTriggering"])
	assert.Equal("**/*.zip", actual["artifactRules"])
}

func idMapVcsRootEntries(v *teamcity.VcsRootEntries) map[string]string {
	out := make(map[string]string)
	for _, item := range v.Items {
		out[item.VcsRoot.ID] = item.Id
	}

	return out
}

func createTestBuildType(t *testing.T, client *teamcity.Client, buildTypeProjectId string) *teamcity.BuildType {
	newProject := getTestProjectData(buildTypeProjectId)

	if _, err := client.Projects.Create(newProject); err != nil {
		t.Fatalf("Failed to create project for buildType: %s", err)
	}

	return createTestBuildTypeWithName(t, client, buildTypeProjectId, "PullRequest")
}

func createTestBuildTypeWithName(t *testing.T, client *teamcity.Client, buildTypeProjectId string, name string) *teamcity.BuildType {
	newBuildType := getTestBuildTypeData(name, "Inspection", buildTypeProjectId)

	createdBuildType, err := client.BuildTypes.Create(buildTypeProjectId, newBuildType)
	if err != nil {
		t.Fatalf("Failed to CreateBuildType: %s", err)
	}

	detailed, _ := client.BuildTypes.GetById(createdBuildType.ID)
	return detailed
}

func createTestBuildStep(t *testing.T, client *teamcity.Client, buildTypeProjectId string) *teamcity.BuildType {
	createdBuildType := createTestBuildType(t, client, buildTypeProjectId)

	builder := teamcity.StepPowershellBuilder
	step := builder.ScriptFile("build.ps1").Build("step1")

	if err := client.BuildTypes.AddStep(createdBuildType.ID, step); err != nil {
		t.Fatalf("Failed to add step to buildType '%s'", createdBuildType.ID)
	}

	updated, _ := client.BuildTypes.GetById(createdBuildType.ID)
	return updated
}

func getTestBuildTypeData(name string, description string, projectId string) *teamcity.BuildType {

	return &teamcity.BuildType{
		Name:        name,
		Description: description,
		ProjectID:   projectId,
	}
}

func cleanUpBuildType(t *testing.T, c *teamcity.Client, id string) {
	if err := c.BuildTypes.Delete(id); err != nil {
		t.Errorf("Unable to delete build type with id = '%s', err: %s", id, err)
		return
	}

	deleted, err := c.BuildTypes.GetById(id)

	if deleted != nil {
		t.Errorf("Build type not deleted during cleanup.")
		return
	}

	if err == nil {
		t.Errorf("Expected 404 Not Found error when getting Deleted Build Type, but no error returned.")
	}
}
