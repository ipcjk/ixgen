[![Docker Repository on Quay](https://quay.io/repository/ipcjk/ixgen/status "Docker Repository on Quay")](https://quay.io/repository/ipcjk/ixgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/ipcjk/ixgen)](https://goreportcard.com/report/github.com/ipcjk/ixgen)
[![Build Status](https://travis-ci.org/ipcjk/ixgen.svg?branch=master)](https://travis-ci.org/ipcjk/ixgen)

# ixgen 
Ixgen is yet-another open-source, multi-platform generator for peering configurations on IXs incorporating the global peeringdb api, but also is able to spin up its own "compatible" server for faster results. Ixgen is configured by an INI- or JSON-style format, producing custom template-driven or fixed json-style configurations, that can be printed on the terminal, to a file or served by HTTP. Direct access to routers REST-APIs and ssh/scp-upload is planned. 

Ixgen is shipped with cross-compiled executables for Darwin, Linux and Windows. Only Linux and Darwin currently support the prefixfilter generator.

# ixgen on docker
Run ixgen from docker with your peering configuration as volume parameter:

docker run -v /Users/joerg/peering.ini:/ixgen/release/configuration/peering.ini joerg/ixgen:latest

## how it works 
Ixgen works by querying the peeringdb-API or its own local API-service for peering members and specific network configurations and populate learned things with a custom router template. 

### ini-style configuration 
Ixgen is configured by _configuration/peering.ini_, that contains a list of Internet Exchanges  and peering as numbers, that the user want to configure on his router.

### peeringdb client and server
Ixgen has not only a peeringdb client module but also can start a limited fast peeringdb-lookalike server,  that can also queried from network. By defaults, it starts an embedded http api server for its own usage. 

### apiserver 
Also you can start-up a standalone version called apiserver. Apiserver can answer limited peeringdb-api-queries but also generate router configuration via HTTP-POST.

 
### flavor / templating 
Ixgen can use different templates for generating router configurations, by default Brocade and Juniper style command line syntax is shipped. The flavor is given on the command line or in the http query with the _-style_ argument, else its always Brocade MLXE (Netiron)  by default. You can create your own templates in the _templates_-directory. Please see the section _Default output and templates_ for more information. 
 
## Quickstart 

### Populating the cache 
Before the first usage, you may want to download the peeringdb start files, else we wont benefit from the  fast local cache. To populate the cache, you need to start-up _ixgen_ with the _-buildcache_ parameter.
 
    ixgen -buildcache

### peering.ini 
The _peering.ini_ is an easy configuration file, that is using
sectors and lists.

For every Internet Exchange that needs to be configured, a section head combination with the original record name of the Exchange and a possible network name , that can be both found on the peeringdb database or website, has to be specified.
 
Several subsections with general options for the Exchange configuration or the peering list can be added.

### Example configuration 
If you need to connect to the DE-CIX in Frankfurt you will add the "DE-CIX Franfurt" name and the network name [Main or Jumbo] separated with a _||_ as an section head in squared brackets.
    
    [DE-CIX Frankfurt||Main]
       
Now you can add a  subsection for general options or the peering list.
   
    [DE-CIX Frankfurt||Main]
    [peers]
    [options]
    
### Add peer to list 
If you want to peer with the AS196922-routers (hint: that is me!) on DE-CIX, you just need to add the right as numbers on a new line after the _[peers]_-section. 

    [DE-CIX Frankfurt||Main]
    [peers]
    196922
    [options]
    

### Run ixgen 

Starting ixgen with default options will now print out the peering bgp configuration for DE-CIX with the wished neighbor statements for all ipv4 and ipv6 routers. 

    ./ixgen.mac 
    
The call will print out my DECIX-configuration for Frankfurt: 

    router bgp
    neighbor 80.81.194.25 remote-as 196922
    neighbor 80.81.194.25 peer-group decix-peer
    address-family ipv4 unicast
    neighbor 80.81.194.25 maximum-prefix 64
    no neighbor 80.81.194.25 shutdown
    exit-address-family
    address-family ipv6 unicast
    neighbor 2001:7f8::3:13a:0:1 remote-as  196922
    neighbor 2001:7f8::3:13a:0:1 peer-group decix-peer6
    neighbor 2001:7f8::3:13a:0:1 maximum-prefix 10
    no neighbor 2001:7f8::3:13a:0:1 shutdown
    exit-address-family
    
 Or a call to output the native format in json:
   
     ./ixgen -style native/json_pretty
     {
     	"additionalconfig": null,
     	"ixname": "DE-CIX Frankfurt||Main",
     	"options": {},
     	"peeringgroups": {},
     	"peers_configured": {
     		"DE-CIX Frankfurt||Main": {
     			"196922": [
     				{
     					"active": true,
     					"asn": "196922",
     					"group": "",
     					"group6": "",
     					"groupenabled": true,
     					"group6_enabled": true,
     					"infoprefixes4": 0,
     					"infoprefixes6": 0,
     					"ipv4addr": "",
     					"ipv6addr": "",
     					"ipv4enabled": true,
     					"ipv6enabled": true,
     					"irrasset": "",
     					"isrs": false,
     					"isrsper": false,
     					"localpreference": 0,
     					"prefixfilter": false
     				}
     			]
     		}
     	},
     	"peersready": [
     		{
     			"active": true,
     			"asn": "196922",
     			"group": "",
     			"group6": "",
     			"groupenabled": false,
     			"group6_enabled": false,
     			"infoprefixes4": 64,
     			"infoprefixes6": 10,
     			"ipv4addr": "80.81.194.25",
     			"ipv6addr": "2001:7f8::3:13a:0:1",
     			"ipv4enabled": true,
     			"ipv6enabled": true,
     			"irrasset": "AS-HOFMEIR",
     			"isrs": false,
     			"isrsper": false,
     			"localpreference": 0,
     			"prefixfilter": false
     		}
     	],
     	"routeserverready": null
     }
    
### Default output and templates 
By default IXgen will output on the standard output channel. The output can be also redirected to a file with the  _-output_ parameter. Be aware, that the output is always sorted by peers ASN.
 
### Default syntax and more info for Brocade Netiron 
The default output syntax is  Brocade Netiron command line syntax, because this is my home box :D. If you are on one of the  Netiron platforms (MLX,CER,MLXE), you can also use my tool _brocadecli_ ( https://github.com/ipcjk/brocadecli) to automatically upload the configuration into your router, such as with an extra cronjob.
  
###  REST-APIs
Newer routers like the Brocade SLX or JunOS 16.X support a configuration with REST and I will support it as soon as I get my hands on.
  
## Incorporate your own templates 
To use your own router, you need to create or use one of the provided templates in the _templates_-folder and set the _-style_ parameter to your flavor, e.g. _-juniper/set_ for Junos set exchange format.
   
Special output like Juniper JSON is integrated in code.  
  
### templates for router snippets 

The templates directory is very easy structured and has a separate layer for vendor and  devices:
 
 - native
   - json
   - json_pretty
 - brocade 
   - netiron
 - juniper
   - set
   - json (fixed in code, no template)
 - cisco
   - (currently almost an one-to-one-copy from the Brocade-template)
   
   The last layer always has a _router.tt_-, an optional _header.tt_-  and _footer.tt_-file. 
   
   _router.tt_ is the main file,
   that is supplied by the _peergenerator_-module. The _header.tt_-file is a custom file, that will be injected before the peering-code, the _footer.tt_ file after. If you need to initialize peering groups or set any other important value, then _header.tt_ is the right place to be.
     

### Exported structures to template engine 

Exported to the template is the type "IX", that is a struct of the member variables:
- AdditionalConfig (array of strings)
- IxName (name of the IX)
- Options (hash map of IX-options from the ini-file)
- PeeringGroups (used peering groups for that IX, generated from the INI)
- PeersINI (Peers as read from the INI-file [dont use this!])
- PeersReady (Peers that have processed and are ready for the templating)
- RouteServerReady (Routeserver records that have been processed and are ready)

## INI-Configuration 
 
### exchange configuration parameters 

When adding an exchange, there are several options and parameters you can add each on a separate line in the 
_[options]_-subsection. Please avoid special characters or whitespaces/tabs inside strings. 

#### ipv4 
 - routeserver_group=$rs_group (group used for peering with $rs_group )
 - peer_group=$peer_group (group used for peering with neighbors for the _[peers]_-list)
 - routeserver_prefixes=$number ($number is used to overwrite the maximum prefix limit from peeringdb)

#### ipv6  
 - routeserver_group6=$rs_group6 (group used for ipv6-peering with $rs_group6 )
 - peer_group6=$peer_group6 (group used for ipv6-peering with neighbors for the _[peers]_-list)
 - routeserver_prefixes6=$number ($number is used to overwrite the maximum prefix limit from peeringdb)

#### iv6 | ipv4 
 - routeserver=(0=disable, 1=auto-detect and configure neighbor statements for route-servers)
 - rsn_asn=$rsn_asn (explicit set the as number value of the expected remote routeserver , this can protect you from rogue route-servers type from peeringdb )

#### wildcard 
 - wildcard= (0=disable [default], 1=enable, 2=enableAll)
 
 Setting wildcard to enable, will configure all possible neighbors of the exchange, that have an open peering policy. 
 
 Setting wildcard to enableAll, will configure all neighbors from the exchange, not respecting the peering policy. This is good for configuration  testing, benchmarking, history ...! Be sure to set the _-myasn_ parameter on start, so that neighbor statements for your own network will be omitted. 
   
#### additional configuration 
It is possible to add custom lines, that are not interpreted by adding the subsection _[additionalConfig]_. The generator will print out this lines before the peering configuration. You can use this code to generate peer group configuration or anything else that you want to add before the single peer configurations. It is also possible to add things into _header.tt_, see templating above.  
   
### peer configuration parameters 
   
   When adding a peer ASN to a _[peers]_-section, there are several options and parameters you can add-on the same line. All options or parameters are delimited by whitespaces or tabs. Future reader will be improved. 
   
    - ipv4=0 (0 = disable neighbor commands with ipv4 addresses, 1 = enable [default])
    - ipv6=0 (0 = disable neighbor commands with ipv6 addresses, 1 = enable [default])
    - active=0 (0 = ignore the ASN for configuration)
    - group4=0 (0 = dont generate peer-group-statement inside IPv4 template, 1=create [default] )
    - group6=0 (0 = dont generate peer-group-statement inside IPv6 template, 1=create [default] )
    - peer_group=$string (if group4 is enabled, don't take the peer-group-name from the exchange options, instead take $string)
    - peer_group6=$string (if group6 is enabled, don't take the peer-group-name from the exchange options, instead take $string)
    - infoprefixes4 = number (number of prefixes for ipv4, only usage is to overwrite the limit from peeringdb, because sometimes the values from peering are not reflecting current values)
    - infoprefixes6 = number (number of prefixes for ipv6 , only usage is to overwrite the limit from peeringdb, because sometimes the values from peering are not reflecting current values)
    - irrasset = (overwrite AS-Macro to use for prefix-filter builder/bgpq3)
    - prefix_filter=(1=build prefix filter, 0=generate prefix-statement with prefix_list from prefix_list or prefix_list6 if enabled )
    - prefixfilter_aggregate=(1=bgpq3-parameter -A: aggregate as much as possible prefix filter)
    - prefix_list=$name (listname for 1) generate or for 2) include statement (external generated)
    - prefix_list6=$name (listname for statement with prefixname)
    - ipv4_addr=$addr (not implemented yet: only generate peering configuration for the specified neighbor address => fixed peering)
    - ipv6_addr=$addr (not implemented yet: only generate peering configuration for the specified neighbor address => fixed peering)
    - local_pref (not implemented yet: generate route-map setting local-preference values)
    - import (not implemented yet: juniper import-policy for single peer)
    - export (not implemented yet: juniper export-policy for single peer)

#### Overview of the command line options ####
    ixgen:
    -api string
    	use a differnt server as sources instead local/api-service. (default "https://www.peeringdb.com/api")
    -buildcache
    	download json files for caching from peeringdb
    -cacheDir string
      	cache directory for json files from peeringdb (default "./cache")
    -config string
    	Path to peering configuration ini-file (default "./configuration/peering.ini")
    -exchange string
    	only generate configuration for given exchange, default: print all
    -exit
    	exit after generator run,  before printing result (performance run)
    -listenapi string
    	listenAddr for local api service (default "localhost:0")
    -myasn int
      	exclude your own asn from the generator
    -noapiservice
    	do NOT create a local thread for the http api server that uses the json file as sources instead peeringdb.com/api-service.
    -output string
    	if set, will output the configuration to a file, else STDOUT
    -prefixfactor float
    	factor for maximum-prefix numbers (default 1.2)
    -style string
    	Style for routing-config by template, e.g. brocade, juniper, cisco. Also possible: native/json or native/json_pretty for outputting the inside structures (default "brocade/netiron")
    -templates string
    	directory for templates (default "./templates")
    -version
    	prints version and exit
   exit status 2



## Apiserver ##

Ixgen has a second standalone executable, called apiserver. Apiserver can run as a daemon or background thread and serve a few peeringdb-like requests, that are mandatory for using ixgen client from the command line. Also Apiserver is capable of generating your routing configurations if you can post the INI-file in text or JSON-format into the http request. That makes it easy to generate the configuration on the router itself (e.g. Brocade SLX with Ubuntu KVM-management). 
      
### Start an apiserver thread 
     ixapiserver -listenAPI localhost:8563
     
Apiserver is now listening on a localhost socket and port 8563. Apiserver runs a fraction of the original peeringdb-api, so in a sample query we can ask for the DE-CIX Frankfurt exchange:

     curl http://localhost:8563/api/ix?name=DE-CIX%20Frankfurt
     {"data":[{"id":31,"city":"Frankfurt","country":"DE","created":"2010-07-29T00:00:00Z",
     "fac_set":null,"ixlan_set":null,"media":"Ethernet","name":"DE-CIX Frankfurt",
     "name_long":"Deutscher Commercial Internet Exchange","notes":"","org":{"id":0,
     "address1":"","address2":"","city":"","country":"","created":"0001-01-01T00:00:00Z",
     "name":"","notes":"","state":"","status":"","updated":"0001-01-01T00:00:00Z",
     "website":"","zipcode":"0"},"org_id":1187,"policy_email":"sales@de-cix.net",
     "policy_phone":"+49 69 1730 902 12","proto_ipv6":true,"proto_multicast":true,
     "proto_unicast":true,"region_continent":"Europe","status":"ok","tech_email":
     "support@de-cix.net","tech_phone":"+49 69 1730 902 11","updated":"2016-10-12T12:18:05Z",
     "url_stats":"https://www.de-cix.net/locations/germany/frankfurt/statistics","website":
     "https://www.de-cix.net/locations/germany/frankfurt"}],"meta":{}}
 
 
### POST-API of ixgen 
 Things become more interesting, when you want to generate configurations over the network. In the next example we use an INI-style configuration file from the local host and post it to the apiserver. 
  
  Contents of peering.ini:
      
      [DE-CIX Frankfurt||Main]
      [options]
      [peers]
      714 ipv6=0
    
  Lets post it to apiserver and request a Brocade SLX-configuration:
  
    $ curl -X POST --data-binary @peering.ini http://localhost:8563/ixgen/brocade/slx
    router bgp
    neighbor 80.81.193.223 remote-as 714
    address-family ipv4 unicast
    neighbor 80.81.193.223 maximum-prefix 10000
    no neighbor 80.81.193.223 shutdown
    exit-address-family
    ....
    
  Also you can post the peering.ini in JSON. The JSON-format is called internal _native/json_ and can be generated by setting the ixgen _-style_-parameter to _native/json_ or _native/json_pretty_. Be sure to set the HTTP content-type to  _"application/json"_.
   
     $ cat peering.json
     [{"additionalconfig":null,"ixname":"DE-CIX Frankfurt||Main","options":{"DE-CIX Frankfurt||Main":{"wildcard":"0"}},
     "peeringgroups":{},"peers_configured":{"DE-CIX Frankfurt||Main":{"714":[{"active":true,"asn":"714","group":"",
     "group6":"","groupenabled":true,"group6_enabled":true,"infoprefixes4":0,"infoprefixes6":0,"ipv4addr":"",
     "ipv6addr":"","ipv4enabled":true,"ipv6enabled":false,"irrasset":"","isrs":false,"isrsper":false,"localpreference":0,
     "prefixfilter":false}]}}}]
     
    $ curl -X POST --data-binary @peering.json http://localhost:8563/ixgen/brocade/netiron -H "Content-type: application/json"
    router bgp
    neighbor 80.81.193.223 remote-as 714
    address-family ipv4 unicast
    neighbor 80.81.193.223 maximum-prefix 10000
    no neighbor 80.81.193.223 shutdown
    exit-address-family
    ....
    
### Generate router configuration 
   
   The URI for posting configurations is formed out of:
   "/ixgen/vendor/style/myasn"
   
   _Vendor_ and _style_ has the same meaning like the _-style_-parameter on the ixgen command line. The argument _myasn_ is an optional argument to omit generating a configuration for that as-number.
   
   As an example, if you want to print out a Juniper set-configuration and want to omit 196922, you would call the the apiserver with:
     
       $ curl -X POST --data-binary @peering.json http://localhost:8563/ixgen/juniper/set/196922 
   
    
    
### Apiserver command line help 
   
     -cacheDir string
       	cache directory for json files from peeringdb (default "./cache")
     -listenAPI string
       	listenAddr for the api service (default "localhost:8443")
     -templates string
       	directory for templates (default "./templates")


 
 