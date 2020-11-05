package cmd

import (
	mdns "github.com/miekg/dns"
)

func nsec(rrset []mdns.RR) string {
	// send data to queue
	for _, rr := range rrset {
		switch rr.Header().Rrtype {
		case mdns.TypeNSEC:
			return NSEC
		case mdns.TypeNSEC3:
			return NSEC3
		}
	}
	return NONSEC
}
