package main

import (
	"github.com/ipcjk/ixgen/libapiserver"
	"github.com/ipcjk/ixgen/peeringdb"
	"net/http"
	"testing"
)

var apiserverDbTest *libapiserver.Apiserver

func init() {
	apiserverDbTest = libapiserver.NewAPIServer("localhost:0", "./cache", "./templates")
	apiserverDbTest.RunAPIServer()
}

func TestGetIX(t *testing.T) {
	peerDB := peeringdb.Peeringdb("https://www.peeringdb.com/api")
	ix, err := peerDB.SearchIXByIxName("DE-CIX Frankfurt")
	if err != nil {
		t.Errorf("Cant search by IxName: %s", err)
	}

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
	peer, err := peerDB.GetNetworkByAsN(196922)

	if err != nil {
		t.Errorf("Cant search by the target ASN: %s", err)
	}

	if peer.Data[0].Asn != 196922 {
		t.Error("Cant find my home asn, something is wrong!")
	}

	if peer.Data[0].Name != "Hofmeir Media GmbH" {
		t.Error("Cant find my home network name, something is wrong!")
	}

	peer, err = peerDB.GetNetworkByAsN(2914)
	if err != nil {
		t.Errorf("Cant search by the target ASN: %s", err)
	}

	if peer.Data[0].Asn != 2914 {
		t.Error("Cant find NTT in the peering db?")
	}

}

func TestRunAPIserver(t *testing.T) {
	_, err := http.Get("http://" + apiserverDbTest.AddrPort + "/api")
	if err != nil {
		t.Error("Cant connect to api service on localhost")
	}
}

func TestQueryAPIserver(t *testing.T) {
	peerDB := peeringdb.Peeringdb("http://" + apiserverDbTest.AddrPort + "/api")
	ix, err := peerDB.SearchIXByIxName("ECIX-MUC / INXS by ecix")

	if err != nil {
		t.Errorf("Cant search by the IxName: %s", err)
	}

	if ix.Data[0].Name != "ECIX-MUC / INXS by ecix" {
		t.Error("Cant find the INXS, something wrong the data-set!")
	}
}

func TestGetASNLocalApi(t *testing.T) {
	peerDB := peeringdb.Peeringdb("http://" + apiserverDbTest.AddrPort + "/api")
	peer, err := peerDB.GetNetworkByAsN(196922)

	if err != nil {
		t.Errorf("Cant search by target ASN: %s", err)
	}

	if peer.Data[0].Asn != 196922 {
		t.Error("Cant find my home asn, something is wrong!")
	}

	if peer.Data[0].Name != "Hofmeir Media GmbH" {
		t.Error("Cant find my home network name, something is wrong!")
	}

	peer, err = peerDB.GetNetworkByAsN(2914)
	if err != nil {
		t.Errorf("Cant search by target ASN: %s", err)
	}

	if peer.Data[0].Asn != 2914 {
		t.Error("Cant find NTT in the peering db?")
	}
}

func BenchmarkAPIserver(b *testing.B) {
	peerDB := peeringdb.Peeringdb("http://" + apiserverDbTest.AddrPort + "/api")
	for i := 0; i < b.N; i++ {
		peerDB.GetNetworkByAsN(196922)
	}
}

func TestGetPeersOnIX(t *testing.T) {
	peerDB := peeringdb.Peeringdb("http://" + apiserverDbTest.AddrPort + "/api")
	myPeers, err := peerDB.GetPeersOnIXByIxName("DE-CIX Frankfurt||Main")

	if err != nil {
		t.Errorf("Cant query the API for peers on IX: %s", err)
	}

	if len(myPeers.Data) < 700 {
		t.Error("DE-CIX Frankfurt||Main too small")
	}
}
