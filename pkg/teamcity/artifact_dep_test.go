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
	propAssert := newPropertyAssertions(t)

	actual, _ := teamcity.NewArtifactDependency("sourceBuild", createDefaultTestingArtifactDependencyOptions())
	require.NotNil(t, actual)

	assert.Equal("sourceBuild", actual.SourceBuildType.ID)
	assert.Equal("artifact_dependency", actual.Type)
	assert.EqualValues(teamcity.NewFalse(), actual.Disabled)

	props := actual.Properties

	// No need to assert property values. This is done in artifact_dep_options_test.go
	propAssert.assertPropertyExists(props, "cleanDestinationDirectory")
	propAssert.assertPropertyExists(props, "pathRules")
	propAssert.assertPropertyExists(props, "revisionName")
	propAssert.assertPropertyExists(props, "revisionValue")
}

func testPathRules() []string {
	return []string{"rule1", "rule2"}
}

func createDefaultTestingArtifactDependencyOptions() *teamcity.ArtifactDependencyOptions {
	c, _ := teamcity.NewArtifactDependencyOptions(testPathRules(), teamcity.LatestSuccessfulBuild, false, "")
	return c
}
