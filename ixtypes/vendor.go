package ixtypes

import "net"

// JunOsJSON is the struct definition
// for decoding into JSON configuration style for Juniper
// will be used in future

type junOsJSON struct {
	Configuration []struct {
		Attributes struct {
			JunosChangedLocaltime string `json:"junos:changed-localtime"`
			JunosChangedSeconds   int64  `json:"junos:changed-seconds,string"`
		} `json:"attributes"`
		Protocols []struct {
			Bgp []struct {
				Group []struct {
					Name struct {
						Data string `json:"data"`
					} `json:"name"`
					Neighbor []struct {
						Family []struct {
							Inet6 []struct {
								LabeledUnicast []struct {
									PrefixLimit []struct {
										Maximum []struct {
											Data int64 `json:"data,string"`
										} `json:"maximum"`
									} `json:"prefix-limit"`
								} `json:"labeled-unicast"`
							} `json:"inet6"`
						} `json:"family"`
						Name struct {
							Data net.IP `json:"data"`
						} `json:"name"`
						PeerAs []struct {
							Data int64 `json:"data,string"`
						} `json:"peer-as"`
					} `json:"neighbor"`
				} `json:"group"`
			} `json:"bgp"`
		} `json:"protocols"`
	} `json:"configuration"`
}
