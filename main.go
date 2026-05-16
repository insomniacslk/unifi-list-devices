package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rgeoghegan/tabulate"
	"github.com/spf13/pflag"
	"github.com/unpoller/unifi"
)

type fieldDef struct {
	header  string
	extract func(idx int, c *unifi.Client) string
}

var fieldDefs = map[string]fieldDef{
	"num":        {"#", func(idx int, _ *unifi.Client) string { return strconv.Itoa(idx + 1) }},
	"id":         {"ID", func(_ int, c *unifi.Client) string { return c.ID }},
	"ip":         {"IP", func(_ int, c *unifi.Client) string { return c.IP }},
	"hostname":   {"Hostname", func(_ int, c *unifi.Client) string { return c.Hostname }},
	"name":       {"Name", func(_ int, c *unifi.Client) string { return c.Name }},
	"mac":        {"MAC", func(_ int, c *unifi.Client) string { return c.Mac }},
	"last-seen":  {"Last seen", func(_ int, c *unifi.Client) string { return time.Unix(int64(c.LastSeen.Val), 0).String() }},
	"uptime":     {"Uptime", func(_ int, c *unifi.Client) string { return (time.Duration(int64(c.Uptime.Val)) * time.Second).String() }},
	"first-seen": {"First seen", func(_ int, c *unifi.Client) string { return time.Unix(int64(c.FirstSeen.Val), 0).String() }},
	"assoc-time": {"Assoc time", func(_ int, c *unifi.Client) string { return time.Unix(int64(c.AssocTime.Val), 0).String() }},
}

var defaultFields = []string{"num", "id", "ip", "hostname", "name", "mac", "last-seen", "uptime"}

var (
	flagUsername = pflag.StringP("username", "u", "", "Unifi controller username")
	flagPassword = pflag.StringP("password", "p", "", "Unifi controller password")
	flagURL      = pflag.StringP("url", "U", "http://127.0.0.1:8443", "Unifi controller URL")
	flagSiteName = pflag.StringP("site", "s", "default", "Site name")
	flagSort     = pflag.String("sort", "name", "Sort by: ip, hostname, name, mac, last-seen, uptime, first-seen, assoc-time")
	flagFields   = pflag.StringSliceP("fields", "f", defaultFields, "Comma-separated fields to display: num, id, ip, hostname, name, mac, last-seen, uptime, first-seen, assoc-time")
)

func main() {
	pflag.Parse()

	defs, err := resolveFields(*flagFields)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

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

	less, err := clientLess(*flagSort, clients)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	sort.Slice(clients, less)

	header := make([]string, 0, len(defs))
	for _, d := range defs {
		header = append(header, d.header)
	}
	rows := make([][]string, 0, len(clients))
	for idx, client := range clients {
		row := make([]string, 0, len(defs))
		for _, d := range defs {
			row = append(row, d.extract(idx, client))
		}
		rows = append(rows, row)
	}
	table, err := tabulate.Tabulate(rows, &tabulate.Layout{Headers: header, Format: tabulate.SimpleFormat})
	if err != nil {
		log.Fatalf("Failed to tabulate text: %v", err)
	}
	fmt.Println(table)
}

func resolveFields(names []string) ([]fieldDef, error) {
	defs := make([]fieldDef, 0, len(names))
	for _, n := range names {
		n = strings.TrimSpace(n)
		def, ok := fieldDefs[n]
		if !ok {
			return nil, fmt.Errorf("invalid field %q", n)
		}
		defs = append(defs, def)
	}
	return defs, nil
}

func clientLess(sortBy string, clients []*unifi.Client) (func(i, j int) bool, error) {
	switch sortBy {
	case "ip":
		return func(i, j int) bool { return clients[i].IP < clients[j].IP }, nil
	case "hostname":
		return func(i, j int) bool { return clients[i].Hostname < clients[j].Hostname }, nil
	case "name":
		return func(i, j int) bool { return clients[i].Name < clients[j].Name }, nil
	case "mac":
		return func(i, j int) bool { return clients[i].Mac < clients[j].Mac }, nil
	case "last-seen":
		return func(i, j int) bool { return clients[i].LastSeen.Val < clients[j].LastSeen.Val }, nil
	case "uptime":
		return func(i, j int) bool { return clients[i].Uptime.Val < clients[j].Uptime.Val }, nil
	case "first-seen":
		return func(i, j int) bool { return clients[i].FirstSeen.Val < clients[j].FirstSeen.Val }, nil
	case "assoc-time":
		return func(i, j int) bool { return clients[i].AssocTime.Val < clients[j].AssocTime.Val }, nil
	default:
		return nil, fmt.Errorf("invalid sort field %q (valid: ip, hostname, name, mac, last-seen, uptime, first-seen, assoc-time)", sortBy)
	}
}
