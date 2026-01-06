package bgpqworkers_test

import (
	"testing"

	. "github.com/ipcjk/ixgen/bgpqworkers"
)

func TestRunBGPQ3Worker(t *testing.T) {
	testAsMacro := "AS196922"
	bgpWorker := NewBGPQ3Worker()

	prefixFilters, err := bgpWorker.GenPrefixList("as196922p4", testAsMacro, 4, true)
	if err != nil {
		t.Errorf("Cant run bgpq3: %s", err)
	}

	if prefixFilters.Name != "as196922p4" {
		t.Errorf("Expected prefix name 'as196922p4', got '%s'", prefixFilters.Name)
	}

	if len(prefixFilters.PrefixRules) <= 2 ||
		len(prefixFilters.PrefixRules) >= 30 {
		t.Errorf("Unexpected number of prefix rules: %d", len(prefixFilters.PrefixRules))
	}

	prefixFilters, err = bgpWorker.GenPrefixList("as196922p6", testAsMacro, 6, false)
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

func TestRunBGPQ3WorkerWithSource(t *testing.T) {
	testAsMacro := "RIPE::AS16097:AS-HLKOMM"
	bgpWorker := NewBGPQ3Worker()

	prefixFilters, err := bgpWorker.GenPrefixList("as16097p4", testAsMacro, 4, true)
	if err != nil {
		t.Errorf("Cant run bgpq3: %s", err)
	}

	if prefixFilters.Name != "as16097p4" {
		t.Errorf("Expected prefix name 'as16097p4', got '%s'", prefixFilters.Name)
	}

	if len(prefixFilters.PrefixRules) <= 2 ||
		len(prefixFilters.PrefixRules) >= 1000 {
		t.Errorf("Unexpected number of prefix rules: %d", len(prefixFilters.PrefixRules))
	}

	prefixFilters, err = bgpWorker.GenPrefixList("as16097p6", testAsMacro, 6, false)
	if err != nil {
		t.Errorf("Cant run bgpq3: %s", err)
	}

	if prefixFilters.Name != "as16097p6" {
		t.Error("Cant find my home ipv6 PrefixName")
	}

	if len(prefixFilters.PrefixRules) == 0 ||
		len(prefixFilters.PrefixRules) > 400 {
		t.Error("Found too less or too many ipv6 prefixes, cant be!")
	}

}

func TestRunBGPQ4Worker(t *testing.T) {
	testAsMacro := "AS196922"
	bgpWorker := NewBGPQ4Worker()

	prefixFilters, err := bgpWorker.GenPrefixList("as196922p4", testAsMacro, 4, true)
	if err != nil {
		t.Errorf("Cant run bgpq4: %s", err)
	}

	if prefixFilters.Name != "as196922p4" {
		t.Errorf("Expected prefix name 'as196922p4', got '%s'", prefixFilters.Name)
	}

	if len(prefixFilters.PrefixRules) <= 2 ||
		len(prefixFilters.PrefixRules) >= 30 {
		t.Errorf("Unexpected number of prefix rules: %d", len(prefixFilters.PrefixRules))
	}

	prefixFilters, err = bgpWorker.GenPrefixList("as196922p6", testAsMacro, 6, false)
	if err != nil {
		t.Errorf("Cant run bgpq4: %s", err)
	}

	if prefixFilters.Name != "as196922p6" {
		t.Error("Cant find my home ipv6 PrefixName")
	}

	if len(prefixFilters.PrefixRules) == 0 ||
		len(prefixFilters.PrefixRules) >= 400 {
		t.Error("Found too less or too many ipv6 prefixes, cant be!", len(prefixFilters.PrefixRules))
	}

}

func TestRunBGPQ4WorkerWithSource(t *testing.T) {
	testAsMacro := "RIPE::AS16097:AS-HLKOMM"
	bgpWorker := NewBGPQ4Worker()

	prefixFilters, err := bgpWorker.GenPrefixList("as16097p4", testAsMacro, 4, true)
	if err != nil {
		t.Errorf("Cant run bgpq4: %s", err)
	}

	/* t.Logf("PrefixFilters: %+v", prefixFilters)
	t.Logf("PrefixFilters Name: %s", prefixFilters.Name)
	t.Logf("Number of PrefixRules: %d", len(prefixFilters.PrefixRules))
	t.Logf("PrefixRules: %+v", prefixFilters.PrefixRules)
	t.Logf("PrefixRules[0]: %+v", prefixFilters.PrefixRules[0])
	*/

	if prefixFilters.Name != "as16097p4" {
		t.Errorf("Expected prefix name 'as16097p4', got '%s'", prefixFilters.Name)
	}

	if len(prefixFilters.PrefixRules) <= 2 ||
		len(prefixFilters.PrefixRules) >= 1000 {
		t.Errorf("Unexpected number of prefix rules: %d", len(prefixFilters.PrefixRules))
	}

	prefixFilters, err = bgpWorker.GenPrefixList("as16097p6", testAsMacro, 6, false)
	if err != nil {
		t.Errorf("Cant run bgpq4: %s", err)
	}

	if prefixFilters.Name != "as16097p6" {
		t.Error("Cant find my home ipv6 PrefixName")
	}

	if len(prefixFilters.PrefixRules) == 0 ||
		len(prefixFilters.PrefixRules) > 400 {
		t.Error("Found too less or too many ipv6 prefixes, cant be!", len(prefixFilters.PrefixRules))
	}
}
