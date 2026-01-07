package bgpqworkers

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	if len(rawOutput) == 0 {
		return filters, fmt.Errorf("empty response from bgpq command")
	}

	for name, rules := range rawOutput {
		filters.Name = name
		filters.PrefixRules = rules
		break
	}

	return filters, nil
}

func NewBGPQ3Worker(path string) *BGPQ3Worker {
	var executable string
	if path != "" {
		if _, err := os.Stat(path); err == nil {
			// check if path is executable on linux or mac
			executable = path
		} else {
			log.Printf("Warning: %s is not a valid executable file, will use default executable", path)
			executable = findExecutable(getExecutableName("bgpq3"))
		}
	} else {
		executable = findExecutable(getExecutableName("bgpq3"))
	}

	return &BGPQ3Worker{
		baseBGPQWorker{config: BGPQConfig{
			Executable: executable,
			Version:    3,
		}},
	}
}

// GenPrefixList generates the prefix lists in bgqp3 way
func (w *BGPQ3Worker) GenPrefixList(prefixListName, asMacro string, ipProtocol int, aggregateMax bool) (ixtypes.PrefixFilters, error) {
	args := buildCommonArgs(prefixListName, asMacro, ipProtocol, aggregateMax)
	if len(args) == 0 {
		return ixtypes.PrefixFilters{}, fmt.Errorf("invalid arguments for prefix list generation")
	}
	return w.runCommand(args)
}

type BGPQ4Worker struct {
	baseBGPQWorker
}

func NewBGPQ4Worker(path string) *BGPQ4Worker {
	var executable string
	if path != "" {
		if _, err := os.Stat(path); err == nil {
			// check if path is executable on linux or mac
			executable = path
		} else {
			log.Printf("Warning: %s is not a valid executable file, will use default executable", path)
			executable = findExecutable(getExecutableName("bgpq4"))
		}
	} else {
		executable = findExecutable(getExecutableName("bgpq4"))
	}

	return &BGPQ4Worker{
		baseBGPQWorker{config: BGPQConfig{
			Executable: executable,
			Version:    4,
		}},
	}
}

func NewBGPQWorker(version int, path string) BGPQ {
	if version == 4 {
		return NewBGPQ4Worker(path)
	}
	return NewBGPQ3Worker(path)
}

// GenPrefixList generates the prefix lists in bgqp4 way
func (w *BGPQ4Worker) GenPrefixList(prefixListName, asMacro string, ipProtocol int, aggregateMax bool) (ixtypes.PrefixFilters, error) {
	args := buildCommonArgs(prefixListName, asMacro, ipProtocol, aggregateMax)
	if len(args) == 0 {
		return ixtypes.PrefixFilters{}, fmt.Errorf("invalid arguments for prefix list generation")
	}
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
			source := parts[0]
			if source == "" {
				log.Printf("Warning: Empty source in IRR format: %s", IrrAsSet)
				return []string{}
			}
			// Split AS sets by : and pass them as separate arguments
			asSets := strings.Split(parts[1], ":")
			// Filter out empty AS sets
			var validAsSets []string
			for _, asSet := range asSets {
				if strings.TrimSpace(asSet) != "" {
					validAsSets = append(validAsSets, asSet)
				}
			}
			if len(validAsSets) == 0 {
				log.Printf("Warning: No valid AS sets found in: %s", IrrAsSet)
				return []string{}
			}
			args = append(args, "-S", source)
			args = append(args, "-j", "-l", prefixListName)
			args = append(args, validAsSets...)
			return args
		} else {
			log.Printf("Warning: Invalid/Unknown IRR source format: %s", IrrAsSet)
			// Return empty args instead of nil to avoid panic
			return []string{}
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
