package main

import (
	"IXgenerator/inireader"
	"IXgenerator/ixtypes"
	"testing"
)

func TestParsePeerFunction(t *testing.T) {
	Line := "196922 ipv4=1 ipv6=0 active=1 prefix_filter=1 local_pref=910 group6=0 group=0 peer_group6=mygroup6 peer_group=mygroup4"
	Peer := inireader.ParsePeerLine(Line, 0)

	if Peer.ASN != "196922" {
		t.Error("Peer ASN is different than expected")
	}
	if Peer.PrefixFilter != true {
		t.Error("Peer Prefix Filter is different than expected")
	}
	if Peer.Ipv6Enabled != false {
		t.Error("Peer Ipv6 value is different than expected")
	}
	if Peer.Ipv4Enabled != true {
		t.Error("Peer Ipv4 value is  different than expected")
	}
	if Peer.Active != true {
		t.Error("Peer Active value is  different than expected")
	}
	if Peer.GroupEnabled != false {
		t.Error("Peer Group enabled/disabled is  different than expected")
	}
	if Peer.Group6Enabled != false {
		t.Error("Peer Group6 enabled/disabled is  different than expected")
	}
	if Peer.Group != "mygroup4" {
		t.Error("Peer Group string is  different than expected")
	}
	if Peer.Group6 != "mygroup6" {
		t.Error("Peer Group6 string is  different than expected")
	}
}

func TestExchangeOption(t *testing.T) {
	/* Init test object */
	var ixConfig = make(ixtypes.ExchangeOptions)
	ixConfig["testIX"] = make(map[string]ixtypes.ExchangeOption, 10)
	inireader.ParseOptionLine("routeserver=1", ixConfig, "testIX")
	inireader.ParseOptionLine("routeserver6=1", ixConfig, "testIX")
	inireader.ParseOptionLine("routeserver_group=IX", ixConfig, "testIX")
	inireader.ParseOptionLine("routeserver_group6=IX6", ixConfig, "testIX")
	inireader.ParseOptionLine("peer_group=peer", ixConfig, "testIX")
	inireader.ParseOptionLine("peer_group6=peer6", ixConfig, "testIX")
	inireader.ParseOptionLine("wildcard=1", ixConfig, "testIX")

	/* Check that we covered all cases from inireader */
	for k := range inireader.PossibleOptions {
		_, ok := ixConfig["testIX"][k]
		if !ok {
			t.Errorf("option  %s not readable ", k)
		}
	}

	/* Check static for every case we give above */
	if ixConfig["testIX"]["routeserver"] != "1" {
		t.Error("Routeserver option is wrong")
	}

	if ixConfig["testIX"]["routeserver6"] != "1" {
		t.Error("Routeserver option is wrong")
	}
	if ixConfig["testIX"]["routeserver_group"] != "IX" {
		t.Error("Routeserver Group option is wrong")
	}

	if ixConfig["testIX"]["routeserver_group6"] != "IX6" {
		t.Error("Routeserver Group6 option is wrong")
	}

	if ixConfig["testIX"]["peer_group"] != "peer" {
		t.Error("Peer Group option is wrong")
	}

	if ixConfig["testIX"]["peer_group6"] != "peer6" {
		t.Error("Peer Group6 option is wrong")
	}

	if ixConfig["testIX"]["wildcard"] != "1" {
		t.Error("Wildcard option is wrong")
	}
}
