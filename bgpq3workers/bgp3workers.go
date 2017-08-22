package bgpqworkers

import (
	"bytes"
	"encoding/json"
	"github.com/ipcjk/ixgen/ixtypes"
	"log"
	"os"
	"os/exec"
	"runtime"
)

type BGPQ3Config struct {
	Executable string
}

type BGPQ3Worker struct {
	BGPQ3Config
}

func NewBGPQ3Worker(Config BGPQ3Config) BGPQ3Worker {

	if runtime.GOOS == "linux" {
		Config.Executable = "bgpq3.linux"
	} else if runtime.GOOS == "darwin" {
		Config.Executable = "bgpq3.mac"
	}

	Config.Executable = findExecutable(Config.Executable)

	return BGPQ3Worker{BGPQ3Config: Config}
}

func (b *BGPQ3Worker) GenPrefixList(prefixListName, asMacro string, ipProtocol int, aggregateMax bool) (ixtypes.PrefixFilters, error) {
	var w = new(bytes.Buffer)
	var prefixFilters ixtypes.PrefixFilters
	var extraArgs []string

	if ipProtocol == 4 {
		extraArgs = append(extraArgs, "-4")
	} else {
		extraArgs = append(extraArgs, "-6")
	}

	if aggregateMax == true {
		extraArgs = append(extraArgs, "-A")
	}

	extraArgs = append(extraArgs, "-j", "-l", prefixListName, asMacro)

	cmd := exec.Command(b.Executable, extraArgs[0:]...)
	cmd.Stdout = w
	cmd.Stderr = w

	err := cmd.Run()
	if err != nil {
		return ixtypes.PrefixFilters{}, err
	}

	err = json.Unmarshal(w.Bytes(), &prefixFilters)
	if err != nil {
		return ixtypes.PrefixFilters{}, err
	}

	return prefixFilters, nil

}

func findExecutable(execName string) string {
	if _, err := os.Stat(execName); err == nil {
		return "./" + execName
	} else if _, err := os.Stat("../" + execName); err == nil {
		return "../" + execName
	} else if path, err := exec.LookPath(execName); err == nil {
		return path
	} else {
		log.Fatalf("Cant find " + execName)
	}
	return execName
}
