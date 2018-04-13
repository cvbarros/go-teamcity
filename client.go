package teamcity

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	userName, password, address string

	HttpClient   *http.Client
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
		HttpClient:   http.DefaultClient,
		RetryTimeout: time.Duration(60 * time.Second),
	}
}

func (c *Client) SetAddress(address string) {
	c.address = address
}

func (c *Client) Validate() (bool, error) {
	var bodyReader io.Reader

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/httpAuth/app/rest/10.0/server", c.address), bodyReader)
	if err != nil {
		return false, err
	}

	if bodyReader != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	var response *http.Response

	response, err = c.HttpClient.Do(req)

	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	return true, nil
}
