package teamcity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TriggerVcsOptionsDefaults(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	props := newPropertyAssertions(t)

	actual, _ := NewTriggerVcsOptions(QuietPeriodDoNotUse, 0)

	require.NotNil(actual)

	assert.Equal(true, actual.QueueOptimization())
	assert.Equal(false, actual.PerCheckinTriggering())
	assert.Equal(QuietPeriodDoNotUse, actual.QuietPeriodMode)

	vcsProps := actual.properties()

	// False properties should be omitted
	props.assertPropertyDoesNotExist(vcsProps, "perCheckinTriggering")
	props.assertPropertyDoesNotExist(vcsProps, "groupCheckinsByCommitter")

	props.assertPropertyValue(vcsProps, "enableQueueOptimization", "true")
	props.assertPropertyValue(actual.properties(), "quietPeriodMode", "DO_NOT_USE")
}

func Test_TriggerVcsOptionsQuietPeriodModeCustom(t *testing.T) {
	assert := assert.New(t)
	props := newPropertyAssertions(t)

	_, err := NewTriggerVcsOptions(QuietPeriodCustom, -1)

	assert.Error(err, "expected error when QuietPeriodCustom and QuietPeriodSeconds is lower than zero")

	_, err = NewTriggerVcsOptions(QuietPeriodCustom, 0)
	assert.Error(err, "expected error when QuietPeriodCustom and QuietPeriodSeconds is zero")

	actual, err := NewTriggerVcsOptions(QuietPeriodCustom, 10)
	assert.NoError(err)
	assert.Equal(QuietPeriodCustom, actual.QuietPeriodMode)
	assert.Equal(10, actual.QuietPeriodInSeconds)

	props.assertPropertyValue(actual.properties(), "quietPeriodMode", "USE_CUSTOM")
	props.assertPropertyValue(actual.properties(), "quietPeriod", "10")
}

func Test_TriggerVcsOptionsQuietPeriodModeNotCustom(t *testing.T) {
	assert := assert.New(t)
	props := newPropertyAssertions(t)

	actual, _ := NewTriggerVcsOptions(QuietPeriodDoNotUse, -1)

	assert.Equal(QuietPeriodDoNotUse, actual.QuietPeriodMode)
	assert.Zero(actual.QuietPeriodInSeconds)
	props.assertPropertyValue(actual.properties(), "quietPeriodMode", "DO_NOT_USE")
	props.assertPropertyDoesNotExist(actual.properties(), "quietPeriod")

	actual, _ = NewTriggerVcsOptions(QuietPeriodUseDefault, 0)

	assert.Equal(QuietPeriodUseDefault, actual.QuietPeriodMode)
	assert.Zero(actual.QuietPeriodInSeconds)
	props.assertPropertyValue(actual.properties(), "quietPeriodMode", "USE_DEFAULT")
	props.assertPropertyDoesNotExist(actual.properties(), "quietPeriod")

	actual, _ = NewTriggerVcsOptions(QuietPeriodUseDefault, 10)
	assert.Equal(QuietPeriodUseDefault, actual.QuietPeriodMode)
	assert.Zero(actual.QuietPeriodInSeconds)
}

func Test_TriggerVcsOptions_EnableQueueOptimizationDisablesPerCheckinTriggering(t *testing.T) {
	assert := assert.New(t)
	props := newPropertyAssertions(t)

	actual, _ := NewTriggerVcsOptions(QuietPeriodDoNotUse, 0)

	actual.SetPerCheckinTriggering(true)
	actual.SetQueueOptimization(true)

	assert.False(actual.perCheckinTriggering)
	assert.True(actual.enableQueueOptimization)

	props.assertPropertyValue(actual.properties(), "enableQueueOptimization", "true")
}

func Test_TriggerVcsOptions_EnablePerCheckinTriggeringDisablesQueueOptimization(t *testing.T) {
	assert := assert.New(t)
	props := newPropertyAssertions(t)

	actual, _ := NewTriggerVcsOptions(QuietPeriodDoNotUse, 0)

	actual.SetQueueOptimization(true)
	actual.SetPerCheckinTriggering(true)

	assert.False(actual.QueueOptimization())
	assert.True(actual.PerCheckinTriggering())

	props.assertPropertyValue(actual.properties(), "perCheckinTriggering", "true")
}
