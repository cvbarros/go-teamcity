package teamcity

import (
	"fmt"
	"strconv"
)

// ConnectionType represents the TeamCity connection type
type ConnectionType string

const (
	// ConnectionTypeOAuthProvider resprents an OAuth Provider connection type
	ConnectionTypeOAuthProvider ConnectionType = "OAuthProvider"
)

// ConnectionProviderType represents the OAuth provider type
type ConnectionProviderType string

const (
	// ConnectionProviderTypeVault represents an Hashicorp Vault provider type
	ConnectionProviderTypeVault ConnectionProviderType = "teamcity-vault"
)

// ConnectionProviderVaultAuthMethod represents the Vault auth method
type ConnectionProviderVaultAuthMethod string

const (
	// ConnectionProviderVaultAuthMethodIAM represents the IAM auth method
	ConnectionProviderVaultAuthMethodIAM ConnectionProviderVaultAuthMethod = "iam"
	// ConnectionProviderVaultAuthMethodApprole represents the approle auth method
	ConnectionProviderVaultAuthMethodApprole ConnectionProviderVaultAuthMethod = "approle"
)

// ConnectionProviderVaultOptions is the required options for the Vault OAuth provider
type ConnectionProviderVaultOptions struct {
	AuthMethod     ConnectionProviderVaultAuthMethod
	DisplayName    string
	Endpoint       string
	FailOnError    bool
	Namespace      string
	ProviderType   ConnectionProviderType
	RoleID         string
	SecretID       string
	URL            string
	VaultNamespace string
}

// ConnectionProviderVault defines the Vault Connection details
type ConnectionProviderVault struct {
	id        string
	projectID string

	Options ConnectionProviderVaultOptions
}

// NewProjectConnectionVault creates a new Vault OAuth Provider connection feature
func NewProjectConnectionVault(projectID string, options ConnectionProviderVaultOptions) *ConnectionProviderVault {
	return &ConnectionProviderVault{
		projectID: projectID,
		Options:   options,
	}
}

// ID returns the ID of this project feature
func (f *ConnectionProviderVault) ID() string {
	return f.id
}

// SetID sets the ID of this project feature
func (f *ConnectionProviderVault) SetID(value string) {
	f.id = value
}

// Type represents the type of this project feature as a string
func (f *ConnectionProviderVault) Type() string {
	return "OAuthProvider"
}

// ProjectID represents the ID of the project the project feature is assigned to.
func (f *ConnectionProviderVault) ProjectID() string {
	return f.projectID
}

// SetProjectID sets the ID of the project the project feature is assigned to.
func (f *ConnectionProviderVault) SetProjectID(value string) {
	f.projectID = value
}

// Properties returns all properties for the Vault OAuth Provider project feature
func (f *ConnectionProviderVault) Properties() *Properties {
	return NewProperties(
		NewProperty("auth-method", string(f.Options.AuthMethod)),
		NewProperty("displayName", string(f.Options.DisplayName)),
		NewProperty("endpoint", string(f.Options.Endpoint)),
		NewProperty("fail-on-error", fmt.Sprintf("%t", f.Options.FailOnError)),
		NewProperty("namespace", string(f.Options.Namespace)),
		NewProperty("providerType", string(ConnectionProviderTypeVault)),
		NewProperty("role-id", string(f.Options.RoleID)),
		NewProperty("secure:secret-id", string(f.Options.SecretID)),
		NewProperty("url", string(f.Options.URL)),
		NewProperty("vault-namespace", string(f.Options.VaultNamespace)),
	)
}

func loadConnectionProviderVault(projectID string, feature projectFeatureJSON) (ProjectFeature, error) {
	settings := &ConnectionProviderVault{
		id:        feature.ID,
		projectID: projectID,
		Options:   ConnectionProviderVaultOptions{},
	}

	// stringProperties := []string{"displayName", "endpoint", "namespace", "role-id", "url", "vault-namespace"}
	// for _, property := range stringProperties {
	// 	if encodedValue, ok := feature.Properties.GetOk("displayName"); ok {
	// 		settings.Options.DisplayName = encodedValue
	// 	}
	// }

	if encodedValue, ok := feature.Properties.GetOk("displayName"); ok {
		settings.Options.DisplayName = encodedValue
	}

	if encodedValue, ok := feature.Properties.GetOk("endpoint"); ok {
		settings.Options.Endpoint = encodedValue
	}

	if encodedValue, ok := feature.Properties.GetOk("namespace"); ok {
		settings.Options.Namespace = encodedValue
	}

	if encodedValue, ok := feature.Properties.GetOk("role-id"); ok {
		settings.Options.RoleID = encodedValue
	}

	if encodedValue, ok := feature.Properties.GetOk("url"); ok {
		settings.Options.URL = encodedValue
	}

	if encodedValue, ok := feature.Properties.GetOk("vault-namespace"); ok {
		settings.Options.VaultNamespace = encodedValue
	}

	if encodedValue, ok := feature.Properties.GetOk("fail-on-error"); ok {
		v, err := strconv.ParseBool(encodedValue)
		if err != nil {
			return nil, err
		}

		settings.Options.FailOnError = v
	}

	if encodedValue, ok := feature.Properties.GetOk("auth-method"); ok {
		settings.Options.AuthMethod = ConnectionProviderVaultAuthMethod(encodedValue)
	}

	if encodedValue, ok := feature.Properties.GetOk("providerType"); ok {
		settings.Options.ProviderType = ConnectionProviderType(encodedValue)
	}

	return settings, nil
}
