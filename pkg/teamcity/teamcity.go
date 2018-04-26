package teamcity

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/dghubble/sling"

	loghttp "github.com/motemen/go-loghttp"
	// Enable HTTP log tracing
	_ "github.com/motemen/go-loghttp/global"
)

//DebugRequests toggle to enable tracing requests to stdout
var DebugRequests = false

//DebugResponses tottle to enable tracing responses to stdout
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
	userName, password, address string
	baseURI                     string

	HTTPClient   *http.Client
	RetryTimeout time.Duration

	commonBase *sling.Sling

	Projects   *ProjectService
	BuildTypes *BuildTypeService
	Server     *ServerService
	VcsRoots   *VcsRootService
	Parameters *ParameterService
}

// New creates a new client for interating with TeamCity API
func New(userName, password string) *Client {
	address := os.Getenv("TEAMCITY_HOST")
	if address == "" {
		address = "http://192.168.99.100:8112"
	}

	sharedClient := sling.New().Base(address+"/httpAuth/app/rest/").
		SetBasicAuth(userName, password).
		Set("Accept", "application/json")

	return &Client{
		userName:   userName,
		password:   password,
		address:    address,
		HTTPClient: http.DefaultClient,
		commonBase: sharedClient,
		Projects:   newProjectService(sharedClient.New()),
		BuildTypes: newBuildTypeService(sharedClient.New()),
		Server:     newServerService(sharedClient.New()),
		VcsRoots:   newVcsRootService(sharedClient.New()),
		Parameters: newParameterService(sharedClient.New()),
	}
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

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Printf("[ERROR] %s\n\n", err)
	}
}
