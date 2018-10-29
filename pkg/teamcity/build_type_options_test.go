package teamcity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PropertiesArtifactRules_EmitEvenWhenEmpty(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults()

	actual := sut.properties()
	pa.assertPropertyValue(actual, "artifactRules", "")
}

func Test_PropertiesAllowTriggeringPersonalBuilds_OmitWhenTrue(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults()

	sut.AllowPersonalBuildTriggering = true
	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "allowPersonalBuildTriggering")
}

func Test_ConvertFromPropertiesAllowTriggringPersonalBuilds_SetsDefaultTrue(t *testing.T) {
	assert := assert.New(t)
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("someProperty", "someValue")

	actual := props.buildTypeOptions(false)

	assert.Equal(actual.AllowPersonalBuildTriggering, true)
}

func Test_PropertiesBuildConfigurationType_OmitWhenDefault(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults()

	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "buildConfigurationType")
}

func Test_ConvertFromPropertiesBuildConfigurationType_SetsDefault(t *testing.T) {
	assert := assert.New(t)
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("someProperty", "someValue")

	actual := props.buildTypeOptions(false)

	assert.Equal(actual.BuildConfigurationType, DefaultBuildConfigurationType)
}

func Test_PropertiesBuildNumberFormat_OmitWhenDefault(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults()

	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "buildNumberPattern")
}

func Test_ConvertFromPropertiesBuildNumberFormat_SetsDefault(t *testing.T) {
	assert := assert.New(t)
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("someProperty", "someValue")

	actual := props.buildTypeOptions(false)

	assert.Equal(actual.BuildNumberFormat, DefaultBuildNumberFormat)
}

func Test_PropertiesEnableStatusWidget_OmitWhenFalse(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults()
	sut.EnableStatusWidget = false //Default, but explicit
	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "enableExternalStatus")
}
func Test_ConvertFromPropertiesEnableStatusWidget_SetsDefault(t *testing.T) {
	assert := assert.New(t)
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("someProperty", "someValue")

	actual := props.buildTypeOptions(false)

	assert.Equal(actual.EnableStatusWidget, false)
}

func Test_PropertiesEnableHangingBuildsDetection_OmitWhenTrue(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults()
	sut.EnableHangingBuildsDetection = true //Default, but explicit
	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "enableHangingBuildsDetection")
}

func Test_ConvertFromPropertiesEnableHangingBuildsDetection_SetsDefault(t *testing.T) {
	assert := assert.New(t)
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("someProperty", "someValue")

	actual := props.buildTypeOptions(false)

	assert.Equal(actual.EnableStatusWidget, false)
}

func Test_PropertiesMaxSimultaneousBuilds_OmitWhenZero(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults()
	sut.MaxSimultaneousBuilds = 0 //Default, but explicit
	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "maximumNumberOfBuilds")
}

func Test_ConvertFromPropertiesMaxSimultaneousBuilds_SetsDefault(t *testing.T) {
	assert := assert.New(t)
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("someProperty", "someValue")

	actual := props.buildTypeOptions(false)

	assert.Equal(actual.MaxSimultaneousBuilds, 0)
}

func Test_PropertiesIfTemplate_OmitBuildCounter(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults()
	sut.Template = true
	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "buildNumberCounter")
}

func Test_Properties_Full(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults()

	//Set all to non-default values
	sut.AllowPersonalBuildTriggering = !sut.AllowPersonalBuildTriggering
	sut.BuildConfigurationType = "DEPLOYMENT"
	sut.EnableHangingBuildsDetection = !sut.EnableHangingBuildsDetection
	sut.EnableStatusWidget = !sut.EnableStatusWidget
	sut.BuildNumberFormat = "1.0.%build.counter%"
	sut.ArtifactRules = []string{"abc", "def"}
	sut.MaxSimultaneousBuilds = 10
	sut.BuildCounter = 5

	actual := sut.properties()

	pa.assertPropertyValue(actual, "allowPersonalBuildTriggering", "false")
	pa.assertPropertyValue(actual, "buildConfigurationType", "DEPLOYMENT")
	pa.assertPropertyValue(actual, "enableHangingBuildsDetection", "false")
	pa.assertPropertyValue(actual, "allowExternalStatus", "true")
	pa.assertPropertyValue(actual, "buildNumberPattern", "1.0.%build.counter%")
	pa.assertPropertyValue(actual, "artifactRules", "abc\ndef")
	pa.assertPropertyValue(actual, "maximumNumberOfBuilds", "10")
	pa.assertPropertyValue(actual, "buildNumberCounter", "5")
}
