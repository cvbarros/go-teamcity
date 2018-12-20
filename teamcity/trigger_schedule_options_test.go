package teamcity

import "testing"

func Test_ForceFalsePropertiesIfApplicable(t *testing.T) {
	sut := NewTriggerScheduleOptions()
	pa := newPropertyAssertions(t)
	// These 3 always get computed with "true" if ommitted from TeamCity
	// Thus, we need to make sure to force and output 'false' when converting to properties
	sut.QueueOptimization = false
	sut.PromoteWatchedBuild = false
	sut.BuildWithPendingChangesOnly = false

	actual := sut.properties()

	pa.assertPropertyValue(actual, "enableQueueOptimization", "false")
	pa.assertPropertyValue(actual, "promoteWatchedBuild", "false")
	pa.assertPropertyValue(actual, "triggerBuildWithPendingChangesOnly", "false")
}
