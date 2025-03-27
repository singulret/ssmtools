package ssmutil

import (
	"log"
	"net"
)

func ResolveHostnameToIP(hostname string) string {
	ips, err := net.LookupHost(hostname)
	if err != nil || len(ips) == 0 {
		log.Printf("Warning: Unable to resolve %s via DNS.", hostname)
		return ""
	}
	return ips[0]
}
