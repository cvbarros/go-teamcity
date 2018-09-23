package teamcity

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func Test_ArtifactDependencyOptions_Invariants(t *testing.T) {
	t.Run("pathRules is required", func(t *testing.T) {
		_, err := NewArtifactDependencyOptions([]string{}, LatestSuccessfulBuild, false, "")
		assert.EqualError(t, err, "pathRules is required", "pathRules is required")
	})

	t.Run("revisionType is required", func(t *testing.T) {
		_, err := NewArtifactDependencyOptions([]string{"rule1"}, "", false, "")
		assert.EqualError(t, err, "revisionType is required", "revisionType is required")
	})

	t.Run("revisionValue is required is using 'BuildWithSpecifiedNumber'", func(t *testing.T) {
		_, err := NewArtifactDependencyOptions([]string{"rule1"}, BuildWithSpecifiedNumber, false, "")
		assert.Error(t, err, "revisionValue is required is using 'BuildWithSpecifiedNumber'")
	})

	t.Run("revisionValue is required is using 'LastBuildFinishedWithTag'", func(t *testing.T) {
		_, err := NewArtifactDependencyOptions([]string{"rule1"}, LastBuildFinishedWithTag, false, "")
		assert.Error(t, err, "revisionValue is required is using 'LastBuildFinishedWithTag'")
	})
}

func Test_ArtifactDependencyOptions_ConstructorDefault(t *testing.T) {
	require := require.New(t)
	props := newPropertyAssertions(t)

	actual, _ := NewArtifactDependencyOptions([]string{"rule1", "rule2"}, LatestSuccessfulBuild, false, "")

	require.NotNil(actual)

	artifactProps := actual.properties()

	props.assertPropertyValue(artifactProps, "cleanDestinationDirectory", "false")
	props.assertPropertyValue(artifactProps, "pathRules", "rule1\r\nrule2")
	props.assertPropertyValue(artifactProps, "revisionValue", "latest.lastSuccessful")
	props.assertPropertyValue(artifactProps, "revisionName", "lastSuccessful")
}

func Test_ArtifactDependencyOptions_ConstructorLatestPinned(t *testing.T) {
	require := require.New(t)
	props := newPropertyAssertions(t)

	actual, _ := NewArtifactDependencyOptions([]string{"rule1", "rule2"}, LatestPinnedBuild, false, "")

	require.NotNil(actual)

	artifactProps := actual.properties()

	props.assertPropertyValue(artifactProps, "revisionValue", "latest.lastPinned")
	props.assertPropertyValue(artifactProps, "revisionName", "lastPinned")
}

func Test_ArtifactDependencyOptions_ConstructorLatestFinished(t *testing.T) {
	require := require.New(t)
	props := newPropertyAssertions(t)

	actual, _ := NewArtifactDependencyOptions([]string{"rule1", "rule2"}, LatestFinishedBuild, false, "")

	require.NotNil(actual)

	artifactProps := actual.properties()

	props.assertPropertyValue(artifactProps, "revisionValue", "latest.lastFinished")
	props.assertPropertyValue(artifactProps, "revisionName", "lastFinished")
}

func Test_ArtifactDependencyOptions_ConstructorSameChain(t *testing.T) {
	require := require.New(t)
	props := newPropertyAssertions(t)

	actual, _ := NewArtifactDependencyOptions([]string{"rule1", "rule2"}, BuildFromSameChain, false, "")

	require.NotNil(actual)

	artifactProps := actual.properties()

	props.assertPropertyValue(artifactProps, "revisionValue", "latest.sameChainOrLastFinished")
	props.assertPropertyValue(artifactProps, "revisionName", "sameChainOrLastFinished")
}

func Test_ArtifactDependencyOptions_ConstructorSpecificBuildNumber(t *testing.T) {
	require := require.New(t)
	props := newPropertyAssertions(t)

	actual, _ := NewArtifactDependencyOptions([]string{"rule1", "rule2"}, BuildWithSpecifiedNumber, false, "123")

	require.NotNil(actual)

	artifactProps := actual.properties()

	props.assertPropertyValue(artifactProps, "revisionValue", "123")
	props.assertPropertyValue(artifactProps, "revisionName", "buildNumber")
}

func Test_ArtifactDependencyOptions_ConstructorLastFinishedWithTag(t *testing.T) {
	require := require.New(t)
	props := newPropertyAssertions(t)

	actual, _ := NewArtifactDependencyOptions([]string{"rule1", "rule2"}, LastBuildFinishedWithTag, false, "tag1")

	require.NotNil(actual)

	artifactProps := actual.properties()

	props.assertPropertyValue(artifactProps, "revisionValue", "tag1.tcbuildtag") // TC UI appends "tag1" to this suffix when calling API
	props.assertPropertyValue(artifactProps, "revisionName", "buildTag")
}

func Test_ArtifactDependencyOptions_BuildTagStrip(t *testing.T) {
	assert := assert.New(t)
	pa := newPropertyAssertions(t)
	sut, _ := NewArtifactDependencyOptions([]string{"rule1"}, LastBuildFinishedWithTag, false, "tag1")

	props := sut.properties()

	pa.assertPropertyValue(props, "revisionValue", "tag1.tcbuildtag")

	actual := props.artifactDependencyOptions()

	assert.Equal("tag1", actual.RevisionNumber)
}
