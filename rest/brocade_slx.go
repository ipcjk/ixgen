package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

type Brocade_SLX struct {
	apiurl   string
	username string
	password string
}

func NewBrocadeSLX(apiurl, username, password string) Brocade_SLX {
	return Brocade_SLX{apiurl: apiurl, username: username, password: password}
}

func (b *Brocade_SLX) postAPI(uri string, i interface{}) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", b.apiurl+uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("User-Agent", "ixgen/golang")

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
