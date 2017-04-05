package main

// ixgen (C) 2017 by Jörg Kost, joerg.kost@gmx.com
// see LICENSE for LICENSING,  TERMS AND CONDITIONS

import (
	"flag"
	"fmt"
	"github.com/ipcjk/ixgen/inireader"
	"github.com/ipcjk/ixgen/ixtypes"
	"github.com/ipcjk/ixgen/peergen"
	"github.com/ipcjk/ixgen/peeringdb"
	"github.com/ipcjk/ixgen/ixworkers"
	"io"
	"log"
	"os"
	"sync"
)

/* Some globals for flag-parsing */
var exchangeOnly string
var peeringConfigFileName string
var peerStyleGenerator string
var templateDir string
var myASN int64
var outputFile string
var exchanges ixtypes.IXs
var peerGenerator *peergen.Peergen
var printOrExit bool
var buildCache bool
var version bool

/* Api server / uri */
var cacheDirectory string
var noapiservice bool
var localAPIServer string
var apiServiceURL string

var wg sync.WaitGroup

func init() {
	flag.StringVar(&peeringConfigFileName, "peerconfig", "./configuration/peering.ini", "Path to peering configuration ini-file")
	flag.StringVar(&peerStyleGenerator, "style", "brocade/netiron", "Style for routing-config by template, e.g. brocade, juniper, cisco. Also possible: native/json or native/json_pretty for outputting the inside structures")
	flag.StringVar(&templateDir, "templates", "./templates", "directory for templates")
	flag.StringVar(&cacheDirectory, "cacheDir", "./cache", "cache directory for json files from peeringdb")
	flag.StringVar(&exchangeOnly, "exchange", "", "only generate configuration for given exchange, default: print all")
	flag.StringVar(&outputFile, "output", "", "if set, will output the configuration to a file, else STDOUT")
	flag.BoolVar(&printOrExit, "exit", false, "exit after generator run,  before printing result (performance run)")
	flag.BoolVar(&buildCache, "buildcache", false, "download json files for caching from peeringdb")
	flag.Int64Var(&myASN, "myasn", 0, "exclude your own asn from the generator")
	flag.BoolVar(&version, "version", false, "prints version and exit")

	/* Api-service-thread on localhost */
	flag.BoolVar(&noapiservice, "noapiservice", false, "do NOT create a local thread for the http api server that uses the json file as sources instead peeringdb.com/api-service.")
	flag.StringVar(&localAPIServer, "listenapi", "localhost:0", "listenAddr for local api service")
	flag.StringVar(&apiServiceURL, "api", "https://www.peeringdb.com/api", "use a differnt server as sources instead local/api-service.")

	flag.Parse()

	if version == true {
		fmt.Println("ixgen 0.1 (C) 2017 by Jörg Kost, joerg.kost@gmx.com")
		os.Exit(0)
	}

	if buildCache == true {
		peeringdb.DownloadCache("https://www.peeringdb.com/api", cacheDirectory)
		os.Exit(0)
	}

	loadConfig()
}

func main() {
	var outputStream io.WriteCloser
	var err error

	wg.Add(len(exchanges))
	ixworkers.WorkerMergePeerConfiguration(&wg, exchanges, apiServiceURL, exchangeOnly, myASN)
	wg.Wait()

	if printOrExit == false {
		if outputFile == "" {
			outputStream = os.Stdout
			defer os.Stdout.Close()
		} else {
			outputStream, err = os.Create(outputFile)
			if err != nil {
				log.Fatal(err)
			}
		}
		defer outputStream.Close()
		peerGenerator.GenerateIXs(exchanges, outputStream)
	}
}

func loadConfig() {
	if len(peeringConfigFileName) == 0 {
		log.Fatal("No peering.ini given")
	}

	file, err := os.Open(peeringConfigFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	exchanges = inireader.ReadPeeringConfig(file)

	peerGenerator = peergen.NewPeerGen(peerStyleGenerator, templateDir)
	if noapiservice == false {
		Apiserver := peeringdb.NewAPIServer(localAPIServer, cacheDirectory)
		Apiserver.RunAPIServer()
		apiServiceURL = fmt.Sprintf("http://%s/api", Apiserver.AddrPort)
	}
}
