package ipify

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	address *url.URL
	httpClient
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func NewClient(address string) (*Client, error) {
	parsedURL, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	return &Client{address: parsedURL}, nil
}

func (c *Client) WithHttpClient(client httpClient) *Client {
	c.httpClient = client
	return c
}

func (c *Client) GetIPAddress() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.address.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("empty response")
	}
	return ioutil.ReadAll(resp.Body)
}
