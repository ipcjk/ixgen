package main

import (
	"bufio"
	"bytes"
	"github.com/ipcjk/ixgen/ixtypes"
	"github.com/ipcjk/ixgen/peergen"
	"html/template"
	"io"
	"log"
	"net"
	"strings"
	"testing"
)

func TestBrocadeTemplate(t *testing.T) {
	var p = peergen.NewPeerGen("brocade/netiron", "./templates")
	var Ix ixtypes.Ix
	var buffer bytes.Buffer

	Ix.PeersReady = []ixtypes.ExchangePeer{
		{
			ASN:           "196922",
			Active:        true,
			Ipv4Enabled:   true,
			Ipv6Enabled:   true,
			PrefixFilter:  false,
			GroupEnabled:  true,
			Group6Enabled: true,
			IsRs:          false, IsRsPeer: true,
			Ipv4Addr:        net.ParseIP("127.0.0.1"),
			Ipv6Addr:        net.ParseIP("3ffe:ffff::/32"),
			InfoPrefixes6:   10,
			InfoPrefixes4:   100,
			LocalPreference: 90,
			IrrAsSet:        "AS-196922",
			Group:           "decix-peer",
			Group6:          "decix-peer6",
		},
		{
			ASN:           "3356",
			Active:        true,
			Ipv4Enabled:   true,
			Ipv6Enabled:   true,
			PrefixFilter:  false,
			GroupEnabled:  true,
			Group6Enabled: true,
			IsRs:          false, IsRsPeer: true,
			Ipv4Addr:        net.ParseIP("127.3.3.56"),
			Ipv6Addr:        net.ParseIP("3ffe:ffff:3356::/32"),
			InfoPrefixes6:   10,
			InfoPrefixes4:   100,
			LocalPreference: 90,
			IrrAsSet:        "AS-3356",
			Group:           "decix-peer",
			Group6:          "decix-peer6",
		},
	}

	writer := bufio.NewWriter(&buffer)
	p.GenerateIXConfiguration(Ix, writer)

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

	if foundSample == false {
		t.Error("Did not find bgp neighbor sample command in template buffer")
	}
	if countLines < 16 {
		t.Error("Template too short or broken, not enough output lines for netiron/brocade")
	}
	if countNeighbor < 8 {
		t.Error("Template too short or broken, not enough bgp neighbor commands")
	}

}

func TestAllTemplates(t *testing.T) {
	templateDir := "./templates/"
	supportedTemplate := []string{
		"brocade/netiron/router.tt",
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
