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
	prefixFiles []string
}

func NewPeerGen(style, templateDir string, configDir string) *Peergen {
	return &Peergen{style: style, templateDir: templateDir,
		peerFiles: []string{
			configDir + "/" + style + "/user.tt",
			templateDir + "/" + style + "/header.tt",
			templateDir + "/" + style + "/router.tt",
			templateDir + "/" + style + "/footer.tt",
		},
		prefixFiles: []string{
			templateDir + "/" + style + "/prefix.tt",
			templateDir + "/" + style + "/prefix6.tt",
		},
	}
}

func (p *Peergen) GenerateIXs(exchanges ixtypes.IXs, w io.Writer) {

	if p.style == "juniper/json" {
		p.ConvertIxToJuniperJSON(exchanges, w, false)
		return
	} else if p.style == "juniper/json_pretty" {
		p.ConvertIxToJuniperJSON(exchanges, w, true)
		return
	} else if p.style == "extreme/slx_json" {
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

func (p *Peergen) GenerateIXPrefixFilter(exchanges ixtypes.IXs, w io.Writer) {
	for k := range exchanges {
		err := p.GeneratePrefixFilter(exchanges[k], w)
		if err != nil {
			log.Print(err)
		}
	}
}

func (p *Peergen) GeneratePrefixFilter(ix ixtypes.Ix, w io.Writer) error {
	var seenFilter = make(map[string]bool)

	for i := range p.prefixFiles {
		_, err := os.Stat(p.prefixFiles[i])
		if err != nil {
			continue
		}

		t, err := template.ParseFiles(p.prefixFiles[i])
		if err != nil {
			return fmt.Errorf("Cant open template file: %s", err)
		}

		/* We try not to output the same prefixFilter for the same IX ,
		but this code will limit us for "one"-shot on one template file
		*/
		for _, peer := range ix.PeersReady {
			if !seenFilter[peer.PrefixFilters.Name] || !seenFilter[peer.PrefixFilters6.Name] {
				if err := t.Execute(w, peer); err != nil {
					return fmt.Errorf("Cant execute template: %s", err)
				}
				seenFilter[peer.PrefixFilters.Name] = true
				seenFilter[peer.PrefixFilters6.Name] = true
			}
		}
	}
	return nil
}
