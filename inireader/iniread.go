package inireader

import (
	"bufio"
	"fmt"
	"github.com/ipcjk/ixgen/ixtypes"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var splitBy = `\s+`
var splitReg = regexp.MustCompile(splitBy)

var PossibleOptions = map[string]bool{
	"ixid":                   true,
	"routeserver":            true,
	"routeserver_group":      true,
	"peer_group":             true,
	"routeserver_group6":     true,
	"peer_group6":            true,
	"wildcard":               true,
	"importpolicy":           true,
	"exportpolicy":           true,
	"routeserver_prefixes":   true,
	"routeserver_prefixes6":  true,
	"rs_asn":                 true,
	"wildcard_prefix_filter": true,
}

const (
	ixsection = iota
	options
	peers
	header
)

func ReadPeeringConfig(r io.Reader) ixtypes.IXs {
	var IXs ixtypes.IXs
	var ix *ixtypes.Ix

	var whichSection = ixsection
	reader := bufio.NewReader(r)
	var currentHead string
	var lineNum int

	for {
		line, err := reader.ReadString('\n')
		lineNum++

		if line == "\n" {
			continue
		}

		line = strings.Replace(line, "\n", "", 1)
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") {
			continue
		} else if strings.HasPrefix(line, "[peers]") {
			whichSection = peers
		} else if strings.HasPrefix(line, "[options]") {
			whichSection = options
		} else if strings.HasPrefix(line, "[additionalConfig]") {
			whichSection = header
		} else if strings.HasPrefix(line, "[") {
			whichSection = ixsection

			/* Add if there was an open IX ahead , else
			allocate a new object*/
			if ix != nil {
				IXs = append(IXs, *ix)
			}

			ix = new(ixtypes.Ix)

			currentHead = strings.Replace(line, "[", "", 1)
			currentHead = strings.Replace(currentHead, "]", "", 1)

			ix.PeersINI = make(ixtypes.ExchangePeers)
			ix.Options = make(ixtypes.ExchangeOptions)
			ix.PeeringGroups = make(ixtypes.PeeringGroups)
			ix.IxName = currentHead

			_, ok := ix.PeersINI[currentHead]
			/* Complicated, we need to init the Exchange map if it does not exist so far */
			if !ok {
				ix.PeersINI[currentHead] = make(map[string][]ixtypes.ExchangePeer, 32)
			}
		} else if line != "" && whichSection == options {
			ParseOptionLine(line, ix.Options, currentHead)
			if ix.Options[currentHead]["routeserver_group"] != "" {
				ix.PeeringGroups[string(ix.Options[currentHead]["routeserver_group"])] = true
			}
			if ix.Options[currentHead]["routeserver_group6"] != "" {
				ix.PeeringGroups[string(ix.Options[currentHead]["routeserver_group6"])] = true
			}
			if ix.Options[currentHead]["peer_group6"] != "" {
				ix.PeeringGroups[string(ix.Options[currentHead]["peer_group6"])] = true
			}
			if ix.Options[currentHead]["peer_group"] != "" {
				ix.PeeringGroups[string(ix.Options[currentHead]["peer_group"])] = true
			}
		} else if line != "" && whichSection == peers {
			Peer := ParsePeerLine(line, lineNum)
			if Peer.Active == true {
				ix.PeersINI[currentHead][Peer.ASN] = append(ix.PeersINI[currentHead][Peer.ASN], Peer)
				if Peer.Group6 != "" {
					ix.PeeringGroups[Peer.Group6] = true
				}
				if Peer.Group != "" {
					ix.PeeringGroups[Peer.Group] = true
				}
			}
		} else if line != "" && whichSection == header {
			ix.AdditionalConfig = append(ix.AdditionalConfig, line)
		}

		if err == io.EOF {
			if ix != nil {
				IXs = append(IXs, *ix)
			}
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	return IXs
}

func ParseOptionLine(line string, exchangeOptions ixtypes.ExchangeOptions, currentHead string) {
	_, ok := exchangeOptions[currentHead]
	if !ok {
		exchangeOptions[currentHead] = make(map[string]ixtypes.ExchangeOption, 10)
	}

	for key := range PossibleOptions {
		if strings.HasPrefix(line, key+"=") {
			value := strings.Replace(line, key+"=", "", 1)
			if PossibleOptions[key] && value != "" {
				exchangeOptions[currentHead][key] = ixtypes.ExchangeOption(value)
			}
		}
	}
}

func ParsePeerLine(line string, lineNumber int) ixtypes.ExchangePeer {
	var Peer = ixtypes.ExchangePeer{Active: true, Ipv4Enabled: true, Ipv6Enabled: true,
		PrefixFilterEnabled: false, GroupEnabled: true, Group6Enabled: true, Unconfigured: false}
	var err error

	for index, value := range splitReg.Split(line, -1) {
		if index == 0 {
			_, err := strconv.Atoi(value)
			if err != nil {
				log.Printf("Ignoring line %d", lineNumber)
				break
			}
			Peer.ASN = value
			continue
		}
		value := strings.Replace(value, "\"", "", -1)
		if value == "ipv4=0" {
			Peer.Ipv4Enabled = false
		} else if value == "ipv6=0" {
			Peer.Ipv6Enabled = false
		} else if value == "ipv4=1" {
			Peer.Ipv4Enabled = true
		} else if value == "ipv6=1" {
			Peer.Ipv6Enabled = true
		} else if value == "active=0" {
			Peer.Active = false
		} else if value == "active=1" {
			Peer.Active = true
		} else if value == "prefix_filter=1" {
			Peer.PrefixFilterEnabled = true
		} else if value == "prefix_filter=0" {
			Peer.PrefixFilterEnabled = false
		} else if value == "prefixfilter_aggregate=1" {
			Peer.PrefixAggregateMax = true
		} else if value == "prefixfilter_aggregate=0" {
			Peer.PrefixAggregateMax = false
		} else if value == "unconfigured=1" {
			Peer.Unconfigured = true
		} else if value == "unconfigured=0" {
			Peer.Unconfigured = false
		} else if strings.HasPrefix(value, "prefix_list=") {
			Peer.PrefixList = strings.TrimPrefix(value, "prefix_list=")
		} else if strings.HasPrefix(value, "prefix_list6=") {
			Peer.PrefixList6 = strings.TrimPrefix(value, "prefix_list6=")
		} else if strings.HasPrefix(value, "local_pref=") {
			localPref, err := strconv.Atoi(strings.TrimPrefix(value, "local_pref="))
			if err != nil {
				log.Printf("Peer %s has local_pref, but no value given", Peer.ASN)
			} else {
				Peer.LocalPreference = localPref
			}
		} else if strings.HasPrefix(value, "ipv4_addr=") {
			ipv4Addr := strings.TrimPrefix(value, "ipv4_addr=")
			fmt.Printf("Found fixed peering of %s", ipv4Addr)
		} else if strings.HasPrefix(value, "ipv6_addr=") {
			ipv6Addr := strings.TrimPrefix(value, "ipv6_addr=")
			fmt.Printf("Found fixed peering of %s", ipv6Addr)
		} else if strings.HasPrefix(value, "peer_group=") {
			Peer.Group = strings.TrimPrefix(value, "peer_group=")
		} else if strings.HasPrefix(value, "peer_group6=") {
			Peer.Group6 = strings.TrimPrefix(value, "peer_group6=")
		} else if strings.HasPrefix(value, "infoprefixes4=") {
			Peer.InfoPrefixes4, err = strconv.ParseInt(strings.TrimPrefix(value, "infoprefixes4="), 10, 64)
			if err != nil {
				log.Println("Wrong prefix limit for IPv4, ignoring")
			}
		} else if strings.HasPrefix(value, "infoprefixes6=") {
			Peer.InfoPrefixes6, err = strconv.ParseInt(strings.TrimPrefix(value, "infoprefixes6="), 10, 64)
			if err != nil {
				log.Println("Wrong prefix limit for IPv6, ignoring")
			}
		} else if value == "group=0" {
			Peer.GroupEnabled = false
		} else if value == "group6=0" {
			Peer.Group6Enabled = false
		} else if value == "group=1" {
			Peer.GroupEnabled = true
		} else if value == "group6=1" {
			Peer.Group6Enabled = true
		} else if strings.HasPrefix(value, "password4=") {
			Peer.Password4 = strings.TrimPrefix(value, "password4=")
		} else if strings.HasPrefix(value, "password6=") {
			Peer.Password6 = strings.TrimPrefix(value, "password6=")
		} else if strings.HasPrefix(value, "irrasset=") {
			Peer.IrrAsSet = strings.TrimPrefix(value, "irrasset=")
		} else {
			log.Printf("Unknown parameter %s for peer on line %d ", value, lineNumber)
		}
	}

	/* set prefixFilterNames automatically, if prefixfilter is enabled
	but name is not given
	*/
	if Peer.PrefixFilterEnabled && Peer.Ipv4Enabled && Peer.PrefixList == "" {
		Peer.PrefixList = fmt.Sprintf("as%s_filter_ip4_filter", Peer.ASN)
	}

	if Peer.PrefixFilterEnabled && Peer.Ipv6Enabled && Peer.PrefixList6 == "" {
		Peer.PrefixList6 = fmt.Sprintf("as%s_filter_ip6_filter", Peer.ASN)
	}

	return Peer
}
