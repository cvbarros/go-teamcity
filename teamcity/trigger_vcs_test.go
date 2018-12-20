package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrigger_Constructor(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	actual, _ := teamcity.NewTriggerVcs([]string{"+:*"}, []string{})

	require.NotNil(actual)
	assert.Equal("vcsTrigger", actual.Type())

	assert.Equal([]string{"+:*"}, actual.Rules)
	assert.Empty(actual.BranchFilter)
}
