package main

import (
	"github.com/ipcjk/ixgen/libapiserver"
	"github.com/ipcjk/ixgen/peeringdb"
	"net/http"
	"os"
	"testing"
)

var apiserverDbTest *libapiserver.Apiserver
var apikey string

func TestApiKey(t *testing.T) {
	apikey = os.Getenv("PEERINGDB_APIKEY")
	if apikey == "" {
		t.Fatal("Please set the PEERINGDB_APIKEY environment variable to your apikey")
	}
}

func init() {
	apiserverDbTest = libapiserver.NewAPIServer("localhost:0", "./cache", "./templates", "./configuration")
	apiserverDbTest.RunAPIServer()
}

func TestGetIXByName(t *testing.T) {
	peerDB := peeringdb.Peeringdb("https://www.peeringdb.com/api", apikey)
	ix, err := peerDB.SearchIXByIxName("DE-CIX Frankfurt")
	if err != nil {
		t.Errorf("Cant search by IxName: %s", err)
	} else {
		if ix.Data[0].Name != "DE-CIX Frankfurt" {
			t.Error("Cant find the DE-CIX Frankfurt, something wrong the data-set!")
		}
	}
}

func TestGetIXByID(t *testing.T) {
	peerDB := peeringdb.Peeringdb("https://www.peeringdb.com/api", apikey)
	ix, err := peerDB.SearchIXByIxId("73")
	if err != nil {
		t.Errorf("Cant search by IxID: %s", err)
		return
	}

	if ix.Data[0].Name != "MegaIX Munich" {
		t.Error("Cant find the MegaIX Munich, something wrong inside the data-set!")
	}
}

func TestRechabilityPeeringDB(t *testing.T) {
	_, err := http.Get("https://www.peeringdb.com/api")
	if err != nil {
		t.Error("Cant connect to peeringdb.com")
	}
}

func TestGetASN(t *testing.T) {
	peerDB := peeringdb.Peeringdb("https://www.peeringdb.com/api", apikey)
	peer, err := peerDB.GetNetworkByAsN(196922)

	if err != nil {
		t.Errorf("Cant search by the target ASN: %s", err)
		return
	}

	if peer.Data[0].Asn != 196922 {
		t.Error("Cant find my home asn, something is wrong!")
	}

	if peer.Data[0].Name != "Hofmeir Media GmbH" {
		t.Error("Cant find my home network name, something is wrong! Found ", peer.Data[0].Name)
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
	peerDB := peeringdb.Peeringdb("http://"+apiserverDbTest.AddrPort+"/api", apikey)
	ix, err := peerDB.SearchIXByIxName("DE-CIX Frankfurt")

	if err != nil {
		t.Errorf("Cant search by the IxName: %s", err)
	}

	if ix.Data[0].Name != "DE-CIX Frankfurt" {
		t.Error("Cant find the INXS, something wrong the data-set!")
	}
}

func TestGetASNLocalApi(t *testing.T) {
	peerDB := peeringdb.Peeringdb("http://"+apiserverDbTest.AddrPort+"/api", apikey)
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
	peerDB := peeringdb.Peeringdb("http://"+apiserverDbTest.AddrPort+"/api", apikey)
	for i := 0; i < b.N; i++ {
		peerDB.GetNetworkByAsN(196922)
	}
}

func TestGetPeersOnIX(t *testing.T) {
	peerDB := peeringdb.Peeringdb("http://"+apiserverDbTest.AddrPort+"/api", apikey)
	myPeers, err := peerDB.GetPeersOnIX("DE-CIX Frankfurt", "0", false)

	if err != nil {
		t.Errorf("Cant query the API for peers on IX: %s", err)
	}

	if len(myPeers.Data) < 700 {
		t.Error("DE-CIX Frankfurt||Main too small")
	}

	/* Query by static id */
	myPeers, err = peerDB.GetPeersOnIX("", "73", true)

	if err != nil {
		t.Errorf("Cant query the API for peers on IX: %s", err)
	}

	if len(myPeers.Data) < 30 {
		t.Error("ECIX-MUC too small")
	}

}

func mini(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

func TestSplitASN(t *testing.T) {
	var list []int
	for x := 0; x < 120; x++ {
		list = append(list, x)
	}

	for {
		if len(list) == 0 {
			break
		}
		mincut := mini(9, len(list))
		list = list[mincut:]
	}

}

func TestGetASNsbyList(t *testing.T) {
	asnList := []string{"196922", "3356"}

	peerDB := peeringdb.Peeringdb("https://www.peeringdb.com/api", apikey)
	peers, err := peerDB.GetASNsbyList(asnList)

	if err != nil {
		t.Errorf("Cant query the API for peers on IX: %s", err)
		return
	}

	if peers.Data[0].Asn != 3356 {
		t.Error("Cant find Lumen in the peering db?")
	}

	if peers.Data[1].Asn != 196922 {
		t.Error("Cant find myself in the peering db?")
	}

}
