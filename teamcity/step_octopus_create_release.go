package teamcity

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// StepOctopusPushPackage represents a a build step of type "octopus.create.release"
type StepOctopusCreateRelease struct {
	ID       string
	Name     string
	stepType string
	stepJSON *stepJSON

	// Specify Octopus web portal URL.
	Host string

	// Specify Octopus API key.
	ApiKey string

	// Specify which version of the Octopus Deploy server you are using.
	OctopusServerVersion string

	// Name of the Octopus project to create a release for.
	Project string

	// Number to use for this release.
	ReleaseNumber string

	// Channel to create the release for.
	ChannelName string

	// Comma separated list of environments to deploy to.
	// Leave empty to create a release without deploying it.
	Environments string

	// Comma separated list of tenants to promote for.
	// Wildcard '*' will promote all tenants currently able to deploy to the above provided environment.
	// Note that when supplying tenant filters then only one environment may be provided above.
	Tenants string

	// Comma separated list of tenant tags that match tenants to deploy for.
	// Note that when supplying tag filters then only one environment may be provided above.
	TenantTags string

	//  If true, the build process will only succeed if the deployment is successful.
	// Output from the deployment will appear in the build output.
	WaitForDeployments bool

	// Additional arguments to be passed to Octo.exe.
	AdditionalCommandLineArguments string
}

func NewStepOctopusCreateRelease(name string) (*StepOctopusCreateRelease, error) {
	return &StepOctopusCreateRelease{
		Name:     name,
		stepType: StepTypeOctopusCreateRelease,
	}, nil
}

func (s *StepOctopusCreateRelease) GetID() string {
	return s.ID
}

func (s *StepOctopusCreateRelease) GetName() string {
	return s.Name
}

func (s *StepOctopusCreateRelease) Type() BuildStepType {
	return StepTypeOctopusCreateRelease
}

func (s *StepOctopusCreateRelease) properties() *Properties {
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("teamcity.step.mode", "default")
	props.AddOrReplaceValue("octopus_host", s.Host)
	props.AddOrReplaceValue("secure:octopus_apikey", s.ApiKey)
	props.AddOrReplaceValue("octopus_version", s.OctopusServerVersion)
	props.AddOrReplaceValue("octopus_project_name", s.Project)
	props.AddOrReplaceValue("octopus_releasenumber", s.ReleaseNumber)
	props.AddOrReplaceValue("octopus_channel_name", s.ChannelName)
	props.AddOrReplaceValue("octopus_deployto", s.Environments)
	props.AddOrReplaceValue("octoups_tenants", s.Tenants)
	props.AddOrReplaceValue("octoups_tenanttags", s.TenantTags)
	props.AddOrReplaceValue("octopus_waitfordeployments", strconv.FormatBool(s.WaitForDeployments))
	props.AddOrReplaceValue("octopus_additionalcommandlinearguments", s.AdditionalCommandLineArguments)

	return props
}

func (s *StepOctopusCreateRelease) MarshalJSON() ([]byte, error) {
	out := &stepJSON{
		ID:         s.ID,
		Name:       s.Name,
		Type:       s.stepType,
		Properties: s.properties(),
	}

	return json.Marshal(out)
}

// UnmarshalJSON implements JSON deserialization for StepOctopusPushPackage
func (s *StepOctopusCreateRelease) UnmarshalJSON(data []byte) error {
	var aux stepJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != string(StepTypeOctopusCreateRelease) {
		return fmt.Errorf("invalid type %s trying to deserialize into StepTypeOctopusCreateRelease entity", aux.Type)
	}
	s.Name = aux.Name
	s.ID = aux.ID
	s.stepType = StepTypeOctopusCreateRelease

	props := aux.Properties
	if v, ok := props.GetOk("octopus_host"); ok {
		s.Host = v
	}

	if v, ok := props.GetOk("secure:octopus_apikey"); ok {
		s.ApiKey = v
	}

	if v, ok := props.GetOk("octopus_version"); ok {
		s.OctopusServerVersion = v
	}

	if v, ok := props.GetOk("octopus_project_name"); ok {
		s.Project = v
	}

	if v, ok := props.GetOk("octopus_releasenumber"); ok {
		s.ReleaseNumber = v
	}

	if v, ok := props.GetOk("octopus_channel_name"); ok {
		s.ChannelName = v
	}

	if v, ok := props.GetOk("octopus_deployto"); ok {
		s.Environments = v
	}

	if v, ok := props.GetOk("octoups_tenants"); ok {
		s.Tenants = v
	}

	if v, ok := props.GetOk("octoups_tenanttags"); ok {
		s.TenantTags = v
	}

	if v, ok := props.GetOk("octopus_waitfordeployments"); ok {
		converted_value, _ := strconv.ParseBool(v)
		s.WaitForDeployments = converted_value
	}

	if v, ok := props.GetOk("octopus_additionalcommandlinearguments"); ok {
		s.AdditionalCommandLineArguments = v
	}

	return nil
}
