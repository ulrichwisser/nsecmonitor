// types.go

// This file contains the various types used by the API

package cmd

// APIError is for errors returned by the RIPE API.
type APIError struct {
	Err struct {
		Status int    `json:"status"`
		Code   int    `json:"code"`
		Detail string `json:"detail"`
		Title  string `json:"title"`
		Errors []struct {
			Source struct {
				Pointer string
			} `json:"source"`
			Detail string
		} `json:"errors"`
	} `json:"error"`
}

// MeasurementRequest contains the different measurement to create/view
type MeasurementRequest struct {
	// see below for definition
	Definitions []Definition `json:"definitions"`

	// requested set of probes
	Probes []ProbeSet `json:"probes"`
	//
	BillTo       string `json:"bill_to,omitempty"`
	IsOneoff     bool   `json:"is_oneoff"`
	SkipDNSCheck bool   `json:"skip_dns_check,omitempty"`
	Times        int    `json:"times,omitempty"`
	StartTime    int    `json:"start_time,omitempty"`
	StopTime     int    `json:"stop_time,omitempty"`
}

// Tags is system and user tags
type Tags struct {
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
}

// ProbeSet is a set of probes obviously
type ProbeSet struct {
	Requested int    `json:"requested"` // number of probes
	Type      string `json:"type"`      // area, country, prefix, asn, probes, msm
	Value     string `json:"value"`     // can be numeric or string
	Tags      Tags   `json:"tags,omitempty"`
}

// Definition is used to create measurements
type Definition struct {
	// Required fields
	Description string `json:"description"`
	Type        string `json:"type"`
	AF          int    `json:"af"`

	// Required for all but "dns"
	Target string `json:"target,omitempty"`

	GroupID        int      `json:"group_id,omitempty"`
	Group          string   `json:"group,omitempty"`
	InWifiGroup    bool     `json:"in_wifi_group,omitempty"`
	Spread         int      `json:"spread,omitempty"`
	Packets        int      `json:"packets,omitempty"`
	PacketInterval int      `json:"packet_interval,omitempty"`
	Tags           []string `json:"tags,omitempty"`

	// Common parameters
	ExtraWait      int  `json:"extra_wait,omitempty"`
	IsOneoff       bool `json:"is_oneoff,omitempty"`
	IsPublic       bool `json:"is_public,omitempty"`
	ResolveOnProbe bool `json:"resolve_on_probe,omitempty"`

	// Default depends on type
	Interval int `json:"interval,omitempty"`

	// dns & traceroute parameters
	Protocol string `json:"protocol,omitempty"`

	// dns parameters
	QueryClass          string `json:"query_class"`
	QueryType           string `json:"query_type"`
	QueryArgument       string `json:"query_argument"`
	Retry               int    `json:"retry"`
	SetCDBit            bool   `json:"set_cd_bit"`
	SetDOBit            bool   `json:"set_do_bit"`
	SetNSIDBit          bool   `json:"set_nsid_bit"`
	SetRDBit            bool   `json:"set_rd_bit"`
	UDPPayloadSize      int    `json:"udp_payload_size,omitempty"`
	UseProbeResolver    bool   `json:"use_probe_resolver"`
	PrependProbeID      bool   `json:"prepend_probe_id"`
	IncludeQbuf         bool   `json:"include_qbuf"`
	IncludeAbuf         bool   `json:"include_abuf"`
	UseMacros           bool   `json:"use_macros"`
	Timeout             int    `json:"timeout,omitempty"`
	TLS                 bool   `json:"tls,omitempty"`
	DefaultClientSubnet bool   `json:"default_client_subnet,omitempty"`
	Cookies             bool   `json:"cookies,omitempty"`
	TTL                 bool   `json:"ttl,omitempty"`
	SkipDNScheck        bool   `json:"skip_dns_check"`

	// ping parameters
	//   none (see target)

	// traceroute parameters
	DestinationOptionSize int  `json:"destination_option_size,omitempty"`
	DontFragment          bool `json:"dont_fragment,omitempty"`
	DuplicateTimeout      int  `json:"duplicate_timeout,omitempty"`
	FirstHop              int  `json:"first_hop,omitempty"`
	HopByHopOptionSize    int  `json:"hop_by_hop_option_size,omitempty"`
	MaxHops               int  `json:"max_hops,omitempty"`
	Paris                 int  `json:"paris,omitempty"`

	// ntp parameters
	//   none (see target)

	// http parameters
	ExtendedTiming     bool   `json:"extended_timing,omitempty"`
	HeaderBytes        int    `json:"header_bytes,omitempty"`
	Method             string `json:"method,omitempty"`
	MoreExtendedTiming bool   `json:"more_extended_timing,omitempty"`
	Path               string `json:"path,omitempty"`
	QueryOptions       string `json:"query_options,omitempty"`
	UserAgent          string `json:"user_agent,omitempty"`
	Version            string `json:"version,omitempty"`

	// sslcert parameters
	//   none (see target)

	// sslcert & traceroute & http & dns parameters
	Port int `json:"port,omitempty"`

	// ping & traceroute parameters
	Size int `json:"size,omitempty"`

	// wifi parameters
	AnonymousIdentity string `json:"anonymous_identity,omitempty"`
	Cert              string `json:"cert,omitempty"`
	EAP               string `json:"eap,omitempty"`
}

// Response to Measurement creation
type MeasurementResponse struct {
	Measurements []int
}
