package main

// ixgen (C) 2017 by Jörg Kost, jk@ip-clear.de
// see LICENSE for LICENSING,  TERMS AND CONDITIONS

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/ipcjk/ixgen/inireader"
	"github.com/ipcjk/ixgen/ixtypes"
	"github.com/ipcjk/ixgen/ixworkers"
	"github.com/ipcjk/ixgen/libapiserver"
	"github.com/ipcjk/ixgen/peergen"
)

/* Some globals for flag-parsing */
var exchangeOnly string
var peeringConfigFileName string
var peerStyleGenerator string
var templateDir, configDir string
var myASN int64
var outputFile string
var exchanges ixtypes.IXs
var peerGenerator *peergen.Peergen
var printOrExit, buildCache, version bool
var prefixFactor float64
var bgpqVersion int

/* Api server / uri */
var cacheDirectory string
var noapiservice bool
var localAPIServer string
var apiServiceURL string
var peeringDBAPIKey string

/* profile vars */
var cpuprofile, memprofile string

func readArgumentsAndSetup() {
	flag.StringVar(&peeringConfigFileName, "config", "./configuration/peering.ini", "Path to peering configuration ini-file")
	flag.StringVar(&peerStyleGenerator, "style", "extreme/netiron", "Style for routing-config by template, e.g. extreme, juniper, cisco. Also possible: native/json or native/json_pretty for outputting the inside structures")
	flag.StringVar(&templateDir, "templates", "./templates", "directory for templates")
	flag.StringVar(&configDir, "configpath", "./configuration", "directory for user.tt to include for reach router type or router")
	flag.StringVar(&cacheDirectory, "cacheDir", "./cache", "cache directory for json files from peeringdb")
	flag.StringVar(&exchangeOnly, "exchange", "", "only generate configuration for given exchange, default: print all")
	flag.StringVar(&outputFile, "output", "", "if set, will output the configuration to a file, else STDOUT")
	flag.BoolVar(&printOrExit, "exit", false, "exit after generator run,  before printing result (performance run)")
	flag.BoolVar(&buildCache, "buildcache", false, "download json files for caching from peeringdb")
	flag.Int64Var(&myASN, "myasn", 0, "exclude your own asn from the generator")
	flag.BoolVar(&version, "version", false, "prints version and exit")
	flag.Float64Var(&prefixFactor, "prefixfactor", 1.2, "factor for maximum-prefix numbers")

	/* Api-service-thread on localhost */
	flag.BoolVar(&noapiservice, "noapiservice", false, "do NOT create a local thread for the http api server that uses the json file as sources instead peeringdb.com/api-service.")
	flag.StringVar(&localAPIServer, "listenapi", "localhost:0", "listenAddr for local api service")
	flag.StringVar(&apiServiceURL, "api", "https://www.peeringdb.com/api", "use a differnt server as sources instead local/api-service.")
	flag.StringVar(&peeringDBAPIKey, "apikey", "", "Peering DB API-Key")
	flag.IntVar(&bgpqVersion, "bgpq", 3, "BGPQ version to use (3 or 4)")

	/* profiling support */
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to `file`")
	flag.StringVar(&memprofile, "memprofile", "", "write memory profile to `file`")

	flag.Parse()

	if version {
		log.Println("ixgen 0.8 (C) 2025 by Jörg Kost, jk@ip-clear.de")
		os.Exit(0)
	}

	if peeringDBAPIKey == "" {
		peeringDBAPIKey = os.Getenv("PEERINGDB_APIKEY")
	}

	if buildCache {
		libapiserver.DownloadCache("https://www.peeringdb.com/api", cacheDirectory, peeringDBAPIKey)
		os.Exit(0)
	}

	loadConfig()
}

func main() {
	var outputStream io.WriteCloser

	readArgumentsAndSetup()

	/* profile support */
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err = pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		defer f.Close()
	}

	/* Merge PeeringDB */
	exchanges = ixworkers.WorkerMergePeerConfiguration(exchanges, apiServiceURL, peeringDBAPIKey, exchangeOnly, myASN, prefixFactor)
	/* Merge BGPq prefixFilters if we are on Mac or Linux */
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		exchanges = ixworkers.WorkerMergePrefixFilters(exchanges, exchangeOnly, bgpqVersion)
	}

	if !printOrExit {
		var err error
		outputStream = os.Stdout
		if outputFile == "" {
			defer os.Stdout.Close()
		} else {
			outputStream, err = os.Create(outputFile)
			if err != nil {
				log.Fatal(err)
			}
		}
		defer outputStream.Close()
		peerGenerator.GenerateIXPrefixFilter(exchanges, outputStream)
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

	if strings.HasSuffix(peeringConfigFileName, ".json") {
		err := json.NewDecoder(file).Decode(&exchanges)
		if err != nil {
			log.Fatalf("JSON input file malformed: %s", err)
			return
		}
	} else {
		exchanges = inireader.ReadPeeringConfig(file)
	}

	peerGenerator = peergen.NewPeerGen(peerStyleGenerator, templateDir, configDir)
	if !noapiservice {
		Apiserver := libapiserver.NewAPIServer(localAPIServer, cacheDirectory, templateDir, configDir)
		Apiserver.RunAPIServer()
		apiServiceURL = fmt.Sprintf("http://%s/api", Apiserver.AddrPort)
	}
}
