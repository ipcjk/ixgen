package peeringdb

import (
	"net"
	"time"
)

// IxLanData is internet-exchange-peering network
// we need to decode the json dataset into a golang struct
type IxLanData struct {
	ID           int64       `json:"id"`
	ArpSponge    interface{} `json:"arp_sponge"`
	Created      time.Time   `json:"created"`
	Descr        string      `json:"descr"`
	Dot1qSupport bool        `json:"dot1q_support"`
	IxID         int64       `json:"ix_id"`
	Mtu          int64       `json:"mtu"`
	Name         string      `json:"name"`
	RsAsn        int64       `json:"rs_asn"`
	Status       string      `json:"status"`
	Updated      time.Time   `json:"updated"`
}

// OrgData is the struct, that is used for decoding the organization object
type OrgData struct {
	ID       int64     `json:"id"`
	Address1 string    `json:"address1"`
	Address2 string    `json:"address2"`
	City     string    `json:"city"`
	Country  string    `json:"country"`
	Created  time.Time `json:"created"`
	Name     string    `json:"name"`
	Notes    string    `json:"notes"`
	State    string    `json:"state"`
	Status   string    `json:"status"`
	Updated  time.Time `json:"updated"`
	Website  string    `json:"website"`
	Zipcode  int64     `json:"zipcode,string"`
}

// FacData is the struct, that is used for decoding the faculty object
type FacData struct {
	ID       int64     `json:"id"`
	Address1 string    `json:"address1"`
	Address2 string    `json:"address2"`
	City     string    `json:"city"`
	Clli     string    `json:"clli"`
	Country  string    `json:"country"`
	Created  time.Time `json:"created"`
	Name     string    `json:"name"`
	NetCount int64     `json:"net_count"`
	Notes    string    `json:"notes"`
	Npanxx   string    `json:"npanxx"`
	OrgID    int64     `json:"org_id"`
	OrgName  string    `json:"org_name"`
	Rencode  string    `json:"rencode"`
	State    string    `json:"state"`
	Status   string    `json:"status"`
	Updated  time.Time `json:"updated"`
	Website  string    `json:"website"`
	Zipcode  string    `json:"zipcode"`
}

// NetixlanData is the struct,
// that is used for decoding the peering network
type NetixlanData struct {
	ID       int64     `json:"id"`
	Asn      int64     `json:"asn"`
	Created  time.Time `json:"created"`
	Ipaddr4  net.IP    `json:"ipaddr4"`
	Ipaddr6  net.IP    `json:"ipaddr6"`
	IsRsPeer bool      `json:"is_rs_peer"`
	IxID     int64     `json:"ix_id"`
	IxlanID  int64     `json:"ixlan_id"`
	Name     string    `json:"name"`
	NetID    int64     `json:"net_id"`
	Notes    string    `json:"notes"`
	Speed    int64     `json:"speed"`
	Status   string    `json:"status"`
	Updated  time.Time `json:"updated"`
}

// netsIxLanDataSortedByASN is the slice of NetixlanData
// that can be used to sort NetixlanData by ASN
type netsIxLanDataSortedByASN []NetixlanData

func (a netsIxLanDataSortedByASN) Len() int           { return len(a) }
func (a netsIxLanDataSortedByASN) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a netsIxLanDataSortedByASN) Less(i, j int) bool { return a[i].Asn < a[j].Asn }

// NetData is the struct,
// that is used for decoding the ASN object
type NetData struct {
	ID              int64     `json:"id"`
	Aka             string    `json:"aka"`
	Asn             int64     `json:"asn"`
	Created         time.Time `json:"created"`
	InfoIpv6        bool      `json:"info_ipv6"`
	InfoMulticast   bool      `json:"info_multicast"`
	InfoPrefixes4   int64     `json:"info_prefixes4"`
	InfoPrefixes6   int64     `json:"info_prefixes6"`
	InfoRatio       string    `json:"info_ratio"`
	InfoScope       string    `json:"info_scope"`
	InfoTraffic     string    `json:"info_traffic"`
	InfoType        string    `json:"info_type"`
	InfoUnicast     bool      `json:"info_unicast"`
	IrrAsSet        string    `json:"irr_as_set"`
	LookingGlass    string    `json:"looking_glass"`
	Name            string    `json:"name"`
	Notes           string    `json:"notes"`
	OrgID           int64     `json:"org_id"`
	PolicyContracts string    `json:"policy_contracts"`
	PolicyGeneral   string    `json:"policy_general"`
	PolicyLocations string    `json:"policy_locations"`
	PolicyRatio     bool      `json:"policy_ratio"`
	PolicyURL       string    `json:"policy_url"`
	RouteServer     string    `json:"route_server"`
	Status          string    `json:"status"`
	Updated         time.Time `json:"updated"`
	Website         string    `json:"website"`
}

// IxData is the struct,
// that is used for decoding the internet exchange object
type IxData struct {
	ID              int64       `json:"id"`
	City            string      `json:"city"`
	Country         string      `json:"country"`
	Created         time.Time   `json:"created"`
	FacSet          []FacData   `json:"fac_set"`
	IxlanSet        []IxLanData `json:"ixlan_set"`
	Media           string      `json:"media"`
	Name            string      `json:"name"`
	NameLong        string      `json:"name_long"`
	Notes           string      `json:"notes"`
	OrgData         `json:"org"`
	OrgID           int64     `json:"org_id"`
	PolicyEmail     string    `json:"policy_email"`
	PolicyPhone     string    `json:"policy_phone"`
	ProtoIpv6       bool      `json:"proto_ipv6"`
	ProtoMulticast  bool      `json:"proto_multicast"`
	ProtoUnicast    bool      `json:"proto_unicast"`
	RegionContinent string    `json:"region_continent"`
	Status          string    `json:"status"`
	TechEmail       string    `json:"tech_email"`
	TechPhone       string    `json:"tech_phone"`
	Updated         time.Time `json:"updated"`
	URLStats        string    `json:"url_stats"`
	Website         string    `json:"website"`
}

// Ix is the combined array "meta"-object
// for the internet exchange
type Ix struct {
	Data []IxData `json:"data"`
	Meta struct {
	} `json:"meta"`
}

// Netixlan is the combined array "meta"-object
// for the peering networks
type Netixlan struct {
	Data []NetixlanData `json:"data"`
	Meta struct {
	} `json:"meta"`
}

// Fac is the combined array "meta"-object for the faculties
type Fac struct {
	Data []FacData `json:"data"`
	Meta struct {
		Generated float64 `json:"generated"`
	} `json:"meta"`
}

// Org is the combined array "meta"-object
// for the organizations
type Org struct {
	Data []OrgData `json:"data"`
	Meta struct {
		Generated float64 `json:"generated"`
	} `json:"meta"`
}

// IxLAN is the combined array "meta"-object
// for the networks that an IX offer
type IxLAN struct {
	Data []IxLanData `json:"data"`
	Meta struct {
	} `json:"meta"`
}

// Net is the combined array "meta"-object
// for the networks (ASN)
type Net struct {
	Data []NetData `json:"data"`
	Meta struct {
	} `json:"meta"`
}
