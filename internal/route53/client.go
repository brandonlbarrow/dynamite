package route53

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

const (
	region       = "us-east-1"
	upsertAction = "UPSERT"
	recordTypeA  = "A"
)

type Client struct {
	client *r53.Client
}

type ClientOption func(c *Client)

func NewClient(opts ...ClientOption) *Client {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithCredentials(creds aws.Config) ClientOption {
	return func(c *Client) {
		c.client = r53.NewFromConfig(creds)
	}
}

func NewConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, config.WithRegion(region))
}

type HostedZoneResponse struct {
	*r53.ListHostedZonesOutput
}

func (c *Client) ListZones(ctx context.Context) (*HostedZoneResponse, error) {
	resp, err := c.client.ListHostedZones(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &HostedZoneResponse{resp}, nil
}

type RecordSetRequest struct {
	*r53.ChangeResourceRecordSetsInput
}

func NewRecordSetRequest(recordName, recordValue, hostedZoneID string, ttl int64) *RecordSetRequest {
	r := &RecordSetRequest{
		&r53.ChangeResourceRecordSetsInput{
			HostedZoneId: &hostedZoneID,
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action: upsertAction,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name: &recordName,
							ResourceRecords: []types.ResourceRecord{
								{
									Value: &recordValue,
								},
							},
							TTL:  &ttl,
							Type: recordTypeA,
						},
					},
				},
			},
		},
	}
	return r
}

type RecordSetResponse struct {
	*r53.ChangeResourceRecordSetsOutput
}

func (c *Client) UpsertRecordSet(ctx context.Context, req *RecordSetRequest) (*RecordSetResponse, error) {
	resp, err := c.client.ChangeResourceRecordSets(ctx, req.ChangeResourceRecordSetsInput)
	if err != nil {
		return nil, err
	}
	return &RecordSetResponse{resp}, nil
}

func (c *Client) UpdateIP(ip string, recordName string, ttl int) error {
	ctx := context.Background()
	zones, err := c.ListZones(ctx)
	if err != nil {
		return err
	}

	for _, z := range zones.HostedZones {
		if z.Id == nil {
			continue
		}
		req := NewRecordSetRequest(recordName, ip, *z.Id, int64(ttl))
		if len(zones.HostedZones) == 1 {
			_, err = c.UpsertRecordSet(ctx, req)
			return err
		}
	}

	return errors.New("record not found")
}
