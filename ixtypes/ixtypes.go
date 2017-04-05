package ixtypes

import "net"

/* This structure will be passed to the normal text template */
type ExchangePeer struct {
	Active          bool   `json:"active"`
	ASN             string `json:"asn"`
	Group           string `json:"group"`
	Group6          string `json:"group6"`
	GroupEnabled    bool   `json:"groupenabled"`
	Group6Enabled   bool   `json:"group6_enabled"`
	InfoPrefixes4   int64  `json:"infoprefixes4"`
	InfoPrefixes6   int64  `json:"infoprefixes6"`
	Ipv4Addr        net.IP `json:"ipv4addr"`
	Ipv6Addr        net.IP `json:"ipv6addr"`
	Ipv4Enabled     bool   `json:"ipv4enabled"`
	Ipv6Enabled     bool   `json:"ipv6enabled"`
	IrrAsSet        string `json:"irrasset"`
	IsRs            bool   `json:"isrs"`
	IsRsPeer        bool   `json:"isrsper"`
	LocalPreference int    `json:"localpreference"`
	PrefixFilter    bool   `json:"prefixfilter"`
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
	AdditionalConfig []string        `json:"additionalconfig"`
	IxName           string          `json:"ixname"`
	Options          ExchangeOptions `json:"options"`
	PeeringGroups    PeeringGroups   `json:"peeringgroups"`
	PeersINI         ExchangePeers   `json:"peers_configured"`
	PeersReady       []ExchangePeer  `json:"peersready"`
	RouteServerReady []ExchangePeer  `json:"routeserverready"`
}

type IXs []Ix

type PeeringGroups map[string]bool
type ExchangePeers map[string]map[string][]ExchangePeer
type ExchangeOptions map[string]map[string]ExchangeOption
