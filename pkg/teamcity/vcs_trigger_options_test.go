package teamcity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VcsTriggerOptionsDefaults(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	props := newPropertyAssertions(t)

	actual := NewVcsTriggerOptions()

	require.NotNil(actual)

	assert.Equal(true, actual.enableQueueOptimization)
	assert.Equal(false, actual.perCheckinTriggering)
	assert.Equal(QuietPeriodDoNotUse, actual.quietPeriodMode)

	vcsProps := actual.vcsTriggerProperties()

	// False properties should be omitted
	props.assertPropertyDoesNotExist(vcsProps, "perCheckinTriggering")
	props.assertPropertyDoesNotExist(vcsProps, "groupCheckinsByCommitter")

	props.assertPropertyValue(vcsProps, "enableQueueOptimization", "true")
	props.assertPropertyValue(actual.vcsTriggerProperties(), "quietPeriodMode", "DO_NOT_USE")
}

func Test_VcsTriggerOptionsGroupUserCheckin(t *testing.T) {
	props := newPropertyAssertions(t)
	actual := NewVcsTriggerOptions()
	actual.SetGroupUserCheckins(true)
	assert.True(t, actual.groupUserCheckins)

	props.assertPropertyValue(actual.vcsTriggerProperties(), "groupCheckinsByCommitter", "true")
}

func Test_VcsTriggerOptionsQuietPeriodModeCustom(t *testing.T) {
	assert := assert.New(t)
	props := newPropertyAssertions(t)

	actual := NewVcsTriggerOptions()

	err := actual.SetQuietPeriodMode(QuietPeriodCustom, -1)
	assert.Error(err, "expected error when QuietPeriodCustom and QuietPeriodSeconds is lower than zero")

	err = actual.SetQuietPeriodMode(QuietPeriodCustom, 0)
	assert.Error(err, "expected error when QuietPeriodCustom and QuietPeriodSeconds is zero")

	err = actual.SetQuietPeriodMode(QuietPeriodCustom, 10)
	assert.NoError(err)
	assert.Equal(QuietPeriodCustom, actual.quietPeriodMode)
	assert.Equal(10, actual.quietPeriodInSeconds)

	props.assertPropertyValue(actual.vcsTriggerProperties(), "quietPeriodMode", "USE_CUSTOM")
	props.assertPropertyValue(actual.vcsTriggerProperties(), "quietPeriod", "10")
}

func Test_VcsTriggerOptionsQuietPeriodModeNotCustom(t *testing.T) {
	assert := assert.New(t)
	props := newPropertyAssertions(t)

	actual := NewVcsTriggerOptions()

	actual.SetQuietPeriodMode(QuietPeriodDoNotUse, -1)
	assert.Equal(QuietPeriodDoNotUse, actual.quietPeriodMode)
	assert.Zero(actual.quietPeriodInSeconds)
	props.assertPropertyValue(actual.vcsTriggerProperties(), "quietPeriodMode", "DO_NOT_USE")
	props.assertPropertyDoesNotExist(actual.vcsTriggerProperties(), "quietPeriod")

	actual.SetQuietPeriodMode(QuietPeriodUseDefault, 0)
	assert.Equal(QuietPeriodUseDefault, actual.quietPeriodMode)
	assert.Zero(actual.quietPeriodInSeconds)
	props.assertPropertyValue(actual.vcsTriggerProperties(), "quietPeriodMode", "USE_DEFAULT")
	props.assertPropertyDoesNotExist(actual.vcsTriggerProperties(), "quietPeriod")

	actual.SetQuietPeriodMode(QuietPeriodUseDefault, 10)
	assert.Equal(QuietPeriodUseDefault, actual.quietPeriodMode)
	assert.Zero(actual.quietPeriodInSeconds)
}

func Test_VcsTriggerOptions_EnableQueueOptimizationDisablesPerCheckinTriggering(t *testing.T) {
	assert := assert.New(t)
	props := newPropertyAssertions(t)

	actual := NewVcsTriggerOptions()

	actual.SetPerCheckinTriggering(true)
	actual.SetQueueOptimization(true)

	assert.False(actual.perCheckinTriggering)
	assert.True(actual.enableQueueOptimization)

	props.assertPropertyValue(actual.vcsTriggerProperties(), "enableQueueOptimization", "true")
}

func Test_VcsTriggerOptions_EnablePerCheckinTriggeringDisablesQueueOptimization(t *testing.T) {
	assert := assert.New(t)
	props := newPropertyAssertions(t)

	actual := NewVcsTriggerOptions()

	actual.SetQueueOptimization(true)
	actual.SetPerCheckinTriggering(true)

	assert.False(actual.enableQueueOptimization)
	assert.True(actual.perCheckinTriggering)

	props.assertPropertyValue(actual.vcsTriggerProperties(), "perCheckinTriggering", "true")
}

type PropertyAssertions struct {
	a *assert.Assertions
	t *testing.T
}

func newPropertyAssertions(t *testing.T) *PropertyAssertions {
	return &PropertyAssertions{a: assert.New(t), t: t}
}

func (p *PropertyAssertions) assertPropertyValue(props *Properties, name string, value string) {
	require.NotNil(p.t, props)

	propMap := props.Map()

	p.a.Contains(propMap, name)
	p.a.Equal(value, propMap[name])
}

func (p *PropertyAssertions) assertPropertyDoesNotExist(props *Properties, name string) {
	require.NotNil(p.t, props)

	propMap := props.Map()

	p.a.NotContains(propMap, name)
}
