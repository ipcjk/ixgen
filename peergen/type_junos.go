package peergen

import "net"

type junosDataString struct {
	Data string `json:"data"`
}

type junosDataIP struct {
	Data net.IP `json:"data"`
}

type junosDataInt64String struct {
	Data string `json:"data"`
}

type junosNeighbor struct {
	Family junosFamily            `json:"family"`
	Name   junosDataIP            `json:"name"`
	PeerAs []junosDataInt64String `json:"peer-as"`
}

type junosMaximumLimit struct {
	junosDataInt64String
}

type junosPrefixLimit struct {
	Maximum []junosMaximumLimit `json:"maximum"`
}

type junosLabeledUnicast struct {
	PrefixLimit []junosPrefixLimit `json:"prefix-limit"`
}

type junosFamilyInet6 struct {
	Inet6Unicast []junosLabeledUnicast `json:"labeled-unicast"`
}

type junosFamilyInet4 struct {
	InetUnicast []junosLabeledUnicast `json:"labeled-unicast"`
}

type junosFamily []struct {
	Inet6 []junosFamilyInet6 `json:"inet6"`
	Inet  []junosFamilyInet4 `json:"inet"`
}

type junosGroup struct {
	Name     junosDataString             `json:"name"`
	Neighbor []junosNeighbor             `json:"neighbor"`
	Type     []struct{ junosDataString } `json:"type"`
}

type junosAttributes struct {
	JunosChangedLocaltime string `json:"junos:changed-localtime"`
	JunosChangedSeconds   int64  `json:"junos:changed-seconds,string"`
}

type junosBgpGroup struct {
	Group []junosGroup `json:"group"`
}

type junosBGPProtocol struct {
	Bgp []junosBgpGroup `json:"bgp"`
}

type junosConfiguration struct {
	Attributes junosAttributes    `json:"attributes"`
	Protocols  []junosBGPProtocol `json:"protocols"`
}

type junOsJSON struct {
	Configuration []junosConfiguration `json:"configuration"`
}
