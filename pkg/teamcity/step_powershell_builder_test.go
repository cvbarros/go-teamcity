package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PropertyAssertions struct {
	a *assert.Assertions
	t *testing.T
}

const testStepName = "stepName"

func TestPowershellBuilderWithScriptFile(t *testing.T) {
	assert := newPropertyAssertions(t)

	psBuilder := teamcity.StepPowershellBuilder

	actual := psBuilder.ScriptFile("build.ps1").Build(testStepName)

	assert.a.Equal(teamcity.StepTypes.Powershell, actual.Type)
	assert.a.Equal(testStepName, actual.Name)
	assert.a.NotNilf(actual.Properties, "Properties expected to be defined")

	assert.assertPropertyValue(actual.Properties, "jetbrains_powershell_script_mode", "FILE")
}

func TestPowershellBuilderWithArgs(t *testing.T) {
	assert := newPropertyAssertions(t)
	expected := "-Target pullrequest"

	psBuilder := teamcity.StepPowershellBuilder

	actual := psBuilder.ScriptFile("build.ps1").Args(expected).Build(testStepName)

	assert.assertPropertyValue(actual.Properties, "jetbrains_powershell_scriptArguments", expected)
}

func TestPowershellScriptAsCode(t *testing.T) {
	assert := newPropertyAssertions(t)
	expected := "Some script code"

	psBuilder := teamcity.StepPowershellBuilder

	actual := psBuilder.Code(expected).Build(testStepName)

	assert.assertPropertyValue(actual.Properties, "jetbrains_powershell_script_code", expected)
	assert.assertPropertyValue(actual.Properties, "jetbrains_powershell_script_mode", "CODE")
}

func newPropertyAssertions(t *testing.T) *PropertyAssertions {
	return &PropertyAssertions{a: assert.New(t), t: t}
}

func (p *PropertyAssertions) assertPropertyValue(props *teamcity.Properties, name string, value string) {
	require.NotNil(p.t, props)

	propMap := props.Map()

	p.a.Contains(propMap, name)
	p.a.Equal(value, propMap[name])
}
