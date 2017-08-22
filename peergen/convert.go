package peergen

import (
	"encoding/json"
	"fmt"
	"github.com/ipcjk/ixgen/ixtypes"
	"io"
	"log"
	"strconv"
	"time"
)

// Functions that convert our Ix to different JSON / XML - configurations
// ConvertIxToJuniperJSON -> Juniper JSON
// ConvertIxToBrocadeSlxJSON -> Brocade SLX

func (p *Peergen) ConvertIxToJson(ixs ixtypes.IXs, w io.Writer) {
	res, err := json.Marshal(ixs)
	if err != nil {
		log.Fatalf("Cant decode IX into native format: err")
	}
	fmt.Fprint(w, string(res))
}

func (p *Peergen) ConvertIxToJsonPretty(ixs ixtypes.IXs, w io.Writer) {
	res, err := json.MarshalIndent(ixs, "", "  ")
	if err != nil {
		log.Fatalf("Cant decode IX into native format: err")
	}
	fmt.Fprint(w, string(res))
}

func (p *Peergen) ConvertIxToBrocadeSlxJSON(ixs ixtypes.IXs, w io.Writer) {
	fmt.Fprint(w, "Not done yet")
}

/* The function ConvertIxToJuniperJSON converts native format to a JUNOS
compatible JSON-configuration.
FIXME NEEDS support for the prefix list
*/
func (p *Peergen) ConvertIxToJuniperJSON(ixs ixtypes.IXs, w io.Writer, pretty bool) {
	var res []byte
	var err error
	var junosConfiguration = junOsJSON{
		[]junosConfiguration{
			{
				Attributes: junosAttributes{},
				Protocols: []junosBGPProtocol{
					{
						Bgp: []junosBgpGroup{
							{Group: []junosGroup{}},
						},
					},
				},
				PolicyOptions: []junosPolicyOptions{{}},
			},
		},
	}

	/* Attention, peering to JSON only works if a PEERGROUP is there
	 */
	for _, ix := range ixs {
		for k := range ix.PeeringGroups {
			var junosPeerConfiguration junosGroup
			for i := range ix.PeersReady {
				if ix.PeersReady[i].Group == k {
					junosPeerConfiguration.Name = junosDataString{k}
					junosPeerConfiguration.Type = []struct{ junosDataString }{
						{junosDataString{"external"}},
					}
					if ix.PeersReady[i].Ipv4Enabled {
						junosPeerConfiguration.Neighbor = append(junosPeerConfiguration.Neighbor,
							junosNeighbor{
								Family: junosFamily{
									{Inet6: []junosFamilyInet6{},
										Inet: []junosFamilyInet4{
											{InetUnicast: []junosLabeledUnicast{
												{
													PrefixLimit: []junosPrefixLimit{
														{Maximum: []junosMaximumLimit{
															{junosDataInt64String: junosDataInt64String{Data: strconv.FormatInt(ix.PeersReady[i].InfoPrefixes4, 10)}},
														}},
													},
												},
											}},
										}},
								},
								Name:   junosDataIP{Data: ix.PeersReady[i].Ipv4Addr},
								PeerAs: []junosDataInt64String{{Data: ix.PeersReady[i].ASN}},
							})
						if ix.PeersReady[i].PrefixFilterEnabled {
							/* Prefix or Router-Filter */

							if ix.PeersReady[i].PrefixAggregateMax {
								junosPolicyStatement := junosPolicyStatement{}
								junosPolicyStatement.Name = junosDataString{ix.PeersReady[i].PrefixFilters.Name}
								junosPolicyStatement.From = []junosPolicyFrom{{}}
								junosPolicyStatement.Then = []junosPolicyThen{}

								junosPeerConfiguration.Neighbor[0].Import = []junosDataString{{Data: ix.PeersReady[i].PrefixFilters.Name}}

								for _, PrefixRule := range ix.PeersReady[i].PrefixFilters.PrefixRules {
									if PrefixRule.Exact {
										junosPolicyStatement.From[0].RouteFilter = append(junosPolicyStatement.From[0].RouteFilter,
											junosRouteFilter{Address: junosDataString{PrefixRule.Prefix}, Exact: []junosDataString{{""}}})
									} else if PrefixRule.GreaterEqual != 0 {
										var data = fmt.Sprintf("/%d-/%d", PrefixRule.GreaterEqual, PrefixRule.LessEqual)
										junosPolicyStatement.From[0].RouteFilter = append(junosPolicyStatement.From[0].RouteFilter,
											junosRouteFilter{Address: junosDataString{PrefixRule.Prefix}, PrefixLengthRange: &junosDataString{data}})
									} else if PrefixRule.LessEqual != 0 {
										var data = fmt.Sprintf("/%d", PrefixRule.LessEqual)
										junosPolicyStatement.From[0].RouteFilter = append(junosPolicyStatement.From[0].RouteFilter,
											junosRouteFilter{Address: junosDataString{PrefixRule.Prefix}, UpTo: &junosDataString{data}})
									}
								}
								junosPolicyStatement.Then = append(junosPolicyStatement.Then, junosPolicyThen{Accept: []*junosDataString{{}}})
								junosConfiguration.Configuration[0].PolicyOptions[0].PolicyStatement =
									append(junosConfiguration.Configuration[0].PolicyOptions[0].PolicyStatement, junosPolicyStatement)

							} else {

								var junosPrefixList junosPrefixList
								junosPrefixList.Name = junosDataString{ix.PeersReady[i].PrefixFilters.Name}
								for _, PrefixRule := range ix.PeersReady[i].PrefixFilters.PrefixRules {
									junosPrefixList.PrefixListItem = append(junosPrefixList.PrefixListItem, junosPrefixListItem{
										Name: junosDataString{Data: PrefixRule.Prefix},
									})
								}

								junosConfiguration.Configuration[0].PolicyOptions[0].PrefixList =
									append(junosConfiguration.Configuration[0].PolicyOptions[0].PrefixList, junosPrefixList)
							}
						}
					}
					if ix.PeersReady[i].Ipv6Enabled {
						junosPeerConfiguration.Neighbor = append(junosPeerConfiguration.Neighbor,
							junosNeighbor{
								Family: junosFamily{
									{Inet6: []junosFamilyInet6{{Inet6Unicast: []junosLabeledUnicast{
										{
											PrefixLimit: []junosPrefixLimit{
												{Maximum: []junosMaximumLimit{
													{junosDataInt64String: junosDataInt64String{Data: strconv.FormatInt(ix.PeersReady[i].InfoPrefixes6, 10)}},
												}},
											},
										},
									}}},
										Inet: []junosFamilyInet4{}},
								},
								Name:   junosDataIP{Data: ix.PeersReady[i].Ipv6Addr},
								PeerAs: []junosDataInt64String{{Data: ix.PeersReady[i].ASN}},
							})
						if ix.PeersReady[i].PrefixFilterEnabled {
							if ix.PeersReady[i].PrefixAggregateMax {
								junosPeerConfiguration.Neighbor[1].Import = []junosDataString{{Data: ix.PeersReady[i].PrefixFilters6.Name}}
								junosPolicyStatement := junosPolicyStatement{}
								junosPolicyStatement.Name = junosDataString{ix.PeersReady[i].PrefixFilters.Name}
								junosPolicyStatement.From = []junosPolicyFrom{{}}
								junosPolicyStatement.Then = []junosPolicyThen{}

								for _, PrefixRule6 := range ix.PeersReady[i].PrefixFilters6.PrefixRules {
									if PrefixRule6.Exact {
										junosPolicyStatement.From[0].RouteFilter = append(junosPolicyStatement.From[0].RouteFilter,
											junosRouteFilter{Address: junosDataString{PrefixRule6.Prefix}, Exact: []junosDataString{{""}}})
									} else if PrefixRule6.GreaterEqual != 0 {
										var data = fmt.Sprintf("/%d-/%d", PrefixRule6.GreaterEqual, PrefixRule6.LessEqual)
										junosPolicyStatement.From[0].RouteFilter = append(junosPolicyStatement.From[0].RouteFilter,
											junosRouteFilter{Address: junosDataString{PrefixRule6.Prefix}, PrefixLengthRange: &junosDataString{data}})
									} else if PrefixRule6.LessEqual != 0 {
										var data = fmt.Sprintf("/%d", PrefixRule6.LessEqual)
										junosPolicyStatement.From[0].RouteFilter = append(junosPolicyStatement.From[0].RouteFilter,
											junosRouteFilter{Address: junosDataString{PrefixRule6.Prefix}, UpTo: &junosDataString{data}})
									}
								}
								junosPolicyStatement.Then = append(junosPolicyStatement.Then, junosPolicyThen{Accept: []*junosDataString{{}}})

								junosConfiguration.Configuration[0].PolicyOptions[0].PolicyStatement =
									append(junosConfiguration.Configuration[0].PolicyOptions[0].PolicyStatement, junosPolicyStatement)

							} else {

								var junosPrefixList junosPrefixList
								junosPrefixList.Name = junosDataString{ix.PeersReady[i].PrefixFilters6.Name}
								for _, PrefixRule := range ix.PeersReady[i].PrefixFilters6.PrefixRules {
									junosPrefixList.PrefixListItem = append(junosPrefixList.PrefixListItem, junosPrefixListItem{
										Name: junosDataString{Data: PrefixRule.Prefix},
									})
								}
								junosConfiguration.Configuration[0].PolicyOptions[0].PrefixList =
									append(junosConfiguration.Configuration[0].PolicyOptions[0].PrefixList, junosPrefixList)
							}
						}
					}
					if len(junosPeerConfiguration.Neighbor) > 0 {
						junosConfiguration.Configuration[0].Protocols[0].Bgp[0].Group = append(
							junosConfiguration.Configuration[0].Protocols[0].Bgp[0].Group,
							junosPeerConfiguration)
					}
				}
			}
		}
	}

	/* Add timestamp */
	t := time.Now()
	junosConfiguration.Configuration[0].Attributes.JunosChangedSeconds = t.Unix()
	junosConfiguration.Configuration[0].Attributes.JunosChangedLocaltime = t.String()

	if pretty {
		res, err = json.MarshalIndent(junosConfiguration, "", "  ")
	} else {
		res, err = json.Marshal(junosConfiguration)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(res))
}
