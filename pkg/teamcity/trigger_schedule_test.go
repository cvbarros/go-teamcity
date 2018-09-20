package teamcity

import (
	"encoding/json"
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

func Test_TriggerScheduleDefaultOptions(t *testing.T) {
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
	assert.Equal("-:*.md", sut.Rules[0])
}

const triggerDailyJSON = `
{
	"id": "TRIGGER_1",
	"type": "schedulingTrigger",
	"properties": {
		"count": 18,
		"property": [
			{
				"name": "cronExpression_dm",
				"value": "*"
			},
			{
				"name": "cronExpression_dw",
				"value": "?"
			},
			{
				"name": "cronExpression_hour",
				"value": "*"
			},
			{
				"name": "cronExpression_min",
				"value": "0"
			},
			{
				"name": "cronExpression_month",
				"value": "*"
			},
			{
				"name": "cronExpression_sec",
				"value": "0"
			},
			{
				"name": "cronExpression_year",
				"value": "*"
			},
			{
				"name": "dayOfWeek",
				"value": "Sunday"
			},
			{
				"name": "enableQueueOptimization",
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
				"name": "schedulingPolicy",
				"value": "daily"
			},
			{
				"name": "timezone",
				"value": "SERVER"
			},
			{
				"name": "triggerBuildWithPendingChangesOnly",
				"value": "true"
			},
			{
				"name": "triggerRules",
				"value": "-:*.md\n+:*"
			}
		]
	}
}
`
