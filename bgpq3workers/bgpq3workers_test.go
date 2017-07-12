package bgpqworkers_test

import (
	"bytes"
	. "github.com/ipcjk/ixgen/bgpq3workers"
	"strings"
	"testing"
)

func TestRunBGP3Worker(t *testing.T) {
	w := new(bytes.Buffer)
	testAsMacro := "AS196922"

	Config := BGPQ3Config{
		Executable: "/Users/joerg/Documents/Programmierung/bgpq3-0.1.21/bgpq3",
		Style:      "brocade/mlx",
		Arguments:  []string{"-4"},
	}

	bgpWorker := NewBGPQ3Worker(Config)

	err := bgpWorker.GenPrefixList(w, "as196922-p4", testAsMacro, 4)
	if err != nil {
		t.Errorf("Cant run bgpq3: %s", err)
	}

	if !strings.Contains(w.String(), "178.248.240.0") {
		t.Error("Cant find my home ipv4 prefix list")
	}

	if !strings.Contains(w.String(), "as196922-p4") {
		t.Error("Expected prefixListName  as196922-p4, but did not find")
	}

	err = bgpWorker.GenPrefixList(w, "as196922-p6", testAsMacro, 6)
	if err != nil {
		t.Errorf("Cant run bgpq3: %s", err)
	}

	if !strings.Contains(w.String(), "2a02:1308") {
		t.Error("Cant find my home ipv6 prefix list")
	}

	if !strings.Contains(w.String(), "as196922-p6") {
		t.Error("Expected prefixListName  as196922-p6, but did not find")
	}

}
