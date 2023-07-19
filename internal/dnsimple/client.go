package dnsimple

import (
	"context"
	"errors"
	"fmt"
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"os"
	"strconv"
)

type Client struct {
	token          string
	domain         string
	dnSimpleClient *dnsimple.Client
}

func NewClient(token string, domain string) *Client {
	tokenClient := dnsimple.StaticTokenHTTPClient(context.Background(), token)
	client := dnsimple.NewClient(tokenClient)
	if os.Getenv("ENVIRONMENT") != "PRODUCTION" {
		client.BaseURL = "https://api.sandbox.dnsimple.com"
	}

	return &Client{dnSimpleClient: client, domain: domain}
}

func (c *Client) UpdateIP(ip string, record string, ttl int) error {
	r, err := c.dnSimpleClient.Identity.Whoami(context.Background())
	if err != nil {
		return err
	}
	if r.Data.Account == nil {
		return errors.New("no account found. (Are you sure you didn't use a User token?)")
	}

	accountID := strconv.FormatInt(r.Data.Account.ID, 10)
	recType := "A"

	if record == "@" {
		record = ""
	}

	records, err := c.dnSimpleClient.Zones.ListRecords(context.Background(), accountID, c.domain, &dnsimple.ZoneRecordListOptions{
		Name: &record,
		Type: &recType,
	})

	if err != nil {
		fmt.Println("records error: ")
		return err
	}

	if len(records.Data) == 1 {
		_, err = c.dnSimpleClient.Zones.UpdateRecord(context.Background(), accountID, c.domain, records.Data[0].ID, dnsimple.ZoneRecordAttributes{
			Content: ip,
			TTL:     ttl,
		})
	}
	return err
}
