package main

import (
	"github.com/ipcjk/ixgen/apiserverlib"
	"github.com/ipcjk/ixgen/peeringdb"
	"net/http"
	"testing"
)

func init() {
	Apiserver := apiserverlib.NewAPIServer("localhost:58412", "./cache", "./templates")
	Apiserver.RunAPIServer()
}

func TestGetIX(t *testing.T) {
	peerDB := peeringdb.Peeringdb("https://www.peeringdb.com/api")
	ix := peerDB.SearchIXByIxName("DE-CIX Frankfurt")

	if ix.Data[0].Name != "DE-CIX Frankfurt" {
		t.Error("Cant find the DE-CIX Frankfurt, something wrong the data-set!")
	}
}

func TestRechabilityPeeringDB(t *testing.T) {
	_, err := http.Get("https://www.peeringdb.com/api")
	if err != nil {
		t.Error("Cant connect to peeringdb.com")
	}
}

func TestGetASN(t *testing.T) {
	peerDB := peeringdb.Peeringdb("https://www.peeringdb.com/api")
	peer := peerDB.GetNetworkByAsN(196922)

	if peer.Data[0].Asn != 196922 {
		t.Error("Cant find my home asn, something is wrong!")
	}

	if peer.Data[0].Name != "Hofmeir Media GmbH" {
		t.Error("Cant find my home network name, something is wrong!")
	}

	peer = peerDB.GetNetworkByAsN(2914)
	if peer.Data[0].Asn != 2914 {
		t.Error("Cant find NTT in the peering db?")
	}

}

func TestRunAPIserver(t *testing.T) {
	_, err := http.Get("http://localhost:58412/api")
	if err != nil {
		t.Error("Cant connect to api service on localhost")
	}
}

func TestQueryAPIserver(t *testing.T) {
	peerDB := peeringdb.Peeringdb("http://localhost:58412/api")
	ix := peerDB.SearchIXByIxName("INXS by ecix")

	if ix.Data[0].Name != "INXS by ecix" {
		t.Error("Cant find the INXS, something wrong the data-set!")
	}
}

func TestGetASNLocalApi(t *testing.T) {
	peerDB := peeringdb.Peeringdb("http://localhost:58412/api")
	peer := peerDB.GetNetworkByAsN(196922)

	if peer.Data[0].Asn != 196922 {
		t.Error("Cant find my home asn, something is wrong!")
	}

	if peer.Data[0].Name != "Hofmeir Media GmbH" {
		t.Error("Cant find my home network name, something is wrong!")
	}

	peer = peerDB.GetNetworkByAsN(2914)
	if peer.Data[0].Asn != 2914 {
		t.Error("Cant find NTT in the peering db?")
	}
}

func BenchmarkAPIserver(b *testing.B) {
	peerDB := peeringdb.Peeringdb("http://localhost:58412/api")
	for i := 0; i < b.N; i++ {
		peerDB.GetNetworkByAsN(196922)
	}
}

func TestGetPeersOnIX(t *testing.T) {
	peerDB := peeringdb.Peeringdb("http://localhost:58412/api")
	myPeers := peerDB.GetPeersOnIXByIxName("DE-CIX Frankfurt/Main")

	if len(myPeers.Data) < 700 {
		t.Error("DE-CIX Frankfurt/Main too small")
	}
}
