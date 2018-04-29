package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSnapshotDependency(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	bt := &teamcity.BuildTypeReference{
		ID:   "someBuildID",
		Name: "someBuildName",
	}

	actual := teamcity.NewSnapshotDependency(bt)

	require.NotNil(actual)
	assert.Equal(bt, actual.SourceBuildType)
	assert.Equal("snapshot_dependency", actual.Type)
	require.NotEmpty(actual.Properties)
	props := actual.Properties.Map()

	assert.Contains(props, "run-build-if-dependency-failed")
	assert.Contains(props, "run-build-if-dependency-failed-to-start")
	assert.Contains(props, "run-build-on-the-same-agent")
	assert.Contains(props, "take-started-build-with-same-revisions")
	assert.Contains(props, "take-successful-builds-only")
}

func TestNewSnapshotDependencyWithOptions(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	bt := &teamcity.BuildTypeReference{
		ID:   "someBuildID",
		Name: "someBuildName",
	}

	opt := &teamcity.SnapshotDependencyOptions{
		OnFailedDependency:       "RUN",
		RunSameAgent:             true,
		TakeSuccessfulBuildsOnly: false,
	}

	actual := teamcity.NewSnapshotDependencyWithOptions(bt, opt)

	require.NotNil(actual)
	assert.Equal(actual.SourceBuildType, bt)
	assert.Equal("snapshot_dependency", actual.Type)

	require.NotEmpty(actual.Properties)
	props := actual.Properties.Map()

	assert.Equal("RUN", props["run-build-if-dependency-failed"])
	assert.Equal("true", props["run-build-on-the-same-agent"])
	assert.Equal("false", props["take-successful-builds-only"])
}
