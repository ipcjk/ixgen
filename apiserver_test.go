package main

import (
	"github.com/ipcjk/ixgen/libapiserver"
	"testing"
	"net/http"
	"bytes"
	"io"
	"bufio"
	"strings"
)

func init() {
	Apiserver := libapiserver.NewAPIServer("localhost:58412", "./cache", "./templates")
	Apiserver.RunAPIServer()
}

func TestApiServer(t *testing.T) {
	_, err := http.Get("http://localhost:58412/api")
	if err != nil {
		t.Error("Cant connect to api service on localhost")
	}
}

func TestPostOnApiServer(t *testing.T) {
	var peering string = `    [DE-CIX Frankfurt/Main]
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

	req, err := http.NewRequest("POST", "http://localhost:58412/ixgen/brocade/netiron/196922", newBuffer)
	if err != nil {
		t.Error(err)
	}

	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("User-Agent", "ixgen/golang")

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
