package teamcity_test

import (
	"testing"
	"time"

	teamcity "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/cvbarros/go-teamcity-sdk/teamcity/acctest"
	"github.com/stretchr/testify/suite"
)

type SuiteBuildTypeTrigger struct {
	suite.Suite
	Project    *teamcity.Project
	BuildType  *teamcity.BuildType
	Client     *teamcity.Client
	TriggerVcs teamcity.Trigger
	Trigger    teamcity.Trigger
	AddTrigger func(teamcity.Trigger) teamcity.Trigger
}

func (suite *SuiteBuildTypeTrigger) SetupSuite() {
	suite.Client = setup()
	suite.TriggerVcs, _ = teamcity.NewTriggerVcs([]string{"+:*"}, []string{})
}

func (suite *SuiteBuildTypeTrigger) SetupTest() {

	suite.Project = createTestProject(suite.T(), suite.Client, acctest.RandomWithPrefix("Project_SuiteBuildTypeTrigger"))
	suite.BuildType = suite.NewBuildType("BuildType_SuiteBuildTypeTrigger")

	suite.AddTrigger = func(t teamcity.Trigger) (created teamcity.Trigger) {
		created, err := suite.Client.TriggerService(suite.BuildType.ID).AddTrigger(t)
		suite.Require().NoError(err)
		suite.Require().NotNil(created)
		return
	}
}

func (suite *SuiteBuildTypeTrigger) TearDownTest() {
	cleanUpProject(suite.T(), suite.Client, suite.Project.ID)
}

func (suite *SuiteBuildTypeTrigger) TestVcsTrigger_Create() {
	actual := suite.AddTrigger(suite.TriggerVcs)
	suite.Equal(suite.BuildType.ID, actual.BuildTypeID())
}

func (suite *SuiteBuildTypeTrigger) TestVcsTrigger_Get() {
	nt := suite.AddTrigger(suite.TriggerVcs)
	suite.RefreshTrigger(nt.ID())

	actual := suite.Trigger
	suite.Equal(nt.ID(), actual.ID())
	suite.Equal(nt.BuildTypeID(), actual.BuildTypeID())
	suite.Equal(nt.Type(), actual.Type())
}

func (suite *SuiteBuildTypeTrigger) TestVcsTrigger_Delete() {
	nt := suite.AddTrigger(suite.TriggerVcs)
	suite.RefreshTrigger(nt.ID())
	suite.AssertDeleted()
}

func (suite *SuiteBuildTypeTrigger) TestBuildFinishTrigger_Create() {
	s := suite.NewBuildType("SourceBuildType_SuiteBuildTypeTrigger")
	t := suite.TriggerBuildFinish(s.ID)
	actual := suite.AddTrigger(t)

	suite.Equal(suite.BuildType.ID, actual.BuildTypeID())
	suite.Equal(teamcity.BuildTriggerBuildFinish, actual.Type())
}

func (suite *SuiteBuildTypeTrigger) TestBuildFinishTrigger_Delete() {
	s := suite.NewBuildType("SourceBuildType_SuiteBuildTypeTrigger")
	t := suite.TriggerBuildFinish(s.ID)
	nt := suite.AddTrigger(t)
	suite.RefreshTrigger(nt.ID())
	suite.AssertDeleted()
}

func (suite *SuiteBuildTypeTrigger) TestBuildFinishTrigger_Get() {
	s := suite.NewBuildType("SourceBuildType_SuiteBuildTypeTrigger")
	t := suite.TriggerBuildFinish(s.ID)
	nt := suite.AddTrigger(t)
	suite.RefreshTrigger(nt.ID())

	actual := suite.Trigger
	suite.Equal(nt.ID(), actual.ID())
	suite.Equal(nt.BuildTypeID(), actual.BuildTypeID())
	suite.Equal(nt.Type(), actual.Type())
}

func (suite *SuiteBuildTypeTrigger) TestScheduledDailyTrigger_Create() {
	t := suite.TriggerScheduledDaily(suite.BuildType.ID)
	nt := suite.AddTrigger(t)

	suite.Equal(suite.BuildType.ID, nt.BuildTypeID())
	suite.Equal(teamcity.BuildTriggerSchedule, nt.Type())
	suite.Require().IsType(&teamcity.TriggerSchedule{}, nt)

	actual := nt.(*teamcity.TriggerSchedule)
	suite.Equal(teamcity.TriggerSchedulingDaily, actual.SchedulingPolicy)
	suite.Equal(uint(12), actual.Hour)
	suite.Equal(uint(30), actual.Minute)
	suite.Equal(false, actual.Options.BuildWithPendingChangesOnly)
	suite.Equal(false, actual.Options.PromoteWatchedBuild)
	suite.Equal(false, actual.Options.QueueOptimization)
	suite.Equal(true, actual.Options.BuildOnAllCompatibleAgents)
}

func (suite *SuiteBuildTypeTrigger) TestScheduledDailyTrigger_Delete() {
	t := suite.TriggerScheduledDaily(suite.BuildType.ID)
	nt := suite.AddTrigger(t)
	suite.RefreshTrigger(nt.ID())
	suite.AssertDeleted()
}

func (suite *SuiteBuildTypeTrigger) TestScheduledDailyTrigger_Get() {
	t := suite.TriggerScheduledDaily(suite.BuildType.ID)
	nt := suite.AddTrigger(t)
	suite.RefreshTrigger(nt.ID())

	created := suite.Trigger
	suite.Equal(nt.ID(), created.ID())
	suite.Equal(nt.BuildTypeID(), created.BuildTypeID())
	suite.Equal(nt.Type(), created.Type())
}

func (suite *SuiteBuildTypeTrigger) TestScheduledWeeklyTrigger_Create() {
	t := suite.TriggerScheduledWeekly(suite.BuildType.ID)
	nt := suite.AddTrigger(t)

	suite.Equal(suite.BuildType.ID, nt.BuildTypeID())
	suite.Equal(teamcity.BuildTriggerSchedule, nt.Type())
	suite.Require().IsType(&teamcity.TriggerSchedule{}, nt)

	actual := nt.(*teamcity.TriggerSchedule)
	suite.Equal(actual.SchedulingPolicy, teamcity.TriggerSchedulingWeekly)
	suite.Equal(actual.Weekday, time.Thursday)
}

func (suite *SuiteBuildTypeTrigger) TestScheduledWeeklyTrigger_Delete() {
	t := suite.TriggerScheduledWeekly(suite.BuildType.ID)
	nt := suite.AddTrigger(t)
	suite.RefreshTrigger(nt.ID())
	suite.AssertDeleted()
}

func (suite *SuiteBuildTypeTrigger) TestScheduledWeeklyTrigger_Get() {
	t := suite.TriggerScheduledWeekly(suite.BuildType.ID)
	nt := suite.AddTrigger(t)
	suite.RefreshTrigger(nt.ID())

	created := suite.Trigger
	suite.Equal(nt.ID(), created.ID())
	suite.Equal(nt.BuildTypeID(), created.BuildTypeID())
	suite.Equal(nt.Type(), created.Type())
}

func (suite *SuiteBuildTypeTrigger) AssertDeleted() {
	ts := suite.Client.TriggerService(suite.BuildType.ID)
	ts.Delete(suite.Trigger.ID())
	_, err := ts.GetByID(suite.Trigger.ID()) // refresh

	suite.Require().Error(err)
	suite.Contains(err.Error(), "404")
}

func (suite *SuiteBuildTypeTrigger) RefreshTrigger(id string) {
	actual, err := suite.Client.TriggerService(suite.BuildType.ID).GetByID(id)
	suite.Require().NoError(err)
	suite.Trigger = actual
}

func (suite *SuiteBuildTypeTrigger) TriggerBuildFinish(source string) teamcity.Trigger {
	t, err := teamcity.NewTriggerBuildFinish(source, teamcity.NewTriggerBuildFinishOptions(true, []string{"+:<default>"}))
	suite.Require().NoError(err)
	return t
}

func (suite *SuiteBuildTypeTrigger) TriggerScheduledDaily(source string) teamcity.Trigger {
	nt, _ := teamcity.NewTriggerScheduleDaily(source, 12, 30, "SERVER", []string{"+:*"})
	opt := teamcity.NewTriggerScheduleOptions()
	opt.QueueOptimization = false
	opt.BuildOnAllCompatibleAgents = true
	opt.PromoteWatchedBuild = false
	opt.BuildWithPendingChangesOnly = false
	nt.Options = opt
	return nt
}

func (suite *SuiteBuildTypeTrigger) TriggerScheduledWeekly(source string) teamcity.Trigger {
	nt, err := teamcity.NewTriggerScheduleWeekly(source, time.Thursday, 12, 0, "SERVER", []string{"+:*"})
	suite.Require().NoError(err)
	return nt
}

func (suite *SuiteBuildTypeTrigger) NewBuildType(name string) *teamcity.BuildType {
	return createTestBuildTypeWithName(suite.T(), suite.Client, suite.Project.ID, acctest.RandomWithPrefix(name), false)
}

func TestSuiteBuildTypeTrigger(t *testing.T) {
	suite.Run(t, new(SuiteBuildTypeTrigger))
}
