package ixworkers

import (
	"github.com/ipcjk/ixgen/bgpq3workers"
	"github.com/ipcjk/ixgen/ixtypes"
	"github.com/ipcjk/ixgen/peeringdb"
	"log"
	"strconv"
	"sync"
)

/* helper function to convert string like 1 or 0 to a boolean */
func isTrue(value string) bool {
	if value == "true" || value == "1" {
		return true
	}
	return false
}

func WorkerMergePeerConfiguration(exchanges ixtypes.IXs, apiServiceURL string, apiKey string, exchangeOnly string, myASN int64, prefixFactor float64) ixtypes.IXs {
	var wg sync.WaitGroup
	peerDB := peeringdb.Peeringdb(apiServiceURL, apiKey)

	/* thread shared lists of asked ASNs, save resources */
	var lockASNList = sync.Mutex{}
	var ASNList = make(map[int64]peeringdb.NetData)

	wg.Add(len(exchanges))
	for k := range exchanges {
		var i = k
		go func() {
			var asnQueryList []string
			var asnQuery = make(map[string]bool)
			var rsAuto bool
			defer wg.Done()
			if exchangeOnly != "" && exchangeOnly != exchanges[i].IxName {
				return
			}
			/* check if routeserver options is set and also if its true or false */
			rsnEnabled, rsnExist := exchanges[i].Options[exchanges[i].IxName]["routeserver"]
			if rsnExist {
				if isTrue(string(rsnEnabled)) {
					rsAuto = true
				}
			}
			rsnASN, rsnOk := exchanges[i].Options[exchanges[i].IxName]["rs_asn"]

			/* check for pinned ixID */
			ixId, ixOk := exchanges[i].Options[exchanges[i].IxName]["ixid"]
			myPeers, err := peerDB.GetPeersOnIX(exchanges[i].IxName, string(ixId), ixOk)
			if err != nil {
				log.Printf("Cant get Peers for IX: %s, error: %s\n", exchanges[i].IxName, err)
				return
			}

			/* peering DB api change, pull asns into a local query list,
			respect other threads if necessary */
			for _, peer := range myPeers.Data {
				/* Check if we or some other thread already know this peer,
				push it to our query list, if we need, else ignore it
				*/
				lockASNList.Lock()
				if _, peerOk := ASNList[peer.Asn]; !peerOk {
					/* We have not seen the peer yet, lets save it to our query list */
					asnQuery[strconv.FormatInt(peer.Asn, 10)] = true
				}
				lockASNList.Unlock()
			}

			/* map to slice */
			for s := range asnQuery {
				asnQueryList = append(asnQueryList, s)
			}

			/* ask for information about every ASN in our local list, save the information
			into the shared thread list */
			for {
				/* limit myself to 150 asns at maximi, */
				mincut := 150

				/* nothing left to do? fine, break the loop */
				if len(asnQueryList) == 0 {
					break
				}

				/* something still to do? fine */
				if len(asnQueryList) < mincut {
					mincut = len(asnQueryList)
				}

				/* reshape our view on the slices */
				questionASN := asnQueryList[:mincut]
				asnQueryList = asnQueryList[mincut:]

				/* send our wishlist to the peering db service */
				data, err := peerDB.GetASNsbyList(questionASN)
				if err != nil {
					log.Println(err)
				}

				/* check if we received the same amount of ASNs that we asked for */
				if len(questionASN) != len(data.Data) {
					log.Printf("Attention, mismatch between askedASNs (%d) and receivedASNs (%d)\n",
						len(questionASN), len(data.Data))
				}

				/* copy data into our ASN map */
				for index := range data.Data {
					lockASNList.Lock()
					ASNList[data.Data[index].Asn] = data.Data[index]
					lockASNList.Unlock()
				}
			}

			for _, peer := range myPeers.Data {
				peerASN := strconv.FormatInt(peer.Asn, 10)

				if peer.Asn == myASN {
					continue
				}

				if _, exists := ASNList[peer.Asn]; !exists {
					log.Printf("Error pulling ASN for peer: %d\n", peer.Asn)
					continue
				}

				peerDbNetwork := ASNList[peer.Asn]
				if err != nil {
					log.Printf("Error pulling ASN for peer: %d, error: %s", peer.Asn, err)
					continue
				}

				if peerDbNetwork.InfoType == "Route Server" {
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
							IrrAsSet: peerDbNetwork.IrrAsSet,
						}
					if rsPeer.Ipv6Addr.To16() != nil {
						rsPeer.Ipv6Enabled = true
					}
					if rsPeer.Ipv4Addr.To4() != nil {
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

					rsPeer.InfoPrefixes4 = peerDbNetwork.InfoPrefixes4
					rsPeer.InfoPrefixes6 = peerDbNetwork.InfoPrefixes6

					if rprefixOk {
						rsPeer.InfoPrefixes4, _ = strconv.ParseInt(string(infoprefixes4), 10, 64)
					}

					if rprefix6Ok {
						rsPeer.InfoPrefixes6, _ = strconv.ParseInt(string(infoprefixes6), 10, 64)
					}

					/* take care of prefix factor if given */
					if prefixFactor != 1.0 && rsPeer.InfoPrefixes4 != 0 {
						rsPeer.InfoPrefixes4 = int64(prefixFactor * float64(rsPeer.InfoPrefixes4))
					}
					/* take care of prefix factor if given */
					if prefixFactor != 1.0 && rsPeer.InfoPrefixes6 != 0 {
						rsPeer.InfoPrefixes6 = int64(prefixFactor * float64(rsPeer.InfoPrefixes6))
					}

					if rsAuto && rsnOk && peerASN != string(rsnASN) {
						log.Printf("Ignoring route-server advertised from ASN %s, but IX ASN shall be %s\n", peerASN, rsnASN)
					} else if rsAuto {
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
						confPeer.InfoPrefixes4 = peerDbNetwork.InfoPrefixes4
					}
					if confPeer.InfoPrefixes6 == 0 {
						confPeer.InfoPrefixes6 = peerDbNetwork.InfoPrefixes6
					}
					// Valid Ipv6? Then, don't load from peerDB
					if confPeer.Ipv6Addr.To16() == nil {
						confPeer.Ipv6Enabled = false
						// if peerDB contains a valid IPv6 address, take the struct
						// and enable the output
						if peer.Ipaddr6.To16() != nil {
							confPeer.Ipv6Addr = peer.Ipaddr6
							confPeer.Ipv6Enabled = true
						}
					}
					// Valid Ipv4? Then, don't load from peerDB
					if confPeer.Ipv4Addr.To4() == nil {
						confPeer.Ipv4Enabled = false
						// if peerDB contains a valid IPv4 address, take the struct
						// and enable the output
						if peer.Ipaddr4.To4() != nil {
							confPeer.Ipv4Addr = peer.Ipaddr4
							confPeer.Ipv4Enabled = true
						}
					}

					/* take care of prefix factor if given */
					if prefixFactor != 1.0 && confPeer.InfoPrefixes4 != 0 {
						confPeer.InfoPrefixes4 = int64(prefixFactor * float64(confPeer.InfoPrefixes4))
					}
					/* take care of prefix factor if given */
					if prefixFactor != 1.0 && confPeer.InfoPrefixes6 != 0 {
						confPeer.InfoPrefixes6 = int64(prefixFactor * float64(confPeer.InfoPrefixes6))
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
					if confPeer.IrrAsSet == "" {
						confPeer.IrrAsSet = peerDbNetwork.IrrAsSet
					}
					exchanges[i].PeersReady = append(exchanges[i].PeersReady, confPeer)
				} else if exchanges[i].Options[exchanges[i].IxName]["wildcard"] == "1" {

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

					if exchanges[i].Options[exchanges[i].IxName]["wildcard_prefix_filter"] == "1" {
						wildPeer.PrefixFilterEnabled = true
						wildPeer.PrefixList = "wildcard" + peerASN + "-4"
						wildPeer.PrefixList6 = "wildcard" + peerASN + "-6"
					}

					if wildPeer.Ipv6Addr.To16() != nil {
						wildPeer.Ipv6Enabled = true
					}
					if wildPeer.Ipv4Addr.To4() != nil {
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

					wildPeer.IrrAsSet = peerDbNetwork.IrrAsSet
					wildPeer.InfoPrefixes4 = peerDbNetwork.InfoPrefixes4
					wildPeer.InfoPrefixes6 = peerDbNetwork.InfoPrefixes6
					exchanges[i].PeersReady = append(exchanges[i].PeersReady, wildPeer)

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

				if exchanges[i].PeersReady[j].IsRs || exchanges[i].PeersReady[j].PrefixFilterEnabled == false {
					continue
				}

				if exchanges[i].PeersReady[j].IrrAsSet != "" {
					asMacro = exchanges[i].PeersReady[j].IrrAsSet
				} else {
					asMacro = "AS" + exchanges[i].PeersReady[j].ASN
				}

				if exchanges[i].PeersReady[j].Ipv4Enabled {
					exchanges[i].PeersReady[j].PrefixFilters, err = bgpWorker.GenPrefixList(exchanges[i].PeersReady[j].PrefixList, asMacro, 4, exchanges[i].PeersReady[j].PrefixAggregateMax)
					if err != nil {
						log.Println(err)
					}
				}

				if exchanges[i].PeersReady[j].Ipv6Enabled {
					exchanges[i].PeersReady[j].PrefixFilters6, err = bgpWorker.GenPrefixList(exchanges[i].PeersReady[j].PrefixList6, asMacro, 6, exchanges[i].PeersReady[j].PrefixAggregateMax)
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
