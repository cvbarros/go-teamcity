package teamcity

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/cenkalti/backoff"
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

func (c *Client) createRequest(method string, resource string, requestBody interface{}) (*http.Request, error) {
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

func (c *Client) doJSONRequest(method, api string, requestBody, out interface{}) error {
	req, err := c.createRequest(method, api, requestBody)
	if err != nil {
		return err
	}

	// Perform the request and retry it if it's not a POST or PUT request
	var resp *http.Response
	if method == "POST" || method == "PUT" {
		resp, err = c.HTTPClient.Do(req)
	} else {
		resp, err = c.doRequestWithRetries(req, c.RetryTimeout)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("API error %s: %s", resp.Status, body)
	}

	// If they don't care about the body, then we don't care to give them one,
	// so bail out because we're done.
	if out == nil {
		// read the response body so http conn can be reused immediately
		io.Copy(ioutil.Discard, resp.Body)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// If we got no body, by default let's just make an empty JSON dict. This
	// saves us some work in other parts of the code.
	if len(body) == 0 {
		body = []byte{'{', '}'}
	}

	return json.Unmarshal(body, &out)
}

// doRequestWithRetries performs an HTTP request repeatedly for maxTime or until
// no error and no acceptable HTTP response code was returned.
func (c *Client) doRequestWithRetries(req *http.Request, maxTime time.Duration) (*http.Response, error) {
	var (
		err  error
		resp *http.Response
		bo   = backoff.NewExponentialBackOff()
		body []byte
	)

	bo.MaxElapsedTime = maxTime

	// Save the body for retries
	if req.Body != nil {
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return resp, err
		}
	}

	operation := func() error {
		if body != nil {
			r := bytes.NewReader(body)
			req.Body = ioutil.NopCloser(r)
		}

		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// 2xx all done
			return nil
		} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// 4xx are not retryable
			return nil
		}

		return fmt.Errorf("Received HTTP status code %d", resp.StatusCode)
	}

	err = backoff.Retry(operation, bo)

	return resp, err
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
