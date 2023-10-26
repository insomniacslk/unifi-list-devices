package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/unpoller/unifi"
)

var (
	flagUsername = flag.String("u", "", "Unifi controller username")
	flagPassword = flag.String("p", "", "Unifi controller password")
	flagURL      = flag.String("U", "http://127.0.0.1:8443", "Unifi controller URL")
	flagSiteName = flag.String("s", "default", "Site name")
)

func main() {
	flag.Parse()
	c := unifi.Config{
		User:     *flagUsername,
		Pass:     *flagPassword,
		URL:      *flagURL,
		ErrorLog: log.Printf,
		DebugLog: nil,
	}
	uni, err := unifi.NewUnifi(&c)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	sites, err := uni.GetSites()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	siteIdx := -1
	for idx, site := range sites {
		if site.Name == *flagSiteName {
			siteIdx = idx
		}
	}
	if siteIdx == -1 {
		log.Fatalf("Site '%s' not found", *flagSiteName)
	}
	clients, err := uni.GetClients([]*unifi.Site{sites[siteIdx]})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	var output string
	idx := 0
	for _, client := range clients {
		fmt.Println(idx+1, client.ID, client.Hostname, client.IP, client.Name, client.Mac, client.LastSeen)
		output += fmt.Sprintf("%s\n", client.IP)
		idx++
	}
}
