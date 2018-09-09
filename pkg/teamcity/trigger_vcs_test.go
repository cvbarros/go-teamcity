package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrigger_Constructor(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	actual, _ := teamcity.NewTriggerVcs("+:*", "")

	require.NotNil(actual)
	assert.Equal("vcsTrigger", actual.Type())

	assert.Equal("+:*", actual.Rules)
	assert.Empty(actual.BranchFilter)
}
