package main

import (
	"fmt"
	"net/http"

	"github.com/brandonlbarrow/dynamite/internal/ipify"
)

var ipifyURL = "https://api.ipify.org"

func main() {
	fmt.Println(getIpAddress(ipifyURL))
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
