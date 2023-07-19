package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/brandonlbarrow/dynamite/internal/dnsimple"
	"github.com/brandonlbarrow/dynamite/internal/ipify"
	"github.com/brandonlbarrow/dynamite/internal/namecheap"
	"github.com/brandonlbarrow/dynamite/internal/route53"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Registrar interface {
	UpdateIP(ip string, record string, ttl int) error
}

var (
	ipifyURL   = flag.String("api", "https://api.ipify.org", "URL to retrieve external API. Optional, defaults to https://api.ipify.org")
	registrar  = flag.String("registrar", "route53", "DNS registrar.  One of: route53, namecheap, dnsimple.  Optional, Defaults to route53")
	recordName = flag.String("record", "", "DNS record to update")
	ttl        = 30
)

func main() {
	flag.Parse()
	selectedRegistrar := validateInput()
	ip := getIpAddress(*ipifyURL)
	fmt.Printf("ip address is %s\n\n", ip)
	err := run(ip, *recordName, ttl, selectedRegistrar)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated %s to %s", *recordName, ip)
}

func validateInput() Registrar {
	if recordName == nil || *recordName == "" {
		panic("--record option required")
	}

	switch *registrar {
	case "route53":
		ctx := context.Background()
		cfg, err := route53.NewConfig(ctx)
		if err != nil {
			return nil
		}
		return route53.NewClient(route53.WithCredentials(cfg))

	case "namecheap":
		domain := os.Getenv("DOMAIN")
		password := os.Getenv("DDNS_PASSWORD")
		if domain == "" || password == "" {
			panic("DOMAIN and DDNS_PASSWORD required in environment")
		}
		netClient := &http.Client{Timeout: 10 * time.Second}
		return namecheap.NewClient(domain, password, netClient)
	case "dnsimple":
		token := os.Getenv("TOKEN")
		domain := os.Getenv("DOMAIN")
		if token == "" || domain == "" {
			panic("TOKEN and DOMAIN required in environment for dnsimple")
		}
		return dnsimple.NewClient(token, domain)
	default:
		panic("Valid options for --registrar are: route53, namecheap")
	}
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

func run(ip string, record string, ttl int, selectedRegistrar Registrar) error {
	return selectedRegistrar.UpdateIP(ip, record, ttl)
}
