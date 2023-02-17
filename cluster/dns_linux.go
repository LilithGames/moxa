//go:build linux

package cluster

import (
	"github.com/miekg/dns"
)

func GetDefaultDnsClientConfig() (*dns.ClientConfig, error) {
	return dns.ClientConfigFromFile("/etc/resolv.conf")
}
