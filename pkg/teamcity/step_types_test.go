package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
)

type PropertyAssertions struct {
	a *assert.Assertions
}

func TestPowershellBuilderWithScriptFile(t *testing.T) {
	assert := &PropertyAssertions{a: assert.New(t)}

	psBuilder := teamcity.PowershellStepBuilder

	actual := psBuilder.ScriptFile("build.ps1").Build("stepName")

	assert.a.Equal(teamcity.StepTypes.Powershell, actual.Type)
	assert.a.Equal("stepName", actual.Name)
	assert.a.NotNilf(actual.Properties, "Properties expected to be defined")

	assert.assertPropertyValue(actual.Properties, "jetbrains_powershell_script_mode", "FILE")
}

func TestPowershellBuilderWithArgs(t *testing.T) {
	assert := &PropertyAssertions{a: assert.New(t)}
	expected := "-Target pullrequest"

	psBuilder := teamcity.PowershellStepBuilder

	actual := psBuilder.ScriptFile("build.ps1").Args(expected).Build("stepName")

	assert.assertPropertyValue(actual.Properties, "jetbrains_powershell_scriptArguments", expected)
}

func (p *PropertyAssertions) assertPropertyValue(props *teamcity.Properties, name string, value string) {
	propMap := props.Map()

	p.a.Contains(propMap, name)
	p.a.Equal(propMap[name], value)
}
