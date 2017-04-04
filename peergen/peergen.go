package peergen

import (
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

func (p *Peergen) GenerateIX(routerTemplate ixtypes.Ix, w io.Writer) {
	if p.style == "juniper/json" {
		p.ConvertIxToJuniperJSON(routerTemplate, w)
		return
	} else if p.style == "brocade/slx" {
		p.ConvertIxToBrocadeSlxJSON(routerTemplate, w)
		return
	}
	p.GenerateIXConfiguration(routerTemplate, w)
}

func (p *Peergen) GenerateIXConfiguration(routerTemplate ixtypes.Ix, w io.Writer) {
	for i := range p.peerFiles {
		_, err := os.Stat(p.peerFiles[i])
		if err != nil {
			continue
		}

		t, err := template.ParseFiles(p.peerFiles[i])
		if err != nil {
			log.Fatalf("Cant open template file: %s", err)
		}

		if err := t.Execute(w, routerTemplate); err != nil {
			log.Fatalf("Cant execute template: %s", err)
		}
	}

}

func (p *Peergen) GeneratePeerConfiguration(routerTemplate ixtypes.PeerTemplate, w io.Writer) {
	for i := range p.peerFiles {
		t, err := template.ParseFiles(p.peerFiles[i])
		if err != nil {
			log.Fatalf("Cant open template file: %s", err)
		}

		if err := t.Execute(w, routerTemplate); err != nil {
			log.Fatalf("Cant execute template: %s", err)
		}
	}
}
