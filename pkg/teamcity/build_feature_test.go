package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildFeature_CommitPublisher_Create(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)

	opt := teamcity.NewCommitStatusPublisherGithubOptionsToken("https://api.github.com", "1234")
	ni, _ := teamcity.NewFeatureCommitStatusPublisherGithub(opt)

	sut := client.BuildFeatureService(buildType.ID)
	actual, err := sut.Create(ni)

	cleanUpProject(t, client, testBuildTypeProjectId)

	require.NoError(t, err)
	assert.NotNil(actual)

	csp := actual.(*teamcity.FeatureCommitStatusPublisher)

	assert.NotEqual("", csp.ID())
	assert.Equal(buildType.ID, csp.BuildTypeID())
	assert.Equal("commit-status-publisher", csp.Type())
	assert.Equal(false, csp.Disabled())
}

func TestBuildFeature_CommitPublisher_Get(t *testing.T) {
	client := setup()
	assert := assert.New(t)
	buildType := createTestBuildType(t, client, testBuildTypeProjectId)

	opt := teamcity.NewCommitStatusPublisherGithubOptionsToken("https://api.github.com", "1234")
	ni, _ := teamcity.NewFeatureCommitStatusPublisherGithub(opt)

	sut := client.BuildFeatureService(buildType.ID)
	actual, err := sut.Create(ni)

	require.NoError(t, err)
	assert.NotNil(actual)

	actual, err = sut.GetByID(actual.ID())

	require.NoError(t, err)
	assert.NotNil(actual)

	csp := actual.(*teamcity.FeatureCommitStatusPublisher)

	cleanUpProject(t, client, testBuildTypeProjectId)
	assert.NotEqual("", csp.ID())
	assert.Equal(buildType.ID, csp.BuildTypeID())
	assert.Equal("commit-status-publisher", csp.Type())
	assert.Equal(false, csp.Disabled())
}
