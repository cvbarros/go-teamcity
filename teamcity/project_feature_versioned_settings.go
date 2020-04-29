package teamcity

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// VersionedSettingsFormat represents the supported formats for the versioned settings project feature.
type VersionedSettingsFormat string

const (
	// VersionedSettingsFormatKotlin represents that the versioned settings are stored as kotlin files.
	VersionedSettingsFormatKotlin VersionedSettingsFormat = "kotlin"
	// VersionedSettingsFormatXML represents that the versioned settings are stored as xml files.
	VersionedSettingsFormatXML VersionedSettingsFormat = "xml"
)

// VersionedSettingsBuildSettings represents the supported strategies for retrieving the current settings.
type VersionedSettingsBuildSettings string

const (
	// VersionedSettingsBuildSettingsPreferCurrent determines that teamcity will default to the settings from the teamcity server.
	VersionedSettingsBuildSettingsPreferCurrent VersionedSettingsBuildSettings = "PREFER_CURRENT"
	// VersionedSettingsBuildSettingsPreferVcs determines that teamcity should always prefer the settings stored in the VCS.
	VersionedSettingsBuildSettingsPreferVcs VersionedSettingsBuildSettings = "PREFER_VCS"
	// VersionedSettingsBuildSettingsAlwaysUseCurrent determines that teamcity should always use the project settings from the server.
	VersionedSettingsBuildSettingsAlwaysUseCurrent VersionedSettingsBuildSettings = "ALWAYS_USE_CURRENT"
)

// CredentialsStorageType represents how credentials should be stored in versioned settings
type CredentialsStorageType string

const (
	// Credentials should be scrambled and stored in VCS
	CredentialsStorageTypeScrambledInVcs CredentialsStorageType = ""

	// Credentials should be stored as Secure Tokens within the Project
	CredentialsStorageTypeCredentialsJSON CredentialsStorageType = "credentialsJSON"
)

// ProjectFeatureVersionedSettingsOptions holds all properties for the versioned settings project feature.
type ProjectFeatureVersionedSettingsOptions struct {
	Enabled                bool
	ShowChanges            bool
	UseRelativeIds         bool
	VcsRootID              string
	Format                 VersionedSettingsFormat
	BuildSettings          VersionedSettingsBuildSettings
	CredentialsStorageType CredentialsStorageType
	ContextParameters      map[string]string
}

// ProjectFeatureVersionedSettings represents the versioned settings feature for a project.
// Can be used to configure https://confluence.jetbrains.com/display/TCD10/Storing+Project+Settings+in+Version+Control.
type ProjectFeatureVersionedSettings struct {
	id        string
	projectID string

	Options ProjectFeatureVersionedSettingsOptions
}

// NewProjectFeatureVersionedSettings creates a new Versioned Settings project feature.
func NewProjectFeatureVersionedSettings(projectID string, options ProjectFeatureVersionedSettingsOptions) *ProjectFeatureVersionedSettings {
	return &ProjectFeatureVersionedSettings{
		projectID: projectID,
		Options:   options,
	}
}

// ID returns the ID of this project feature.
func (f *ProjectFeatureVersionedSettings) ID() string {
	return f.id
}

// SetID sets the ID of this project feature.
func (f *ProjectFeatureVersionedSettings) SetID(value string) {
	f.id = value
}

// Type represents the type of the project feature as a string.
func (f *ProjectFeatureVersionedSettings) Type() string {
	return "versionedSettings"
}

// ProjectID represents the ID of the project the project feature is assigned to.
func (f *ProjectFeatureVersionedSettings) ProjectID() string {
	return f.projectID
}

// SetProjectID sets the ID of the project the project feature is assigned to.
func (f *ProjectFeatureVersionedSettings) SetProjectID(value string) {
	f.projectID = value
}

// Properties returns all properties for the versioned settings project feature.
func (f *ProjectFeatureVersionedSettings) Properties() *Properties {
	props := NewProperties(
		NewProperty("buildSettings", string(f.Options.BuildSettings)),
		NewProperty("format", string(f.Options.Format)),
		NewProperty("rootId", f.Options.VcsRootID),
		NewProperty("showChanges", fmt.Sprintf("%t", f.Options.ShowChanges)),
		NewProperty("useRelativeIds", fmt.Sprintf("%t", f.Options.UseRelativeIds)),
		NewProperty("enabled", fmt.Sprintf("%t", f.Options.Enabled)),
	)

	// TeamCity doesn't send this if unset
	if f.Options.CredentialsStorageType != "" {
		props.Add(NewProperty("credentialsStorageType", string(f.Options.CredentialsStorageType)))
	}

	// these are `context.foo` rather than a nested object
	for k, v := range f.Options.ContextParameters {
		key := fmt.Sprintf("context.%s", k)
		props.Add(NewProperty(key, v))
	}

	return props
}

func loadProjectFeatureVersionedSettings(projectID string, feature projectFeatureJSON) (ProjectFeature, error) {
	settings := &ProjectFeatureVersionedSettings{
		id:        feature.ID,
		projectID: projectID,
		Options:   ProjectFeatureVersionedSettingsOptions{},
	}

	if encodedValue, ok := feature.Properties.GetOk("buildSettings"); ok {
		settings.Options.BuildSettings = VersionedSettingsBuildSettings(encodedValue)
	}

	// Context Parameters are a map, rather than a nested object
	// as the key `context.foo` and `context.bar`
	contextParams := make(map[string]string)
	for k, value := range feature.Properties.Map() {
		if !strings.HasPrefix(k, "context.") {
			continue
		}

		key := strings.TrimPrefix(k, "context.")
		contextParams[key] = value
	}
	settings.Options.ContextParameters = contextParams

	if encodedValue, ok := feature.Properties.GetOk("credentialsStorageType"); ok {
		settings.Options.CredentialsStorageType = CredentialsStorageType(encodedValue)
	} else {
		// defaulted but not returned
		settings.Options.CredentialsStorageType = CredentialsStorageTypeScrambledInVcs
	}

	if encodedValue, ok := feature.Properties.GetOk("enabled"); ok {
		v, err := strconv.ParseBool(encodedValue)
		if err != nil {
			return nil, err
		}

		settings.Options.Enabled = v
	}

	if encodedValue, ok := feature.Properties.GetOk("format"); ok {
		settings.Options.Format = VersionedSettingsFormat(encodedValue)
	}

	if encodedValue, ok := feature.Properties.GetOk("rootId"); ok {
		settings.Options.VcsRootID = encodedValue
	}

	if encodedValue, ok := feature.Properties.GetOk("showChanges"); ok {
		value := false
		if err := json.Unmarshal([]byte(encodedValue), &value); err != nil {
			return nil, err
		}
		settings.Options.ShowChanges = value
	}

	if encodedValue, ok := feature.Properties.GetOk("useRelativeIds"); ok {
		value := false
		if err := json.Unmarshal([]byte(encodedValue), &value); err != nil {
			return nil, err
		}
		settings.Options.UseRelativeIds = value
	}

	return settings, nil
}
