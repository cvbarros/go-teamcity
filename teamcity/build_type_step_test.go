package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/stretchr/testify/suite"
)

type SuiteBuildTypeSteps struct {
	suite.Suite
	TC                    *TestContext
	BuildTypeContext      *BuildTypeContext
	BuildTypeID           string
	StepPowershell        teamcity.Step
	StepCmdLineExecutable teamcity.Step
	StepCmdLineScript     teamcity.Step
	AddStep               func(teamcity.Step) teamcity.Step
}

func NewSuiteBuildTypeSteps(t *testing.T) *SuiteBuildTypeSteps {
	return &SuiteBuildTypeSteps{TC: NewTc("SuiteBuildTypeSteps", t), BuildTypeContext: new(BuildTypeContext)}
}

func (suite *SuiteBuildTypeSteps) SetupSuite() {
	suite.StepPowershell, _ = teamcity.NewStepPowershellScriptFile("step1", "build.ps1", "")
	suite.StepCmdLineExecutable, _ = teamcity.NewStepCommandLineExecutable("step_exe", "./script.sh", "hello")
	suite.StepCmdLineScript, _ = teamcity.NewStepCommandLineExecutable("step_exe", "./script.sh", "hello")
	script := `echo "Hello World
	echo "World, Hello!
	export HELLO_WORLD=1
	`
	suite.StepCmdLineScript, _ = teamcity.NewStepCommandLineScript("step_exe", script)
}

func (suite *SuiteBuildTypeSteps) SetupTest() {
	suite.BuildTypeContext.Setup(suite.TC)
	suite.BuildTypeID = suite.BuildTypeContext.BuildType.ID
	suite.AddStep = func(s teamcity.Step) (created teamcity.Step) {
		created, err := suite.TC.Client.BuildTypes.AddStep(suite.BuildTypeID, s)
		suite.Require().NoError(err)
		suite.Require().NotNil(created)
		return
	}
}

func (suite *SuiteBuildTypeSteps) TearDownTest() {
	suite.BuildTypeContext.Teardown()
}

func (suite *SuiteBuildTypeSteps) TestAdd_StepPowershell() {
	suite.AddStep(suite.StepPowershell)
}

func (suite *SuiteBuildTypeSteps) TestAdd_StepCmdLineExecutable() {
	suite.AddStep(suite.StepCmdLineExecutable)
}

func (suite *SuiteBuildTypeSteps) TestAdd_StepCmdLineScript() {
	suite.AddStep(suite.StepCmdLineScript)
}

func (suite *SuiteBuildTypeSteps) GetSteps(buildTypeID string) []teamcity.Step {
	out, err := suite.TC.Client.BuildTypes.GetSteps(suite.BuildTypeID)
	suite.Require().NoError(err)
	return out
}

func (suite *SuiteBuildTypeSteps) TestGet_All() {
	step1 := suite.AddStep(suite.StepCmdLineScript)
	step2 := suite.AddStep(suite.StepCmdLineScript)

	actual := suite.GetSteps(suite.BuildTypeID)
	suite.Contains(actual, step1)
	suite.Contains(actual, step2)
}

func (suite *SuiteBuildTypeSteps) TestDelete() {
	step1 := suite.AddStep(suite.StepCmdLineScript)
	sut := suite.TC.Client.BuildTypes
	sut.DeleteStep(suite.BuildTypeID, step1.GetID())

	actual := suite.GetSteps(suite.BuildTypeID)
	suite.Empty(actual)
}

func (suite *SuiteBuildTypeSteps) TestGet_Inline() {
	step1 := suite.AddStep(suite.StepCmdLineScript)
	step2 := suite.AddStep(suite.StepCmdLineScript)
	expected := []teamcity.Step{step1, step2}

	actual := suite.GetSteps(suite.BuildTypeID) // refresh

	suite.NotEmpty(actual)
	suite.Equal(expected, actual)
}

func (suite *SuiteBuildTypeSteps) TestAdd_Inline() {
	bt := suite.BuildTypeContext.BuildType
	newSteps := []teamcity.Step{suite.StepCmdLineScript, suite.StepPowershell}
	bt.Steps = append(bt.Steps, newSteps[0])
	bt.Steps = append(bt.Steps, newSteps[1])
	sut := suite.TC.Client.BuildTypes

	actual, err := sut.Update(bt)
	suite.Require().NoError(err)
	suite.Require().NotNil(actual)

	suite.Equal(len(newSteps), len(actual.Steps))
	for i := 0; i < len(newSteps); i++ {
		suite.Equal(newSteps[i].GetName(), actual.Steps[i].GetName())
		suite.Equal(newSteps[i].Type(), actual.Steps[i].Type())
	}
}

func TestSuiteBuildTypeSteps(t *testing.T) {
	s := NewSuiteBuildTypeSteps(t)
	suite.Run(t, s)
}

func createTestBuildStep(t *testing.T, client *teamcity.Client, step teamcity.Step, buildTypeProjectId string) (*teamcity.BuildType, teamcity.Step) {
	createdBuildType := createTestBuildType(t, client, buildTypeProjectId)

	created, err := client.BuildTypes.AddStep(createdBuildType.ID, step)
	if err != nil {
		t.Fatalf("Failed to add step to buildType '%s'", createdBuildType.ID)
	}

	updated, _ := client.BuildTypes.GetByID(createdBuildType.ID)
	return updated, created
}
