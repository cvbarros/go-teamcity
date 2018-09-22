package teamcity

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func Test_TriggerScheduleDeserializeDaily(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var dt triggerJSON
	var sut TriggerSchedule
	if err := json.Unmarshal([]byte(triggerDailyJSON), &dt); err != nil {
		t.Error(err)
	}

	err := sut.read(&dt)
	require.NoError(err)

	assert.Equal(TriggerSchedulingDaily, sut.SchedulingPolicy)
	assert.Equal(uint(12), sut.Hour)
	assert.Equal(uint(0), sut.Minute)
	assert.Equal("SERVER", sut.Timezone)
	assert.Equal("-:*.md", sut.Rules[0])
	assert.Equal("+:*", sut.Rules[1])
}

func Test_TriggerScheduleSerializeDaily(t *testing.T) {
	require := require.New(t)
	pa := newPropertyAssertions(t)

	var dt, _ = NewTriggerScheduleDaily("someBuild", 12, 0, "SERVER", []string{"+:*", "-:*.md"})
	jsonBytes, err := dt.MarshalJSON()

	require.NoError(err)

	var actual triggerJSON
	if err := json.Unmarshal([]byte(jsonBytes), &actual); err != nil {
		t.Error(err)
	}

	props := actual.Properties

	pa.assertPropertyValue(props, "schedulingPolicy", dt.SchedulingPolicy)
	pa.assertPropertyValue(props, "timezone", dt.Timezone)
	pa.assertPropertyValue(props, "hour", fmt.Sprint(dt.Hour))
	pa.assertPropertyValue(props, "minute", fmt.Sprint(dt.Minute))
	pa.assertPropertyValue(props, "triggerRules", "+:*\n-:*.md")

	//Default Options
	pa.assertPropertyValue(props, "enableQueueOptimization", "true")
	pa.assertPropertyValue(props, "promoteWatchedBuild", "true")
	pa.assertPropertyValue(props, "triggerBuildWithPendingChangesOnly", "true")
	pa.assertPropertyValue(props, "revisionRuleBuildBranch", "<default>")
	//Unset options must not be exported to properties
	pa.assertPropertyDoesNotExist(props, "triggerBuildIfWatchedBuildChanges")
	pa.assertPropertyDoesNotExist(props, "triggerBuildOnAllCompatibleAgents")
	pa.assertPropertyDoesNotExist(props, "enforceCleanCheckout")
	pa.assertPropertyDoesNotExist(props, "enforceCleanCheckoutForDependencies")
}

func Test_TriggerDeserializeScheduleOptions(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var dt triggerJSON
	var sut TriggerSchedule
	if err := json.Unmarshal([]byte(triggerScheduleOptions), &dt); err != nil {
		t.Error(err)
	}

	err := sut.read(&dt)
	require.NoError(err)

	assert.Equal(TriggerSchedulingDaily, sut.SchedulingPolicy)
	require.NotNil(sut.Options)

	opt := sut.Options
	assert.Equal(true, opt.QueueOptimization)
	assert.Equal(true, opt.EnforceCleanCheckout)
	assert.Equal(true, opt.EnforceCleanCheckoutForDependencies)
	assert.Equal(true, opt.PromoteWatchedBuild)
	assert.Equal(LatestFinishedBuild, opt.RevisionRule)
	assert.Equal("Project1_ReleasetoTesting", opt.RevisionRuleSourceBuildID)
	assert.Equal(true, opt.TriggerIfWatchedBuildChanges)
	assert.Equal(true, opt.BuildOnAllCompatibleAgents)
	assert.Equal(true, opt.BuildWithPendingChangesOnly)
}

const triggerDailyJSON = `
{
	"id": "TRIGGER_1",
	"type": "schedulingTrigger",
	"properties": {
		"count": 5,
		"property": [
			{
				"name": "hour",
				"value": "12"
			},
			{
				"name": "minute",
				"value": "0"
			},
			{
				"name": "schedulingPolicy",
				"value": "daily"
			},
			{
				"name": "timezone",
				"value": "SERVER"
			},
			{
				"name": "triggerRules",
				"value": "-:*.md\n+:*"
			}
		]
	}
}
`
const triggerScheduleOptions = `
{
	"id": "TRIGGER_12",
	"type": "schedulingTrigger",
	"properties": {
		"property": [
			{
				"name": "enableQueueOptimization",
				"value": "true"
			},
			{
				"name": "enforceCleanCheckout",
				"value": "true"
			},
			{
				"name": "enforceCleanCheckoutForDependencies",
				"value": "true"
			},
			{
				"name": "hour",
				"value": "12"
			},
			{
				"name": "minute",
				"value": "0"
			},
			{
				"name": "promoteWatchedBuild",
				"value": "true"
			},
			{
				"name": "revisionRule",
				"value": "lastFinished"
			},
			{
				"name": "revisionRuleBuildBranch",
				"value": "<default>"
			},
			{
				"name": "revisionRuleDependsOn",
				"value": "Project1_ReleasetoTesting"
			},
			{
				"name": "schedulingPolicy",
				"value": "daily"
			},
			{
				"name": "timezone",
				"value": "SERVER"
			},
			{
				"name": "triggerBuildIfWatchedBuildChanges",
				"value": "true"
			},
			{
				"name": "triggerBuildOnAllCompatibleAgents",
				"value": "true"
			},
			{
				"name": "triggerBuildWithPendingChangesOnly",
				"value": "true"
			},
			{
				"name": "triggerRules",
				"value": "+:*\n-:*.md"
			}
		]
	}
}
`
