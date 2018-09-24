package teamcity

const DefaultBuildNumberFormat = "%build.counter%"
const DefaultBuildConfigurationType = "REGULAR"

//BuildTypeOptions represents settings for a Build Configuration
type BuildTypeOptions struct {
	AllowPersonalBuildTriggering bool     `prop:"allowPersonalBuildTriggering"`
	ArtifactRules                []string `prop:"artifactRules" separator:"\n"`
	EnableHangingBuildsDetection bool     `prop:"enableHangingBuildsDetection" force:""`
	EnableStatusWidget           bool     `prop:"allowExternalStatus"`
	BuildCounter                 int      `prop:"buildNumberCounter"`
	BuildNumberFormat            string   `prop:"buildNumberPattern"`
	BuildConfigurationType       string   `prop:"buildConfigurationType"`
	MaxSimultaneousBuilds        int      `prop:"maximumNumberOfBuilds"`

	BuildTypeID int
}

//NewBuildTypeOptionsWithDefaults returns a new instance of default settings, the same as presented in the TeamCity UI when a new build configuration is created.
func NewBuildTypeOptionsWithDefaults(artifactRules []string) *BuildTypeOptions {
	return &BuildTypeOptions{
		AllowPersonalBuildTriggering: true,
		ArtifactRules:                artifactRules,
		EnableHangingBuildsDetection: true,
		EnableStatusWidget:           false,
		MaxSimultaneousBuilds:        0,
		BuildConfigurationType:       DefaultBuildConfigurationType,
		BuildCounter:                 1,
		BuildNumberFormat:            DefaultBuildNumberFormat,
	}
}

func (o *BuildTypeOptions) properties() *Properties {
	props := serializeToProperties(o)

	//TeamCity API for build settings has a very weird behavior to omit some properties when they assume their "default" value.
	//In this case, in order to keep consistent behaviour between reads/writes, the property raw model is adjusted for this behaviour.

	//Omit allowPersonalBuildTriggering if equals to default 'true'
	if o.AllowPersonalBuildTriggering {
		props.Remove("allowPersonalBuildTriggering")
	}

	//Omit enableHangingBuildsDetection if equals to default 'true'
	if o.EnableHangingBuildsDetection {
		props.Remove("enableHangingBuildsDetection")
	}

	//Omit if buildConfigurationType == "REGULAR"
	if v, ok := props.GetOk("buildConfigurationType"); ok && v == DefaultBuildConfigurationType {
		props.Remove("buildConfigurationType")
	}

	//Omit if buildNumberPattern == "%build.counter%"
	if v, ok := props.GetOk("buildNumberPattern"); ok && v == DefaultBuildNumberFormat {
		props.Remove("buildNumberPattern")
	}

	if v, ok := props.GetOk("maximumNumberOfBuilds"); ok && v == "0" {
		props.Remove("maximumNumberOfBuilds")
	}
	return props
}

// func (b buildTypeSettingsBuilder) Build() *Properties {
// 	var props []*Property

// 	props = appendPropertyIfApplicable(b, props, "buildConfigurationType")
// 	props = appendPropertyIfApplicable(b, props, "allowPersonalBuildTriggering")
// 	props = appendPropertyIfApplicable(b, props, "enableHangingBuildsDetection")
// 	props = appendPropertyIfApplicable(b, props, "artifactRules")
// 	props = appendPropertyIfApplicable(b, props, "maximumNumberOfBuilds")
// 	props = appendPropertyIfApplicable(b, props, "buildCounter")
// 	props = appendPropertyIfApplicable(b, props, "buildNumberPattern")

// 	return NewProperties(props...)
// }
