package bgpqworkers

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ipcjk/ixgen/ixtypes"
)

type BGPQ interface {
	GenPrefixList(string, string, int, bool) (ixtypes.PrefixFilters, error)
}

type BGPQConfig struct {
	Executable string
	Version    int
}

type baseBGPQWorker struct {
	config BGPQConfig
}

type BGPQ3Worker struct {
	baseBGPQWorker
}

// helper to execute the command and parse the JSON output
func (b *baseBGPQWorker) runCommand(args []string) (ixtypes.PrefixFilters, error) {
	var out = new(bytes.Buffer)
	var filters ixtypes.PrefixFilters
	var rawOutput map[string][]ixtypes.PrefixRule

	cmd := exec.Command(b.config.Executable, args...)
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Run(); err != nil {
		return filters, err
	}

	err := json.Unmarshal(out.Bytes(), &rawOutput)
	if err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		return filters, err
	}

	// Extract the first (and only) key from the map
	for name, rules := range rawOutput {
		filters.Name = name
		filters.PrefixRules = rules
		break
	}

	return filters, nil
}

func NewBGPQ3Worker() *BGPQ3Worker {
	return &BGPQ3Worker{
		baseBGPQWorker{config: BGPQConfig{
			Executable: findExecutable(getExecutableName("bgpq3")),
			Version:    3,
		}},
	}
}

// GenPrefixList generates the prefix lists in bgqp3 way
func (w *BGPQ3Worker) GenPrefixList(prefixListName, asMacro string, ipProtocol int, aggregateMax bool) (ixtypes.PrefixFilters, error) {
	args := buildCommonArgs(prefixListName, asMacro, ipProtocol, aggregateMax)
	return w.runCommand(args)
}

type BGPQ4Worker struct {
	baseBGPQWorker
}

func NewBGPQ4Worker() *BGPQ4Worker {
	return &BGPQ4Worker{
		baseBGPQWorker{config: BGPQConfig{
			Executable: findExecutable(getExecutableName("bgpq4")),
			Version:    4,
		}},
	}
}

func NewBGPQWorker(version int) BGPQ {
	if version == 4 {
		return NewBGPQ4Worker()
	}
	return NewBGPQ3Worker()
}

// GenPrefixList generates the prefix lists in bgqp4 way
func (w *BGPQ4Worker) GenPrefixList(prefixListName, asMacro string, ipProtocol int, aggregateMax bool) (ixtypes.PrefixFilters, error) {
	args := buildCommonArgs(prefixListName, asMacro, ipProtocol, aggregateMax)
	return w.runCommand(args)
}

// helper to construct common args
func buildCommonArgs(prefixListName, IrrAsSet string, ipProtocol int, aggregateMax bool) []string {
	var args []string
	if ipProtocol == 4 {
		args = append(args, "-4")
	} else {
		args = append(args, "-6")
	}
	if aggregateMax {
		args = append(args, "-A")
	}

	// Handle IRR source format (e.g. RIPE::AS16097:AS-HLKOMM)
	// Format: SOURCE::AS-SET1:AS-SET2
	// SOURCE: IRR database (RIPE, ARIN, APNIC, ...)
	// AS-SETs: Colon-separated list of AS sets to query
	if strings.Contains(IrrAsSet, "::") {
		parts := strings.Split(IrrAsSet, "::")
		if len(parts) == 2 {
			IrrAsSet = parts[0]
			// Split AS sets by : and pass them as separate arguments
			asSets := strings.Split(parts[1], ":")
			args = append(args, "-S", IrrAsSet)
			args = append(args, "-j", "-l", prefixListName)
			args = append(args, asSets...)
			return args
		} else {
			log.Printf("Warning: Invalid/Unknown IRR source format: %s", IrrAsSet)
			return nil
		}
	}

	// Handle default format (e.g. AS-MYSET)
	args = append(args, "-j", "-l", prefixListName, IrrAsSet)
	return args
}

// helper to find executable
func findExecutable(execName string) string {
	if _, err := os.Stat(execName); err == nil {
		return "./" + execName
	} else if _, err := os.Stat("../" + execName); err == nil {
		return "../" + execName
	} else if path, err := exec.LookPath(execName); err == nil {
		return path
	} else {
		log.Printf("Can't find %s\n", execName)
	}
	return execName
}

// helper to get the right executable for the platform
func getExecutableName(base string) string {
	switch runtime.GOOS {
	case "linux":
		if runtime.GOARCH == "arm64" {
			return base + "_arm64.linux"
		}
		return base + ".linux"
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return base + "_arm64.mac"
		}
		return base + ".mac"
	default:
		log.Fatalf("Unsupported OS: %s", runtime.GOOS)
		return ""
	}
}
