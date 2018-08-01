package teamcity_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
)

func Test_ArtifactDependency_Invariants(t *testing.T) {
	t.Run("sourceBuildID is required", func(t *testing.T) {
		_, err := teamcity.NewArtifactDependency("", &teamcity.ArtifactDependencyOptions{})
		assert.EqualError(t, err, "sourceBuildTypeID is required")
	})

	t.Run("opt must be non-nil and valid", func(t *testing.T) {
		_, err := teamcity.NewArtifactDependency("sourceBuild", nil)
		assert.EqualError(t, err, "options must be valid")
	})
}

func Test_ArtifactDependency_Constructor(t *testing.T) {
	assert := assert.New(t)

	actual, _ := teamcity.NewArtifactDependency("sourceBuild", createDefaultTestingArtifactDependencyOptions())
	require.NotNil(t, actual)

	assert.Equal("sourceBuild", actual.SourceBuildTypeID())
	assert.Equal("artifact_dependency", actual.Type())
	assert.EqualValues(false, actual.Disabled())
}

func testPathRules() []string {
	return []string{"rule1", "rule2"}
}

func createDefaultTestingArtifactDependencyOptions() *teamcity.ArtifactDependencyOptions {
	c, _ := teamcity.NewArtifactDependencyOptions(testPathRules(), teamcity.LatestSuccessfulBuild, false, "")
	return c
}
