package namecheap

import (
	"errors"
	"fmt"
	"net/http"
)

const baseUrl = "https://dynamicdns.park-your-domain.com/update"

type Client struct {
	domain    string
	password  string
	netClient *http.Client
}

func NewClient(domain string, password string, client *http.Client) *Client {
	return &Client{
		domain:    domain,
		password:  password,
		netClient: client,
	}
}

func (c *Client) UpdateIP(ip string, recordName string, ttl int) error {
	host := fmt.Sprintf(baseUrl+"?host=%s&domain=%s&password=%s&ip=%s", recordName, c.domain, c.password, ip)
	resp, err := c.netClient.Get(host)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Error updating record: " + resp.Status)
	}

	return nil
}
