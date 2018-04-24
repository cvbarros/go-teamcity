package teamcity

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/sling"
)

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
