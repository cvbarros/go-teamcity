package teamcity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FinishTriggerOptionsConstructor(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	props := newPropertyAssertions(t)

	actual := NewFinishBuildTriggerOptions(false, nil)

	require.NotNil(actual)

	assert.Nil(actual.BranchFilter)
	assert.False(actual.AfterSuccessfulBuildOnly)

	triggerProps := actual.properties()

	// False properties should be omitted
	props.assertPropertyDoesNotExist(triggerProps, "afterSuccessfulBuildOnly")
	props.assertPropertyDoesNotExist(triggerProps, "branchFilter")
}

func Test_FinishTriggerOptionsConvertToProperties(t *testing.T) {
	require := require.New(t)
	props := newPropertyAssertions(t)

	actual := NewFinishBuildTriggerOptions(true, []string{"+:<default>", "-:/refs/(pull/*)/head"})

	require.NotNil(actual)

	triggerProps := actual.properties()

	// False properties should be omitted
	props.assertPropertyValue(triggerProps, "afterSuccessfulBuildOnly", "true")
	props.assertPropertyValue(triggerProps, "branchFilter", "+:<default>\n-:/refs/(pull/*)/head")
}
