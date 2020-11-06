package cmd

var probes = []ProbeSet{
	{
		Requested: 200,
		Type:      "country",
		Value:     "SE",
		Tags: Tags{
			Include: []string{"system-resolves-a-correctly"},
		},
	},
	{
		Requested: 200,
		Type:      "country",
		Value:     "NO",
		Tags: Tags{
			Include: []string{"system-resolves-a-correctly"},
		},
	},
	{
		Requested: 200,
		Type:      "country",
		Value:     "FI",
		Tags: Tags{
			Include: []string{"system-resolves-a-correctly"},
		},
	},
	{
		Requested: 200,
		Type:      "country",
		Value:     "DK",
		Tags: Tags{
			Include: []string{"system-resolves-a-correctly"},
		},
	},
	{
		Requested: 5000,
		Type:      "area",
		Value:     "WW",
		Tags: Tags{
			Include: []string{"system-resolves-a-correctly"},
		},
	},
}

var probesV4 = []ProbeSet{
	{
		Requested: 200,
		Type:      "country",
		Value:     "SE",
		Tags: Tags{
			Include: []string{"system-resolves-a-correctly", "system-ipv4-stable-1d"},
		},
	},
	{
		Requested: 200,
		Type:      "country",
		Value:     "NO",
		Tags: Tags{
			Include: []string{"system-resolves-a-correctly", "system-ipv4-stable-1d"},
		},
	},
	{
		Requested: 200,
		Type:      "country",
		Value:     "FI",
		Tags: Tags{
			Include: []string{"system-resolves-a-correctly", "system-ipv4-stable-1d"},
		},
	},
	{
		Requested: 200,
		Type:      "country",
		Value:     "DK",
		Tags: Tags{
			Include: []string{"system-resolves-a-correctly", "system-ipv4-stable-1d"},
		},
	},
	{
		Requested: 5000,
		Type:      "area",
		Value:     "WW",
		Tags: Tags{
			Include: []string{"system-resolves-a-correctly", "system-ipv4-stable-1d"},
		},
	},
}

var probesV6 = []ProbeSet{
	{
		Requested: 200,
		Type:      "country",
		Value:     "SE",
		Tags: Tags{
			Include: []string{"system-resolves-aaaa-correctly", "system-ipv6-stable-1d"},
		},
	},
	{
		Requested: 200,
		Type:      "country",
		Value:     "NO",
		Tags: Tags{
			Include: []string{"system-resolves-aaaa-correctly", "system-ipv6-stable-1d"},
		},
	},
	{
		Requested: 200,
		Type:      "country",
		Value:     "FI",
		Tags: Tags{
			Include: []string{"system-resolves-aaaa-correctly", "system-ipv6-stable-1d"},
		},
	},
	{
		Requested: 200,
		Type:      "country",
		Value:     "DK",
		Tags: Tags{
			Include: []string{"system-resolves-aaaa-correctly", "system-ipv6-stable-1d"},
		},
	},
	{
		Requested: 5000,
		Type:      "area",
		Value:     "WW",
		Tags: Tags{
			Include: []string{"system-resolves-aaaa-correctly", "system-ipv6-stable-1d"},
		},
	},
}

var definition1 = Definition{
	Type:             "dns",
	AF:               4,
	Description:      "",
	IsPublic:         true,
	QueryClass:       "IN",
	QueryType:        "TXT",
	QueryArgument:    "",
	Retry:            0,
	SetCDBit:         false,
	SetDOBit:         true,
	SetNSIDBit:       true,
	SetRDBit:         true,
	UseProbeResolver: true,
	IncludeAbuf:      true,
	UseMacros:        false,
	TTL:              false,
	Interval:         60,
	Spread:           60,
}

var definition2 = Definition{
	Type:             "dns",
	AF:               4,
	Description:      "Direct to Authoritative",
	IsPublic:         true,
	QueryClass:       "IN",
	QueryType:        "NS",
	QueryArgument:    "",
	Retry:            0,
	SetCDBit:         false,
	SetDOBit:         true,
	SetNSIDBit:       true,
	SetRDBit:         false,
	UseProbeResolver: false,
	ResolveOnProbe:   false,
	IncludeAbuf:      true,
	UseMacros:        false,
	TTL:              false,
	Interval:         60,
	Spread:           30,
	Protocol:         "UDP",
	UDPPayloadSize:   512,
	SkipDNScheck:     false,
	Timeout:          5000,
}

/*
curl --dump-header - -H "Content-Type: application/json" -H "Accept: application/json" -X POST -d '{
 "definitions": [
  {
   "target": "ns.setest.se",
   "af": 6,
   "query_class": "IN",
   "query_type": "NS",
   "query_argument": "nxdomain.setest.se",
   "use_macros": false,
   "description": "TEST v6",
   "interval": 240,
   "use_probe_resolver": false,
   "resolve_on_probe": false,
   "set_nsid_bit": true,
   "protocol": "UDP",
   "udp_payload_size": 512,
   "retry": 0,
   "skip_dns_check": false,
   "include_qbuf": false,
   "include_abuf": true,
   "prepend_probe_id": false,
   "set_rd_bit": false,
   "set_do_bit": true,
   "set_cd_bit": false,
   "timeout": 5000,
   "type": "dns"
  }
 ],
 "probes": [
  {
   "tags": {
    "include": [
     "system-resolves-aaaa-correctly",
     "system-ipv6-stable-1d"
    ],
    "exclude": []
   },
   "type": "country",
   "value": "SE",
   "requested": 208
  }
 ],
 "is_oneoff": false,
 "bill_to": "ulrich@wisser.se",
 "start_time": 1604580920,
 "stop_time": 1604581220
}' https://atlas.ripe.net/api/v2/measurements//?key=YOUR_KEY_HERE

{
  "definitions": [
    {
      "description": "Direct to Authoritative",
      "type": "dns",
      "af": 6,
      "target": "ns.setest.se",
      "spread": 30,
      "is_public": true,
      "interval": 60,
      "protocol": "UDP",
      "query_class": "IN",
      "query_type": "NS",
      "query_argument": "nxdomain.setest.se",
      "retry": 0,
      "set_cd_bit": false,
      "set_do_bit": true,
      "set_nsid_bit": true,
      "set_rd_bit": false,
      "udp_payload_size": 512,
      "use_probe_resolver": false,
      "prepend_probe_id": false,
      "include_qbuf": false,
      "include_abuf": true,
      "use_macros": false,
      "timeout": 5000,
      "ttl": true,
      "skip_dns_check": false
    }
  ],
  "probes": [
    {
      "requested": 200,
      "type": "country",
      "value": "SE",
      "tags": {
        "include": [
          "system-resolves-aaaa-correctly",
          "system-ipv6-stable-1d"
        ]
      }
    }
  ],
  "bill_to": "ulrich@wisser.se",
  "is_oneoff": false,
  "start_time": 1604589315,
  "stop_time": 1604589615
}





*/
/*
curl --dump-header - -H "Content-Type: application/json" -H "Accept: application/json" -X POST -d '{
 "definitions": [
  {
   "af": 4,
   "query_class": "IN",
   "query_type": "TXT",
   "query_argument": "invalid.nu",
   "use_macros": false,
   "description": "DNS measurement",
   "interval": 240,
   "use_probe_resolver": true,
   "resolve_on_probe": false,
   "set_nsid_bit": true,
   "protocol": "UDP",
   "udp_payload_size": 512,
   "retry": 0,
   "skip_dns_check": false,
   "include_qbuf": false,
   "include_abuf": true,
   "prepend_probe_id": false,
   "set_rd_bit": true,
   "set_do_bit": true,
   "set_cd_bit": false,
   "timeout": 5000,
   "type": "dns"
  },
  {
   "af": 4,
   "query_class": "IN",
   "query_type": "TXT",
   "query_argument": "transition2nsec.nu",
   "use_macros": false,
   "description": "DNS measurement",
   "interval": 240,
   "use_probe_resolver": true,
   "resolve_on_probe": false,
   "set_nsid_bit": true,
   "protocol": "UDP",
   "udp_payload_size": 512,
   "retry": 0,
   "skip_dns_check": false,
   "include_qbuf": false,
   "include_abuf": true,
   "prepend_probe_id": false,
   "set_rd_bit": true,
   "set_do_bit": true,
   "set_cd_bit": false,
   "timeout": 5000,
   "type": "dns"
  },
  {
   "af": 4,
   "query_class": "IN",
   "query_type": "TXT",
   "query_argument": "$r-$p-$t-invalid.nu",
   "use_macros": true,
   "description": "DNS measurement",
   "interval": 240,
   "use_probe_resolver": true,
   "resolve_on_probe": false,
   "set_nsid_bit": true,
   "protocol": "UDP",
   "udp_payload_size": 512,
   "retry": 0,
   "skip_dns_check": false,
   "include_qbuf": false,
   "include_abuf": true,
   "prepend_probe_id": false,
   "set_rd_bit": true,
   "set_do_bit": true,
   "set_cd_bit": false,
   "timeout": 5000,
   "type": "dns"
  }
 ],
 "probes": [
  {
   "tags": {
    "include": [
     "system-ipv4-works",
     "system-resolves-a-correctly"
    ],
    "exclude": []
   },
   "type": "country",
   "value": "SE",
   "requested": 500
  }
 ],
 "is_oneoff": false,
 "bill_to": "ulrich@wisser.se",
 "start_time": 1603464624,
 "stop_time": 1603493424
}' https://atlas.ripe.net/api/v2/measurements//?key=YOUR_KEY_HERE
*/
