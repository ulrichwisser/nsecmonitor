package cmd

import (
	mdns "github.com/miekg/dns"
)

var RCODES = []int{
	mdns.RcodeSuccess,
	mdns.RcodeFormatError,
	mdns.RcodeServerFailure,
	mdns.RcodeNameError,
	mdns.RcodeNotImplemented,
	mdns.RcodeRefused,
	mdns.RcodeYXDomain,
	mdns.RcodeYXRrset,
	mdns.RcodeNXRrset,
	mdns.RcodeNotAuth,
	mdns.RcodeNotZone,
	mdns.RcodeBadSig,
	mdns.RcodeBadVers,
	mdns.RcodeBadKey,
	mdns.RcodeBadTime,
	mdns.RcodeBadMode,
	mdns.RcodeBadName,
	mdns.RcodeBadAlg,
	mdns.RcodeBadTrunc,
	mdns.RcodeBadCookie,
}
