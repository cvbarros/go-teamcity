package teamcity

// ProjectFeatureGoogleCloudProfileOptions represents the options available on Google Cloud Profiles
type ProjectFeatureGoogleCloudProfileOptions struct {
	Enabled             bool   `prop:"enabled"`
	ProfileID           string `prop:"profileId"`
	Name                string `prop:"name"`
	Description         string `prop:"description"`
	CloudCode           string `prop:"cloud-code"`
	ProfileServerURL    string `prop:"profileServerUrl"`
	AgentPushPreset     string `prop:"agentPushPreset"` //
	TotalWorkTime       int    `prop:"total-work-time"`
	CredentialsType     string `prop:"credentialsType"` // TODO "key"
	NextHour            string `prop:"next-hour"`       // TODO
	TerminateAfterBuild bool   `prop:"terminate-after-build"`
	TerminateIdleTime   int    `prop:"terminate-idle-time"`
	AccessKey           string `prop:"secure:accessKey"`
}

// ProjectFeatureGoogleCloudProfile represents the Google Cloud Profile feature for a project.
type ProjectFeatureGoogleCloudProfile struct {
	id        string
	projectID string

	Options ProjectFeatureGoogleCloudProfileOptions
}

// NewProjectFeatureGoogleCloudProfile creates a new Google Cloud Profile project feature.
func NewProjectFeatureGoogleCloudProfile(projectID string, options ProjectFeatureGoogleCloudProfileOptions) *ProjectFeatureGoogleCloudProfile {
	options.CloudCode = "google"

	return &ProjectFeatureGoogleCloudProfile{
		projectID: projectID,
		Options:   options,
	}
}

// ID returns the Cloud Profile ID
func (f *ProjectFeatureGoogleCloudProfile) ID() string {
	return f.id
}

// SetID sets the Cloud Profile ID
func (f *ProjectFeatureGoogleCloudProfile) SetID(value string) {
	f.id = value
}

// Type returns the Feature type, which is always "CloudProfile"
func (f *ProjectFeatureGoogleCloudProfile) Type() string {
	return "CloudProfile"
}

// ProjectID returns the Project ID for the project to which this profile is attached
func (f *ProjectFeatureGoogleCloudProfile) ProjectID() string {
	return f.projectID
}

// SetProjectID sets the Project ID for the project to which this profile should be attached
func (f *ProjectFeatureGoogleCloudProfile) SetProjectID(value string) {
	f.projectID = value
}

// Properties returns a typed representation of the Properties of this profile
func (f *ProjectFeatureGoogleCloudProfile) Properties() *Properties {
	f.Options.CloudCode = "google"

	return serializeToProperties(&f.Options)
}

func loadProjectFeatureGoogleCloudProfile(projectID string, feature projectFeatureJSON) (ProjectFeature, error) {
	settings := &ProjectFeatureGoogleCloudProfile{
		id:        feature.ID,
		projectID: projectID,
		Options:   ProjectFeatureGoogleCloudProfileOptions{},
	}

	fillStructFromProperties(&settings.Options, feature.Properties)
	return settings, nil
}
