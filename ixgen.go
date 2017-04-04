package main

// ixgen (C) 2017 by Jörg Kost, joerg.kost@gmx.com
// see LICENSE for LICENSING,  TERMS AND CONDITIONS

import (
	"github.com/ipcjk/ixgen/inireader"
	"github.com/ipcjk/ixgen/ixtypes"
	"github.com/ipcjk/ixgen/peergen"
	"github.com/ipcjk/ixgen/peeringdb"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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
var apiservice bool
var localAPIServer string
var apiServiceURL string

var wg sync.WaitGroup

func init() {
	flag.StringVar(&peeringConfigFileName, "peerconfig", "./configuration/peering.ini", "Path to peering configuration ini-file")
	flag.StringVar(&peerStyleGenerator, "style", "brocade/netiron", "Style for routing-config by template, e.g. brocade, juniper, cisco...")
	flag.StringVar(&templateDir, "templates", "./templates", "directory for templates")
	flag.StringVar(&cacheDirectory, "cacheDir", "./cache", "cache directory for json files from peeringdb")
	flag.StringVar(&exchangeOnly, "exchange", "", "only generate configuration for given exchange, default: print all")
	flag.StringVar(&outputFile, "output", "", "if set, will output the configuration to a file, else STDOUT")
	flag.BoolVar(&printOrExit, "print", true, "print or just do a performance run")
	flag.BoolVar(&buildCache, "buildcache", false, "download json files for caching from peeringdb")
	flag.Int64Var(&myASN, "myasn", 0, "exclude your own asn from the generator")
	flag.BoolVar(&version, "version", false, "prints version and exit")

	/* Api-service-thread on localhost */
	flag.BoolVar(&apiservice, "apiservice", false, "create a local thread for the http api server that uses the json file as sources instead peeringdb.com/api-service.")
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
	workerMergePeerConfiguration()
	wg.Wait()
	if printOrExit == true {
		workerPrintPeerConfiguration()
	}
}

func workerPrintPeerConfiguration() {
	var outputStream io.WriteCloser
	var err error
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
	for k := range exchanges {
		peerGenerator.GenerateIX(exchanges[k], outputStream)
	}

}

func workerMergePeerConfiguration() {
	peerDB := peeringdb.Peeringdb(apiServiceURL)
	wg.Add(len(exchanges))
	for k := range exchanges {
		var i = k
		go func() {
			defer wg.Done()
			if exchangeOnly != "" && exchangeOnly != exchanges[i].IxName {
				return
			}
			_, rs_auto := exchanges[i].Options[exchanges[i].IxName]["routeserver"]

			myPeers := peerDB.GetPeersOnIXByIxName(exchanges[i].IxName)
			for _, peer := range myPeers.Data {
				peerASN := strconv.FormatInt(peer.Asn, 10)
				if peer.Asn == myASN {
					continue
				}
				peerDbNetwork := peerDB.GetNetworkByAsN(peer.Asn)
				if len(peerDbNetwork.Data) != 1 {
					continue
				}
				if peerDbNetwork.Data[0].InfoType == "Route Server" {
					rgroup, rgOk := exchanges[i].Options[exchanges[i].IxName]["routeserver_group"]
					rgroup6, rg6Ok := exchanges[i].Options[exchanges[i].IxName]["routeserver_group6"]
					rsPeer :=
						ixtypes.ExchangePeer{
							ASN:           peerASN,
							Active:        true,
							Ipv4Enabled:   false,
							Ipv6Enabled:   false,
							PrefixFilter:  false,
							GroupEnabled:  false,
							Group6Enabled: false,
							IsRs:          true, IsRsPeer: false,
							Ipv4Addr: peer.Ipaddr4,
							Ipv6Addr: peer.Ipaddr6,
						}
					if peer.Ipaddr6 != nil {
						rsPeer.Ipv6Enabled = true
					}
					if peer.Ipaddr4 != nil {
						rsPeer.Ipv4Enabled = true
					}
					if rgOk {
						rsPeer.Group = string(rgroup)
						rsPeer.GroupEnabled = true
					}
					if rg6Ok {
						rsPeer.Group6 = string(rgroup6)
						rsPeer.Group6Enabled = true
					}
					rsPeer.InfoPrefixes4 = peerDbNetwork.Data[0].InfoPrefixes4
					rsPeer.InfoPrefixes6 = peerDbNetwork.Data[0].InfoPrefixes6
					if rs_auto {
						exchanges[i].PeersReady = append(exchanges[i].PeersReady, rsPeer)
					}
					continue
				}
				_, ok := exchanges[i].PeersINI[exchanges[i].IxName][peerASN]
				pg, pgOk := exchanges[i].Options[exchanges[i].IxName]["peer_group"]
				pg6, pg6Ok := exchanges[i].Options[exchanges[i].IxName]["peer_group6"]
				// ok = this peer is configured in the INI-file
				if ok {
					/* Fix me, add support for looping same ASN with different configs */
					confPeer := exchanges[i].PeersINI[exchanges[i].IxName][peerASN][0]

					/* take care of values from the INI-file */
					if confPeer.InfoPrefixes4 == 0 {
						confPeer.InfoPrefixes4 = peerDbNetwork.Data[0].InfoPrefixes4
					}
					if confPeer.InfoPrefixes6 == 0 {
						confPeer.InfoPrefixes6 = peerDbNetwork.Data[0].InfoPrefixes6
					}
					if confPeer.Ipv6Addr == nil {
						confPeer.Ipv6Addr = peer.Ipaddr6
					}
					if confPeer.Ipv4Addr == nil {
						confPeer.Ipv4Addr = peer.Ipaddr4
					}

					if pgOk && confPeer.GroupEnabled == true && confPeer.Group == "" {
						confPeer.Group = string(pg)
					} else if confPeer.Group == "" {
						confPeer.GroupEnabled = false
					}
					if pg6Ok && confPeer.Group6Enabled == true && confPeer.Group6 == "" {
						confPeer.Group6 = string(pg6)
					} else if confPeer.Group6 == "" {
						confPeer.Group6Enabled = false
					}
					exchanges[i].PeersReady = append(exchanges[i].PeersReady, confPeer)
				} else if exchanges[i].Options[exchanges[i].IxName]["wildcard"] == "1" {
					// Wildcard, we take everything
					wildPeer :=
						ixtypes.ExchangePeer{
							ASN:          peerASN,
							Active:       true,
							Ipv4Enabled:  false,
							Ipv6Enabled:  false,
							PrefixFilter: false,
							Ipv4Addr:     peer.Ipaddr4,
							Ipv6Addr:     peer.Ipaddr6,
						}
					if peer.Ipaddr6 != nil {
						wildPeer.Ipv6Enabled = true
					}
					if peer.Ipaddr4 != nil {
						wildPeer.Ipv4Enabled = true
					}
					if pgOk && wildPeer.GroupEnabled == true && wildPeer.Group == "" {
						wildPeer.Group = string(pg)
					} else {
						wildPeer.GroupEnabled = false
					}
					if pg6Ok && wildPeer.Group6Enabled == true && wildPeer.Group6 == "" {
						wildPeer.Group6 = string(pg6)
					} else {
						wildPeer.Group6Enabled = false
					}

					if peerDbNetwork.Data != nil {
						wildPeer.InfoPrefixes4 = peerDbNetwork.Data[0].InfoPrefixes4
						wildPeer.InfoPrefixes6 = peerDbNetwork.Data[0].InfoPrefixes6
						exchanges[i].PeersReady = append(exchanges[i].PeersReady, wildPeer)
					} else {
						log.Printf("No data for ASN %s", peerASN)
					}
				}
			}
		}()
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
	if apiservice == true {
		Apiserver := peeringdb.NewAPIServer(localAPIServer, cacheDirectory)
		Apiserver.RunAPIServer()
		apiServiceURL = fmt.Sprintf("http://%s/api", Apiserver.AddrPort)
	}
}
