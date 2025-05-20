package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/ipcjk/ixgen/ixtypes"
	"github.com/ipcjk/ixgen/peergen"
	"html/template"
	"io"
	"log"
	"net"
	"strings"
	"testing"
)

func TestExtremeNetworksIX(t *testing.T) {
	var p = peergen.NewPeerGen("extreme/netiron", "./templates", "./configuration")
	var Ix ixtypes.Ix
	var buffer bytes.Buffer
	var err error

	Ix.PeersReady = []ixtypes.ExchangePeer{
		{
			ASN:                 "196922",
			Active:              true,
			Ipv4Enabled:         true,
			Ipv6Enabled:         true,
			PrefixFilterEnabled: false,
			GroupEnabled:        true,
			Group6Enabled:       true,
			IsRs:                false, IsRsPeer: true,
			Ipv4Addr:        net.ParseIP("127.0.0.1"),
			Ipv6Addr:        net.ParseIP("3ffe:ffff::/128"),
			InfoPrefixes6:   10,
			InfoPrefixes4:   100,
			LocalPreference: 90,
			IrrAsSet:        "AS-196922",
			Group:           "decix-peer",
			Group6:          "decix-peer6",
		},
		{
			ASN:                 "3356",
			Active:              true,
			Ipv4Enabled:         true,
			Ipv6Enabled:         true,
			PrefixFilterEnabled: false,
			GroupEnabled:        true,
			Group6Enabled:       true,
			IsRs:                false, IsRsPeer: true,
			Ipv4Addr:        net.ParseIP("127.3.3.56"),
			Ipv6Addr:        net.ParseIP("3ffe:ffff:3356::"),
			InfoPrefixes6:   10,
			InfoPrefixes4:   100,
			LocalPreference: 90,
			IrrAsSet:        "AS-3356",
			Group:           "decix-peer",
			Group6:          "decix-peer6",
			PrefixFilters: ixtypes.PrefixFilters{
				Name: "3356peer",
				PrefixRules: []ixtypes.PrefixRule{
					{Prefix: "178.248.240.4/24", Exact: true},
				},
			},
			PrefixFilters6: ixtypes.PrefixFilters{
				Name: "3356peer",
				PrefixRules: []ixtypes.PrefixRule{
					{Prefix: "2a02:1308::/32", GreaterEqual: 32, LessEqual: 48},
				},
			},
		},
	}

	// validate ipv4/ipv6 like peergen would to
	for x := range Ix.PeersReady {
		// Valid Ipv6? Then, don't load from peerDB
		if Ix.PeersReady[x].Ipv6Addr.To16() == nil {
			Ix.PeersReady[x].Ipv6Enabled = false
		}
		// Valid Ipv4? Then, don't load from peerDB
		if Ix.PeersReady[x].Ipv4Addr.To4() == nil {
			Ix.PeersReady[x].Ipv4Enabled = false
		}
	}

	for _, peer := range Ix.PeersReady {
		if peer.Ipv4Enabled && peer.Ipv4Addr.To4() == nil {
			t.Errorf("IPv4 is enabled for %s but IPv4 is not set", peer.ASN)
		}
		if peer.Ipv6Enabled && peer.Ipv6Addr.To16() == nil {
			t.Errorf("IPv6 is enabled for %s but IPv6 is not set", peer.ASN)
		}
	}
	writer := bufio.NewWriter(&buffer)
	err = p.GenerateIXConfiguration(Ix, writer)
	if err != nil {
		log.Fatal(err)
	}

	err = writer.Flush()
	if err != nil {
		log.Fatal("Cant flush generated configuration into buffer")
	}

	var countLines, countNeighbor int
	var foundSample bool
	reader := bufio.NewReader(&buffer)
	for {
		var line string
		line, err = reader.ReadString('\n')

		if strings.HasPrefix(line, "neighbor") {
			countNeighbor++
		}
		if strings.HasPrefix(line, "neighbor 127.0.0.1 remote-as 196922") {
			foundSample = true
		}
		countLines++

		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error("Error reading from template buffer")
		}
	}

	if !foundSample {
		t.Error("Did not find bgp neighbor sample command in template buffer")
	}
	if countLines < 16 {
		t.Error("Template too short or broken, not enough output lines for netiron/extreme")
	}
	if countNeighbor < 8 {
		t.Error("Template too short or broken, not enough bgp neighbor commands")
	}

	filterWriter := bufio.NewWriter(new(bytes.Buffer))
	err = p.GeneratePrefixFilter(Ix, filterWriter)
	if err != nil {
		t.Error(err)
	}

	err = filterWriter.Flush()
	if err != nil {
		t.Error("Cant flush generated configuration into buffer")
	}
}

func TestExtremeNetironPrefixFilterTemplate(t *testing.T) {
	var p = peergen.NewPeerGen("extreme/netiron", "./templates", "./configuration")
	var Ix ixtypes.Ix

	Ix.PeersReady = []ixtypes.ExchangePeer{
		{
			ASN:                 "3356",
			Active:              true,
			Ipv4Enabled:         true,
			Ipv6Enabled:         true,
			IrrAsSet:            "AS-3356",
			PrefixFilterEnabled: true,
			PrefixFilters: ixtypes.PrefixFilters{
				Name: "3356peer-4",
				PrefixRules: []ixtypes.PrefixRule{
					{Prefix: "178.248.240.0/21", Exact: true},
					{Prefix: "178.248.241.0/24", Exact: true},
				},
			},
			PrefixFilters6: ixtypes.PrefixFilters{
				Name: "3356peer-6",
				PrefixRules: []ixtypes.PrefixRule{
					{Prefix: "2a02:1308::/32", GreaterEqual: 32, LessEqual: 48},
					{Prefix: "2a02:1308::/48", Exact: true},
				},
			},
		},
	}

	buffer := new(bytes.Buffer)
	err := p.GeneratePrefixFilter(Ix, buffer)
	if err != nil {
		t.Error(err)
	}

	if !strings.ContainsAny(buffer.String(), "2a02:1308::/48") {
		t.Error("Cant find my home prefix in ipv6-prefixlist")
	}

	if !strings.ContainsAny(buffer.String(), "178.248.240.0/24") {
		t.Error("Cant find my home prefix in ip prefixlist")
	}

}

func TestAllTemplates(t *testing.T) {
	templateDir := "./templates/"
	supportedTemplate := []string{
		"extreme/netiron/router.tt",
		"extreme/netiron/prefix.tt",
		"extreme/netiron/prefix6.tt",
		"extreme/slx/prefix6.tt",
		"extreme/slx/router.tt",
		"extreme/slx/header.tt",
		"extreme/slx/prefix.tt",
		"extreme/vdx/header.tt",
		"extreme/vdx/router.tt",
		"cisco/ios/prefix.tt",
		"cisco/ios/prefix6.tt",
		"cisco/ios/router.tt",
		"juniper/set/header.tt",
		"juniper/set/prefix.tt",
		"juniper/set/router.tt",
	}

	for _, v := range supportedTemplate {
		_, err := template.New("test").ParseFiles(templateDir + v)
		if err != nil {
			t.Errorf("broken template: %s, %s", v, err)
		} else {
			t.Logf("tt %s ok", v)
		}
	}
}

func TestIXConfigFromJson(t *testing.T) {
	var testJSON = `{"additionalconfig":null,"ixname":"DE-CIX Frankfurt","options":{},"peeringgroups":{},"peers_configured":{"DE-CIX Frankfurt":{"196922":[{"active":true,"asn":"196922","group":"","group6":"","groupenabled":true,"group6_enabled":true,"infoprefixes4":0,"infoprefixes6":0,"ipv4addr":"","ipv6addr":"","ipv4enabled":true,"ipv6enabled":true,"irrasset":"","isrs":false,"isrsper":false,"localpreference":0,"prefixfilter":false}]}},"peersready":[{"active":true,"asn":"196922","group":"","group6":"","groupenabled":false,"group6_enabled":false,"infoprefixes4":64,"infoprefixes6":10,"ipv4addr":"80.81.194.25","ipv6addr":"2001:7f8::3:13a:0:1","ipv4enabled":true,"ipv6enabled":true,"irrasset":"AS-HOFMEIR","isrs":false,"isrsper":false,"localpreference":0,"prefixfilter":false}],"routeserverready":null}`
	var p = peergen.NewPeerGen("extreme/netiron", "./templates", "./configuration")
	var buffer bytes.Buffer

	ix := ixtypes.Ix{}

	if err := json.Unmarshal([]byte(testJSON), &ix); err != nil {
		t.Errorf("error decoding JSON into format, some code has changed? Error %s", err)
	}

	if ix.IxName != "DE-CIX Frankfurt" {
		t.Error("IX-Name has changed, not expected")
	}

	writer := bufio.NewWriter(&buffer)
	p.GenerateIXConfiguration(ix, writer)

	err := writer.Flush()
	if err != nil {
		log.Fatal("Cant flush generated configuration into buffer")
	}

	var countLines, countNeighbor int
	var foundSample bool
	reader := bufio.NewReader(&buffer)
	for {
		line, err := reader.ReadString('\n')

		if strings.HasPrefix(line, "neighbor") {
			countNeighbor++
		}
		if strings.HasPrefix(line, "neighbor 80.81.194.25 remote-as 196922") {
			foundSample = true
		}
		countLines++

		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error("Error reading from template buffer")
		}
	}

	if !foundSample {
		t.Error("Did not find any bgp neighbor sample command in template buffer")
	}

	if countLines < 8 {
		t.Error("Template too short or broken, not enough output lines for netiron/extreme")
	}
	if countNeighbor < 2 {
		t.Error("Template too short or broken, not enough bgp neighbor commands")
	}

}
