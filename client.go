package teamcity

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

//Client represents the base for connecting to TeamCity
type Client struct {
	userName, password, address string
	baseURI                     string

	HTTPClient   *http.Client
	RetryTimeout time.Duration
}

func NewClient(userName, password string) *Client {
	address := os.Getenv("TEAMCITY_HOST")
	if address == "" {
		address = "http://192.168.99.100:8112"
	}

	return &Client{
		userName:     userName,
		password:     password,
		address:      address,
		HTTPClient:   http.DefaultClient,
		RetryTimeout: time.Duration(60 * time.Second),
	}
}

// SetAddress changes the base address for the Client
func (c *Client) SetAddress(address string) {
	c.address = address
}

// Validate tests if the client is properly configured and can be used
func (c *Client) Validate() (bool, error) {
	req, err := c.createRequest("GET", "server", nil)
	if err != nil {
		return false, err
	}

	var response *http.Response

	response, err = c.HTTPClient.Do(req)

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

	defer response.Body.Close()

	return true, nil
}

func (c *Client) uriForResource(resource string) (string, error) {
	apiBase, err := url.Parse(c.address + "/httpAuth/app/rest/" + resource)
	if err != nil {
		return "", err
	}

	return apiBase.String(), nil
}

func (c *Client) createRequest(method string, resource string, reqbody interface{}) (*http.Request, error) {
	var bodyReader io.Reader

	uri, err := c.uriForResource(resource)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, uri, bodyReader)

	if err != nil {
		return nil, err
	}
	// TODO: Change this later to setup text/plain requests
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", basicAuth(c.userName, c.password))

	if bodyReader != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
