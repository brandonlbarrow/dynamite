package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/brandonlbarrow/dynamite/internal/ipify"
	"github.com/brandonlbarrow/dynamite/internal/route53"

	_ "github.com/joho/godotenv/autoload"
)

var (
	ipifyURL   = "https://api.ipify.org"
	recordName = "dynamitedemo.brandonbarrow.io"
	ttl        = 30
)

func main() {
	ip := getIpAddress(ipifyURL)
	fmt.Printf("ip address is %s\n\n", ip)
	resp, err := run(ip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("status is: %s\n", resp.ChangeInfo.Status)
	fmt.Printf("submitted at: %s\n", resp.ChangeInfo.SubmittedAt)

}

func getIpAddress(url string) string {
	client, err := ipify.NewClient(url)
	if err != nil {
		panic(err)
	}
	client = client.WithHttpClient(&http.Client{})
	addr, err := client.GetIPAddress()
	if err != nil {
		panic(err)
	}
	return string(addr)

}

func run(ip string) (*route53.RecordSetResponse, error) {
	ctx := context.Background()
	cfg, err := route53.NewConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := route53.NewClient(route53.WithCredentials(cfg))
	zones, err := client.ListZones(ctx)
	if err != nil {
		return nil, err
	}
	for _, z := range zones.HostedZones {
		if z.Id == nil {
			continue
		}
		req := route53.NewRecordSetRequest(recordName, ip, *z.Id, int64(ttl))
		if len(zones.HostedZones) == 1 {
			return client.UpsertRecordSet(ctx, req)
		}
	}
	return nil, nil
}
