package ixtypes

import "net"

/* This structure will be passed to the normal text template */
type ExchangePeer struct {
	ASN             string
	Ipv6Enabled     bool
	Ipv4Enabled     bool
	PrefixFilter    bool
	Active          bool
	Ipv4Addr        net.IP
	Ipv6Addr        net.IP
	LocalPreference int
	InfoPrefixes6   int64
	InfoPrefixes4   int64
	Group           string
	GroupEnabled    bool
	Group6          string
	Group6Enabled   bool
	IsRs            bool
	IsRsPeer        bool
}

type ExchangeOption string

type PeerTemplate struct {
	Peers           []ExchangePeer
	ExchangeOptions ExchangeOptions
}

/* Ix, is the whole definition of the exchange
read by the ini-configuration, and then completed
by filling PeersReady. This type will be exported
to the template function.

*/
type Ix struct {
	AdditionalConfig []string
	IxName           string
	Options          ExchangeOptions
	PeeringGroups    PeeringGroups
	PeersINI         ExchangePeers
	PeersReady       []ExchangePeer
	RouteServerReady []ExchangePeer
}

type IXs []Ix

type PeeringGroups map[string]bool
type ExchangePeers map[string]map[string][]ExchangePeer
type ExchangeOptions map[string]map[string]ExchangeOption
