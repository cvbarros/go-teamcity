package teamcity

import "testing"

func Test_PropertiesAllowTriggeringPersonalBuilds_OmitWhenTrue(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults(nil)

	sut.AllowPersonalBuildTriggering = true
	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "allowPersonalBuildTriggering")
}

func Test_PropertiesBuildConfigurationType_OmitWhenDefault(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults(nil)

	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "buildConfigurationType")
}

func Test_PropertiesBuildNumberFormat_OmitWhenDefault(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults(nil)

	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "buildNumberPattern")
}

func Test_PropertiesEnableStatusWidget_OmitWhenFalse(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults(nil)
	sut.EnableStatusWidget = false //Default, but explicit
	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "enableExternalStatus")
}

func Test_PropertiesEnableHangingBuildsDetection_OmitWhenTrue(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults(nil)
	sut.EnableHangingBuildsDetection = true //Default, but explicit
	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "enableHangingBuildsDetection")
}

func Test_PropertiesMaxSimultaneousBuilds_OmitWhenZero(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults(nil)
	sut.MaxSimultaneousBuilds = 0 //Default, but explicit
	actual := sut.properties()

	pa.assertPropertyDoesNotExist(actual, "maximumNumberOfBuilds")
}

func Test_Properties_Full(t *testing.T) {
	pa := newPropertyAssertions(t)
	sut := NewBuildTypeOptionsWithDefaults(nil)

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

	pa.assertPropertyDoesNotExist(actual, "allowPersonalBuildTriggering")
	pa.assertPropertyValue(actual, "buildConfigurationType", "DEPLOYMENT")
	pa.assertPropertyValue(actual, "enableHangingBuildsDetection", "false")
	pa.assertPropertyValue(actual, "allowExternalStatus", "true")
	pa.assertPropertyValue(actual, "buildNumberPattern", "1.0.%build.counter%")
	pa.assertPropertyValue(actual, "artifactRules", "abc\ndef")
	pa.assertPropertyValue(actual, "maximumNumberOfBuilds", "10")
	pa.assertPropertyValue(actual, "buildNumberCounter", "5")
}
