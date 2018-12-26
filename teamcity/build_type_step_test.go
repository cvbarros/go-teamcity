package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity-sdk/teamcity/acctest"

	"github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/stretchr/testify/suite"
)

type SuiteBuildTypeSteps struct {
	suite.Suite
	Project               *teamcity.Project
	BuildType             *teamcity.BuildType
	Client                *teamcity.Client
	StepPowershell        teamcity.Step
	StepCmdLineExecutable teamcity.Step
	StepCmdLineScript     teamcity.Step
	AddStep               func(teamcity.Step) teamcity.Step
}

func (suite *SuiteBuildTypeSteps) SetupSuite() {
	suite.Client = setup()
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

	suite.Project = createTestProject(suite.T(), suite.Client, acctest.RandomWithPrefix("Project_SuiteBuildTypeSteps"))
	suite.BuildType = createTestBuildTypeWithName(suite.T(), suite.Client, suite.Project.ID, acctest.RandomWithPrefix("BuildType_SuiteBuildTypeSteps"), false)

	suite.AddStep = func(s teamcity.Step) (created teamcity.Step) {
		created, err := suite.Client.BuildTypes.AddStep(suite.BuildType.ID, s)
		suite.Require().NoError(err)
		suite.Require().NotNil(created)
		return
	}
}

func (suite *SuiteBuildTypeSteps) TearDownTest() {
	cleanUpProject(suite.T(), suite.Client, suite.Project.ID)
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

func (suite *SuiteBuildTypeSteps) TestGet_All() {
	step1 := suite.AddStep(suite.StepCmdLineScript)
	step2 := suite.AddStep(suite.StepCmdLineScript)

	actual, err := suite.Client.BuildTypes.GetSteps(suite.BuildType.ID)
	suite.NoError(err)
	suite.Contains(actual, step1)
	suite.Contains(actual, step2)
}

func (suite *SuiteBuildTypeSteps) TestDelete() {
	step1 := suite.AddStep(suite.StepCmdLineScript)

	suite.Client.BuildTypes.DeleteStep(suite.BuildType.ID, step1.GetID())

	actual, err := suite.Client.BuildTypes.GetSteps(suite.BuildType.ID)

	suite.NoError(err)
	suite.Empty(actual)
}

func (suite *SuiteBuildTypeSteps) TestGet_Inline() {
	step1 := suite.AddStep(suite.StepCmdLineScript)
	step2 := suite.AddStep(suite.StepCmdLineScript)
	expected := []teamcity.Step{step1, step2}

	actual, err := suite.Client.BuildTypes.GetByID(suite.BuildType.ID) // refresh

	suite.Require().NoError(err)
	suite.Require().NotNil(actual.Steps)
	suite.NotEmpty(actual.Steps)
	suite.Equal(actual.Steps, expected)
}

func (suite *SuiteBuildTypeSteps) TestAdd_Inline() {
	bt := suite.BuildType
	newSteps := []teamcity.Step{suite.StepCmdLineScript, suite.StepPowershell}
	bt.Steps = append(bt.Steps, newSteps[0])
	bt.Steps = append(bt.Steps, newSteps[1])

	actual, err := suite.Client.BuildTypes.Update(bt)
	suite.Require().NoError(err)
	suite.Require().NotNil(actual)

	suite.Equal(len(newSteps), len(actual.Steps))
	for i := 0; i < len(newSteps); i++ {
		suite.Equal(newSteps[i].GetName(), actual.Steps[i].GetName())
		suite.Equal(newSteps[i].Type(), actual.Steps[i].Type())
	}
}

func TestSuiteBuildTypeSteps(t *testing.T) {
	suite.Run(t, new(SuiteBuildTypeSteps))
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
