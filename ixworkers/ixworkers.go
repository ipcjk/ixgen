package ixworkers

import (
	"github.com/ipcjk/ixgen/bgpq3workers"
	"github.com/ipcjk/ixgen/ixtypes"
	"github.com/ipcjk/ixgen/peeringdb"
	"log"
	"strconv"
	"sync"
)

func WorkerMergePeerConfiguration(exchanges ixtypes.IXs, apiServiceURL string, exchangeOnly string, myASN int64) ixtypes.IXs {
	var wg sync.WaitGroup
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
			rsnASN, rsnOk := exchanges[i].Options[exchanges[i].IxName]["rs_asn"]

			myPeers, err := peerDB.GetPeersOnIXByIxName(exchanges[i].IxName)
			if err != nil {
				log.Printf("Cant get Peers for IX: %s, error: %s", exchanges[i].IxName, err)
				return
			}
			for _, peer := range myPeers.Data {
				peerASN := strconv.FormatInt(peer.Asn, 10)
				if peer.Asn == myASN {
					continue
				}
				peerDbNetwork, err := peerDB.GetNetworkByAsN(peer.Asn)
				if err != nil {
					log.Printf("Error pulling ASN for peer: %s, error: %s", peer.Asn, err)
					continue
				}
				if len(peerDbNetwork.Data) != 1 {
					continue
				}
				if peerDbNetwork.Data[0].InfoType == "Route Server" {
					rgroup, rgOk := exchanges[i].Options[exchanges[i].IxName]["routeserver_group"]
					rgroup6, rg6Ok := exchanges[i].Options[exchanges[i].IxName]["routeserver_group6"]
					infoprefixes4, rprefixOk := exchanges[i].Options[exchanges[i].IxName]["routeserver_prefixes"]
					infoprefixes6, rprefix6Ok := exchanges[i].Options[exchanges[i].IxName]["routeserver_prefixes6"]

					rsPeer :=
						ixtypes.ExchangePeer{
							ASN:                 peerASN,
							Active:              true,
							Ipv4Enabled:         false,
							Ipv6Enabled:         false,
							PrefixFilterEnabled: false,
							GroupEnabled:        false,
							Group6Enabled:       false,
							IsRs:                true, IsRsPeer: false,
							Ipv4Addr: peer.Ipaddr4,
							Ipv6Addr: peer.Ipaddr6,
							IrrAsSet: peerDbNetwork.Data[0].IrrAsSet,
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

					if rprefixOk {
						rsPeer.InfoPrefixes4, _ = strconv.ParseInt(string(infoprefixes4), 10, 64)
					}

					if rprefix6Ok {
						rsPeer.InfoPrefixes6, _ = strconv.ParseInt(string(infoprefixes6), 10, 64)
					}

					if rs_auto && rsnOk && peerASN != string(rsnASN) {
						log.Printf("Probably rogue route-server advertised in %s\n", peerASN)
					} else if rs_auto {
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
					confPeer.IrrAsSet = peerDbNetwork.Data[0].IrrAsSet
					exchanges[i].PeersReady = append(exchanges[i].PeersReady, confPeer)
				} else if exchanges[i].Options[exchanges[i].IxName]["wildcard"] == "1" {
					// Wildcard, we take everything
					wildPeer :=
						ixtypes.ExchangePeer{
							ASN:                 peerASN,
							Active:              true,
							Ipv4Enabled:         false,
							Ipv6Enabled:         false,
							PrefixFilterEnabled: false,
							Ipv4Addr:            peer.Ipaddr4,
							Ipv6Addr:            peer.Ipaddr6,
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
						wildPeer.IrrAsSet = peerDbNetwork.Data[0].IrrAsSet
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
	wg.Wait()
	return exchanges
}

func WorkerMergePrefixFilters(exchanges ixtypes.IXs, exchangeOnly string) ixtypes.IXs {
	var wg sync.WaitGroup

	bgpWorker := bgpqworkers.NewBGPQ3Worker(bgpqworkers.BGPQ3Config{})

	wg.Add(len(exchanges))
	for k := range exchanges {
		var i = k
		go func() {
			defer wg.Done()
			if exchangeOnly != "" && exchangeOnly != exchanges[i].IxName {
				return
			}

			for j := range exchanges[i].PeersReady {
				var asMacro string
				var err error

				if exchanges[i].PeersReady[j].IsRs {
					continue
				}

				if exchanges[i].PeersReady[j].IrrAsSet != "" {
					asMacro = exchanges[i].PeersReady[j].IrrAsSet
				} else {
					asMacro = "AS" + exchanges[i].PeersReady[j].ASN
				}

				if exchanges[i].PeersReady[j].Ipv4Enabled {
					exchanges[i].PeersReady[j].PrefixFilters, err = bgpWorker.GenPrefixList(exchanges[i].PeersReady[j].PrefixList, asMacro, 4)
					if err != nil {
						log.Println(err)
					}
				}

				if exchanges[i].PeersReady[j].Ipv6Enabled {
					exchanges[i].PeersReady[j].PrefixFilters6, err = bgpWorker.GenPrefixList(exchanges[i].PeersReady[j].PrefixList6, asMacro, 6)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}()
	}
	wg.Wait()
	return exchanges
}
