package rest

import (
	"encoding/json"
	"github.com/ipcjk/ixgen/ixtypes"
	"log"
	"net/http"
)

type ExtremeSLX struct {
	apiUrl   string
	username string
	password string
}

func NewExtremeSLX(apiUrl, username, password string) ExtremeSLX {
	return ExtremeSLX{apiUrl: apiUrl, username: username, password: password}
}

func (b *ExtremeSLX) postAPI(uri string, i interface{}) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", b.apiUrl+uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	req.Header.Add("User-Agent", "ixgen/golang")
	req.SetBasicAuth(b.username, b.password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTP request not successful: %s", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("HTTP Api server responded with %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&i)
	if err != nil {
		log.Fatalf("Problems decoding http api output: %s", err)
	}
}

func (b *ExtremeSLX) readConfiguration() {

}

func (b *ExtremeSLX) configureBgpPeers(ix ixtypes.IXs) {
	/* Read ix and post configuration */

}

func (b *ExtremeSLX) unConfigureBgpPeer() {

}

func (b *ExtremeSLX) generatePrefixList() {

}
