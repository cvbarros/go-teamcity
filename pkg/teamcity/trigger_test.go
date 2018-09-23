package teamcity_test

import (
	"testing"
	"time"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrigger_CreateTriggerVcs(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerVcs("+:*", "")

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	assert.Equal(created.BuildTypeID(), bt.ID)
	cleanUpProject(t, client, bt.ProjectID)
}

func TestTrigger_CreateTriggerBuildFinish(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)
	st := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "SourceBuild", false)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerBuildFinish(st.ID, teamcity.NewTriggerBuildFinishOptions(true, []string{"+:<default>"}))

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	assert.Equal(created.BuildTypeID(), bt.ID)
	cleanUpProject(t, client, bt.ProjectID)
}

func TestTrigger_CreateTriggerScheduleDaily(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerScheduleDaily(bt.ID, 12, 30, "SERVER", []string{"+:*"})
	opt := teamcity.NewTriggerScheduleOptions()
	opt.QueueOptimization = false
	opt.BuildOnAllCompatibleAgents = true
	opt.PromoteWatchedBuild = false
	opt.BuildWithPendingChangesOnly = false
	nt.Options = opt
	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	assert.IsType(&teamcity.TriggerSchedule{}, created)
	assert.Equal(created.BuildTypeID(), bt.ID)
	cleanUpProject(t, client, bt.ProjectID)

	actual := created.(*teamcity.TriggerSchedule)
	assert.Equal(teamcity.TriggerSchedulingDaily, actual.SchedulingPolicy)
	assert.Equal(uint(12), actual.Hour)
	assert.Equal(uint(30), actual.Minute)
	assert.Equal(false, actual.Options.BuildWithPendingChangesOnly)
	assert.Equal(false, actual.Options.PromoteWatchedBuild)
	assert.Equal(false, actual.Options.QueueOptimization)
	assert.Equal(true, actual.Options.BuildOnAllCompatibleAgents)
}

func TestTrigger_CreateTriggerScheduleWeekly(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerScheduleWeekly(bt.ID, time.Thursday, 12, 0, "SERVER", []string{"+:*"})

	created, err := sut.AddTrigger(nt)

	require.Nil(err)
	assert.IsType(&teamcity.TriggerSchedule{}, created)
	assert.Equal(created.BuildTypeID(), bt.ID)
	cleanUpProject(t, client, bt.ProjectID)

	actual := created.(*teamcity.TriggerSchedule)
	assert.Equal(actual.SchedulingPolicy, teamcity.TriggerSchedulingWeekly)
	assert.Equal(actual.Weekday, time.Thursday)
}

func TestTrigger_GetTriggerVcs(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerVcs("+:*", "")

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	actual, err := sut.GetByID(created.ID())

	require.NoError(err)
	assert.Equal(created.ID(), actual.ID())
	assert.Equal(created.BuildTypeID(), actual.BuildTypeID())
	assert.Equal(created.Type(), actual.Type())

	cleanUpProject(t, client, bt.ProjectID)
}

func TestTrigger_GetTriggerBuildFinish(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)
	st := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "SourceBuild", false)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerBuildFinish(st.ID, teamcity.NewTriggerBuildFinishOptions(true, []string{"master", "feature"}))

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	actual, err := sut.GetByID(created.ID())

	cleanUpProject(t, client, bt.ProjectID)
	require.NoError(err)
	require.IsType(&teamcity.TriggerBuildFinish{}, actual)
	actualT := actual.(*teamcity.TriggerBuildFinish)

	assert.Equal(st.ID, actualT.SourceBuildID)
	assert.Equal(created.ID(), actual.ID())
	assert.Equal(created.BuildTypeID(), actual.BuildTypeID())
	assert.Equal(created.Type(), actual.Type())
	assert.Equal([]string{"master", "feature"}, actualT.Options.BranchFilter)
}

func TestTrigger_GetTriggerScheduleDaily(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerScheduleDaily(bt.ID, 12, 0, "SERVER", []string{"+:*"})

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	actual, err := sut.GetByID(created.ID())

	require.NoError(err)
	assert.Equal(created.ID(), actual.ID())
	assert.Equal(created.BuildTypeID(), actual.BuildTypeID())
	assert.Equal(created.Type(), actual.Type())

	cleanUpProject(t, client, bt.ProjectID)
}

func TestTrigger_DeleteTriggerVcs(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerVcs("+:*", "")

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	sut.Delete(created.ID())
	_, err = sut.GetByID(created.ID()) // refresh

	require.Error(err)
	assert.Contains(err.Error(), "404")
	cleanUpProject(t, client, bt.ProjectID)
}

func TestTrigger_DeleteTriggerBuildFinish(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)
	st := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "SourceBuild", false)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerBuildFinish(st.ID, teamcity.NewTriggerBuildFinishOptions(true, []string{"+:<default>"}))

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	sut.Delete(created.ID())
	_, err = sut.GetByID(created.ID()) // refresh

	require.Error(err)
	assert.Contains(err.Error(), "404")
	cleanUpProject(t, client, bt.ProjectID)
}

func TestTrigger_DeleteTriggerScheduleDaily(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := setup()

	bt := createTestBuildTypeWithName(t, client, "BuildTriggerProject", "BuildRelease", true)

	sut := client.TriggerService(bt.ID)
	nt, _ := teamcity.NewTriggerScheduleDaily(bt.ID, 12, 0, "SERVER", []string{"+:*"})

	created, err := sut.AddTrigger(nt)

	require.Nil(err)

	sut.Delete(created.ID())
	_, err = sut.GetByID(created.ID()) // refresh

	require.Error(err)
	assert.Contains(err.Error(), "404")
	cleanUpProject(t, client, bt.ProjectID)
}
