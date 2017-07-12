package bgpqworkers

import (
	"fmt"
	"io"
	"os/exec"
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
	return BGPQ3Worker{BGPQ3Config: Config}
}

func (b *BGPQ3Worker) GenPrefixList(w io.Writer, prefixListName, asMacro string, ipProtocol int) error {
	var cmd *exec.Cmd
	var err error
	var styleParameter = "-j"
	var ipParameter string

	if ipProtocol == 4 {
		ipParameter = "-4"
	} else {
		ipParameter = "-6"
	}

	if asMacro == "" {
		return fmt.Errorf("No valid AS or Macro given on command line: %s", asMacro)
	}

	cmd = exec.Command(b.Executable, ipParameter, styleParameter, "-l", prefixListName, asMacro)
	cmd.Stdout = w
	cmd.Stderr = w

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil

}
