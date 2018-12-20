package teamcity_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

const testBuildTemplateProjectId = "BuildTemplateProject"

func Test_AttachBuildTemplate(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTemplateProjectId)
	buildTemplate := createTestBuildTypeTemplateWithName(t, client, testBuildTemplateProjectId, "Template", false)

	sut := client.BuildTemplateService(buildType.ID)
	_, err := sut.Attach(buildTemplate.ID)

	require.NoError(t, err)

	actual, _ := client.BuildTypes.GetByID(buildType.ID)

	cleanUpProject(t, client, testBuildTemplateProjectId)
	assert.NotNil(actual.Templates)
	assert.Equal(int32(1), actual.Templates.Count)
	assert.Equal(buildTemplate.ID, actual.Templates.Items[0].ID)
}

func Test_DetachBuildTemplate(t *testing.T) {
	client := setup()
	require := require.New(t)
	buildType := createTestBuildType(t, client, testBuildTemplateProjectId)
	template1 := createTestBuildTypeTemplateWithName(t, client, testBuildTemplateProjectId, "Template1", false)
	template2 := createTestBuildTypeTemplateWithName(t, client, testBuildTemplateProjectId, "Template2", false)

	sut := client.BuildTemplateService(buildType.ID)
	_, err := sut.Attach(template1.ID)
	require.NoError(err)
	_, err = sut.Attach(template2.ID)
	require.NoError(err)

	err = sut.Detach(template1.ID)
	require.NoError(err)

	actual, _ := client.BuildTypes.GetByID(buildType.ID)

	cleanUpProject(t, client, testBuildTemplateProjectId)
	require.NotNil(actual.Templates)
	require.Equal(int32(1), actual.Templates.Count)
	require.Equal(template2.ID, actual.Templates.Items[0].ID)
}
