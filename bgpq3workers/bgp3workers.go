package bgpqworkers

import (
	"os/exec"
	"io"
	"fmt"
)

type BGPQ3Config struct {
	Executable string
	Arguments  []string
	Style      string
}

type BGPQ3Worker struct {
	BGPQ3Config
}

func NewBGPQ3Worker(Config BGPQ3Config) BGPQ3Worker {
	return BGPQ3Worker{BGPQ3Config:Config}
}

func (b *BGPQ3Worker) GenPrefixList(w io.Writer, prefixListName, asMacro string, ipProtocol int) error {
	var cmd *exec.Cmd
	var err error
	var styleParameter, ipParameter string

	switch b.Style {
	case "junos/set":
		styleParameter = "-J"
	case "junos/json":
		styleParameter = "-j"
	case "cisco/iosxr":
		styleParameter = "-X"
	case "bird":
		styleParameter = "b"
	default:
		// Brocade SLX / MLX / IOS
		// no styleParameter necessary
	}

	if ipProtocol == 4 {
		ipParameter = "-4"
	} else {
		ipParameter = "-6"
	}

	if asMacro == "" {
		return fmt.Errorf("No valid AS or Macro given on command line: %s", asMacro)
	}

	if styleParameter == "" {
		cmd  = exec.Command(b.Executable, ipParameter,  "-l", prefixListName, asMacro)
	} else {
		cmd  = exec.Command(b.Executable, ipParameter,  styleParameter, "-l", prefixListName, asMacro)
	}

	cmd.Stdout = w
	cmd.Stderr = w

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil

}
