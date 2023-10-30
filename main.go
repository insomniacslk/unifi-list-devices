package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/rgeoghegan/tabulate"
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
	rows := make([][]string, 0)
	header := []string{"#", "ID", "IP", "Hostname", "Name", "MAC", "Last seen"}
	for _, client := range clients {
		lastSeen := time.Unix(int64(client.LastSeen.Val), 0)
		rows = append(rows, []string{
			strconv.FormatInt(int64(idx+1), 10),
			client.ID,
			client.IP,
			client.Hostname,
			client.Name,
			client.Mac,
			lastSeen.String(),
		})
		//fmt.Printf("% 2d) %s\t%s %s %s %s %v\n", idx+1, client.ID, client.IP, client.Hostname, client.Name, client.Mac, lastSeen)
		output += fmt.Sprintf("%s\n", client.IP)
		idx++
	}
	table, err := tabulate.Tabulate(rows, &tabulate.Layout{Headers: header, Format: tabulate.SimpleFormat})
	if err != nil {
		log.Fatalf("Failed to tabulate text: %v", err)
	}
	fmt.Println(table)
}
