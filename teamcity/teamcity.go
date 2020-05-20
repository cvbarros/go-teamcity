package teamcity

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/dghubble/sling"

	"github.com/motemen/go-loghttp"
	// Enable HTTP log tracing
	_ "github.com/motemen/go-loghttp/global"
)

type Auth interface{}
type basicAuth struct {
	username, password string
}
type tokenAuth struct {
	token string
}

func BasicAuth(username, password string) Auth {
	return basicAuth{username, password}
}

func TokenAuth(token string) Auth {
	return tokenAuth{token}
}

//DebugRequests toggle to enable tracing requests to stdout
var DebugRequests = false

//DebugResponses toggle to enable tracing responses to stdout
var DebugResponses = false

func init() {
	loghttp.DefaultTransport.LogRequest = func(resp *http.Request) {
		if DebugRequests {
			debug(httputil.DumpRequest(resp, true))
		}
	}
	loghttp.DefaultTransport.LogResponse = func(resp *http.Response) {
		if DebugResponses {
			debug(httputil.DumpResponse(resp, true))
		}
	}
}

//Client represents the base for connecting to TeamCity
type Client struct {
	address string
	baseURI string

	HTTPClient   *http.Client
	RetryTimeout time.Duration

	commonBase *sling.Sling

	Projects        *ProjectService
	BuildTypes      *BuildTypeService
	Server          *ServerService
	VcsRoots        *VcsRootService
	Groups          *GroupService
	RoleAssignments *RoleAssignmentService
}

func NewClient(auth Auth, httpClient *http.Client) (*Client, error) {
	address := os.Getenv("TEAMCITY_ADDR")
	if address == "" {
		return nil, fmt.Errorf("TEAMCITY_ADDR environment variable not set, specify address explicit by setting the variable or using NewWithAddress constructor")
	}

	return newClientInstance(auth, address, httpClient)
}

func NewClientWithAddress(auth Auth, address string, httpClient *http.Client) (*Client, error) {
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}
	return newClientInstance(auth, address, httpClient)
}

func newClientInstance(auth Auth, address string, httpClient *http.Client) (*Client, error) {

	sharedClient := sling.New().
		Set("Accept", "application/json").
		Set("Origin", address)

	switch a := auth.(type) {
	case tokenAuth:
		sharedClient = sharedClient.
			Base(address+"/app/rest/").
			Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
	case basicAuth:
		sharedClient = sharedClient.
			Base(address+"/httpAuth/app/rest/").
			SetBasicAuth(a.username, a.password)

	default:
		return nil, errors.New("unsupported authentication")
	}

	return &Client{
		address:         address,
		HTTPClient:      httpClient,
		commonBase:      sharedClient,
		Projects:        newProjectService(sharedClient.New(), httpClient),
		BuildTypes:      newBuildTypeService(sharedClient.New(), httpClient),
		Server:          newServerService(sharedClient.New()),
		VcsRoots:        newVcsRootService(sharedClient.New(), httpClient),
		Groups:          newGroupService(sharedClient.New(), httpClient),
		RoleAssignments: newRoleAssignmentService(sharedClient.New(), httpClient),
	}, nil
}

// New creates a new client for server address specified at TEAMCITY_ADDR environment variable
// Deprecated
func New(userName, password string, httpClient *http.Client) (*Client, error) {
	return NewClient(BasicAuth(userName, password), httpClient)
}

// NewWithAddress creates a new client by using the explicit server address from the parameter
// Deprecated
func NewWithAddress(userName, password, address string, httpClient *http.Client) (*Client, error) {
	return NewClientWithAddress(BasicAuth(userName, password), address, httpClient)
}

//AgentRequirementService returns a service to manage agent requirements for a build configuration with given id
func (c *Client) AgentRequirementService(id string) *AgentRequirementService {
	return newAgentRequirementService(id, c.HTTPClient, c.commonBase.New())
}

//BuildFeatureService returns a service to manage agent requirements for a build configuration with given id
func (c *Client) BuildFeatureService(id string) *BuildFeatureService {
	return newBuildFeatureService(id, c.HTTPClient, c.commonBase.New())
}

// ProjectFeatureService returns a service to manage project features for a project with given id
func (c *Client) ProjectFeatureService(id string) *ProjectFeatureService {
	return newProjectFeatureService(id, c.HTTPClient, c.commonBase.New())
}

//DependencyService returns a service to manage snapshot and artifact dependencies for a build configuration with given id
func (c *Client) DependencyService(id string) *DependencyService {
	return NewDependencyService(id, c.HTTPClient, c.commonBase.New())
}

//BuildTemplateService returns a service to manage template associations for a build configuration with given id
func (c *Client) BuildTemplateService(id string) *BuildTemplateService {
	return NewBuildTemplateService(id, c.HTTPClient, c.commonBase.New())
}

//TriggerService returns a service to manage build triggers for a build configuration with given id
func (c *Client) TriggerService(buildTypeID string) *TriggerService {
	return newTriggerService(buildTypeID, c.HTTPClient, c.commonBase.New())
}

// Validate tests if the client is properly configured and can be used
func (c *Client) Validate() (bool, error) {
	response, err := c.commonBase.Get("server").ReceiveSuccess(nil)

	if err != nil {
		return false, err
	}

	if response.StatusCode != 200 && response.StatusCode != 403 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("API error %s: %s", response.Status, body)
	}

	return true, nil
}
