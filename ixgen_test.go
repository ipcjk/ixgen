package main

import (
	"reflect"
	"regexp"
	"testing"
)

var RegExpPeering = `^(?P<ASN>\d+)\s?(?:ipv4=)?(?P<IPV4>\d+)?\s?(?:ipv6=)?(?P<IPV6>\d+)?\s?(?:active=)?(?P<ACTIVE>\d+)?\s?(?:filter=)?(?P<FILTER>\d+)?`

func TestRegExpForPeering(t *testing.T) {
	_, err := regexp.Compile(RegExpPeering)
	if err != nil {
		t.Error("Cant compile Regex")
	}
}

func TestRegExpSubExpNames(t *testing.T) {
	var MySlice = []string{
		"", "ASN", "IPV4", "IPV6", "ACTIVE", "FILTER"}

	matchPeer, err := regexp.Compile(RegExpPeering)
	if err != nil {
		t.Error("Cant compile Regex")
	}

	if reflect.DeepEqual(matchPeer.SubexpNames(), MySlice) == false {
		t.Error("SubExpNames not Equal with expected results", matchPeer.SubexpNames())
	}
}
