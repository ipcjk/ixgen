package bgpqworkers_test

import (
	. "github.com/ipcjk/ixgen/bgpq3workers"
	"testing"
)

func TestRunBGP3Worker(t *testing.T) {
	testAsMacro := "AS196922"

	Config := BGPQ3Config{
		Style:     "brocade/mlx",
		Arguments: []string{"-4"},
	}

	bgpWorker := NewBGPQ3Worker(Config)

	prefixFilters, err := bgpWorker.GenPrefixList("as196922p4", testAsMacro, 4)
	if err != nil {
		t.Errorf("Cant run bgpq3: %s", err)
	}

	if prefixFilters.Name != "as196922p4" {
		t.Error("Cant find my home ipv4 PrefixName")
	}

	if len(prefixFilters.PrefixRules) <= 4 ||
		len(prefixFilters.PrefixRules) >= 30 {
		t.Error("Found too less or too many ipv4 prefixes, cant be!")
	}

	prefixFilters, err = bgpWorker.GenPrefixList("as196922p6", testAsMacro, 6)
	if err != nil {
		t.Errorf("Cant run bgpq3: %s", err)
	}

	if prefixFilters.Name != "as196922p6" {
		t.Error("Cant find my home ipv6 PrefixName")
	}

	if len(prefixFilters.PrefixRules) == 0 ||
		len(prefixFilters.PrefixRules) >= 30 {
		t.Error("Found too less or too many ipv6 prefixes, cant be!")
	}

}
