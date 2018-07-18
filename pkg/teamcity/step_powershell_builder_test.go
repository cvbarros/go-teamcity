package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
)

const testStepName = "stepName"

func TestPowershellStepBuilder_ScriptFile(t *testing.T) {
	assert := newPropertyAssertions(t)

	psBuilder := teamcity.StepPowershellBuilder

	actual := psBuilder.ScriptFile("build.ps1").Build(testStepName)

	assert.a.Equal(teamcity.StepTypes.Powershell, actual.Type)
	assert.a.Equal(testStepName, actual.Name)
	assert.a.NotNilf(actual.Properties, "Properties expected to be defined")

	assert.assertPropertyValue(actual.Properties, "jetbrains_powershell_script_mode", "FILE")
	assert.assertPropertyValue(actual.Properties, "jetbrains_powershell_script_file", "build.ps1")
	assert.assertPropertyDoesNotExist(actual.Properties, "jetbrains_powershell_scriptArguments")
}

func TestPowershellStepBuilder_Args(t *testing.T) {
	assert := newPropertyAssertions(t)
	expected := "-Target pullrequest"

	psBuilder := teamcity.StepPowershellBuilder

	actual := psBuilder.ScriptFile("build.ps1").Args(expected).Build(testStepName)

	assert.assertPropertyValue(actual.Properties, "jetbrains_powershell_scriptArguments", expected)
}

func TestPowershellStepBuilder_Code(t *testing.T) {
	assert := newPropertyAssertions(t)
	expected := "Some script code"

	psBuilder := teamcity.StepPowershellBuilder

	actual := psBuilder.Code(expected).Build(testStepName)

	assert.assertPropertyValue(actual.Properties, "jetbrains_powershell_script_code", expected)
	assert.assertPropertyValue(actual.Properties, "jetbrains_powershell_script_mode", "CODE")
}

func TestPowershellStepBuilder_MultipleTimes(t *testing.T) {
	assert := newPropertyAssertions(t)
	expected := "Some script code"

	psBuilder := teamcity.StepPowershellBuilder.ScriptFile("script.ps1")
	psBuilder = psBuilder.Args("someargs")

	actual1 := psBuilder.Build("step1")

	assert.assertPropertyValue(actual1.Properties, "jetbrains_powershell_script_mode", "FILE")
	assert.assertPropertyValue(actual1.Properties, "jetbrains_powershell_script_file", "script.ps1")
	assert.assertPropertyValue(actual1.Properties, "jetbrains_powershell_scriptArguments", "someargs")

	actual2 := teamcity.StepPowershellBuilder.Code(expected).Build("step2")

	assert.assertPropertyValue(actual2.Properties, "jetbrains_powershell_script_code", expected)
	assert.assertPropertyValue(actual2.Properties, "jetbrains_powershell_script_mode", "CODE")
	assert.assertPropertyDoesNotExist(actual2.Properties, "jetbrains_powershell_scriptArguments")
	assert.assertPropertyDoesNotExist(actual2.Properties, "jetbrains_powershell_script_file")
}
