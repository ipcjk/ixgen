package peergen

import (
	"fmt"
	"github.com/ipcjk/ixgen/ixtypes"
	"html/template"
	"io"
	"log"
	"os"
)

type Peergen struct {
	style       string
	templateDir string
	peerFiles   []string
}

func NewPeerGen(style, templateDir string) *Peergen {
	return &Peergen{style, templateDir,
		[]string{
			templateDir + "/" + style + "/header.tt",
			templateDir + "/" + style + "/router.tt",
			templateDir + "/" + style + "/footer.tt",
		}}
}

func (p *Peergen) GenerateIXs(exchanges ixtypes.IXs, w io.Writer) {

	if p.style == "juniper/json" {
		p.ConvertIxToJuniperJSON(exchanges, w)
		return
	} else if p.style == "brocade/slx_json" {
		p.ConvertIxToBrocadeSlxJSON(exchanges, w)
		return
	} else if p.style == "native/json" {
		p.ConvertIxToJson(exchanges, w)
		return
	} else if p.style == "native/json_pretty" {
		p.ConvertIxToJsonPretty(exchanges, w)
		return
	}

	for k := range exchanges {
		err := p.GenerateIXConfiguration(exchanges[k], w)
		if err != nil {
			log.Print(err)
		}
	}
}

func (p *Peergen) GenerateIXConfiguration(ix ixtypes.Ix, w io.Writer) error {
	for i := range p.peerFiles {
		_, err := os.Stat(p.peerFiles[i])
		if err != nil {
			continue
		}

		t, err := template.ParseFiles(p.peerFiles[i])
		if err != nil {
			return fmt.Errorf("Cant open template file: %s", err)
		}

		if err := t.Execute(w, ix); err != nil {
			return fmt.Errorf("Cant execute template: %s", err)
		}
	}
	return nil
}
