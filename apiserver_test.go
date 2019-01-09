package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/ipcjk/ixgen/ixtypes"
	"github.com/ipcjk/ixgen/libapiserver"
	"io"
	"net/http"
	"strings"
	"testing"
)

var Apiserver *libapiserver.Apiserver

func init() {
	Apiserver = libapiserver.NewAPIServer("localhost:0", "./cache", "./templates", "./configuration")
	Apiserver.RunAPIServer()
}

func TestApiServer(t *testing.T) {
	_, err := http.Get("http://" + Apiserver.AddrPort + "/api")
	if err != nil {
		t.Error("Cant connect to api service on localhost")
	}
}

func TestPostOnApiServer(t *testing.T) {
	var peering = `[DE-CIX Frankfurt]
    [peers]
    714 ipv6=1 ipv4=1
    196922
    [options]
    peer_group=decix-peer
    peer_group6=decix-peer6`

	var newBuffer = bytes.NewBuffer([]byte(peering))
	var lineNum int
	var testCases int

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://"+Apiserver.AddrPort+"/ixgen/brocade/netiron/196922/1.2", newBuffer)
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Do(req)

	if err != nil {
		t.Errorf("HTTP request to apiserver not successful: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Wrong statuscode from apiserver: %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		lineNum++
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("Problems reading from apiserver: %s", err)
		}

		if strings.HasPrefix(line, "neighbor") {
			testCases++
		} else if strings.HasPrefix(line, "address-family") {
			testCases++
		}

		if strings.Contains(line, "196922") {
			t.Error("Saw 196922 but shall not because exluding on post-URL.")
		}

	}

	if lineNum < 12 {
		t.Error("apiserver did not return enough lines for input code")
	}

	if testCases < 8 {
		t.Errorf("Not enough stringSearch cases work. Only %d matched.", testCases)
	}

}

func TestPostJsonOnApiServer(t *testing.T) {
	var peering = `[{"additionalconfig":null,"ixname":"DE-CIX Frankfurt",
	"options":{"DE-CIX Frankfurt":{"wildcard":"0"}},
	"peeringgroups":{},"peers_configured":{"DE-CIX Frankfurt":{"714":[{"active":
	true,"asn":"714","group":"","group6":"","groupenabled":true,"group6_enabled":true,
	"infoprefixes4":0,"infoprefixes6":0,"ipv4addr":"","ipv6addr":"","ipv4enabled":true,
	"ipv6enabled":false,"irrasset":"AS-APPLE","isrs":false,"isrsper":false,"localpreference":0,
	"prefixfilter":true}]}}}]`

	var newBuffer = bytes.NewBuffer([]byte(peering))
	var ixs ixtypes.IXs

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://"+Apiserver.AddrPort+"/ixgen/native/json", newBuffer)
	if err != nil {
		t.Error(err)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		t.Errorf("HTTP request to apiserver not successful: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Wrong statuscode from apiserver: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&ixs)
	if err != nil {
		t.Errorf("Cant decode apiserver results: %s", err)
	}

	if ixs[0].IxName != "DE-CIX Frankfurt" {
		t.Error("Wrong or no IX in result set")
	}

	if len(ixs[0].PeersINI) != 1 {
		t.Error("Wrong numbers of Peer in the JSON-configuration found")
	}

}
