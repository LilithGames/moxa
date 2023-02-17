package cluster

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

func DnsResolveRecordA(client *dns.Client, config *dns.ClientConfig, domain string) ([]net.IP, error) {
	m := &dns.Msg{}
	m.SetQuestion(domain, dns.TypeA)
	resp, _, err := client.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port))
	if err != nil {
		return nil, fmt.Errorf("dns.Exchange err: %w", err)
	}
	if resp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("dns.Exchange err code: %d", resp.Rcode)
	}
	ips := make([]net.IP, 0, len(resp.Answer))
	for _, rr := range resp.Answer {
		if record, ok := rr.(*dns.A); ok {
			ips = append(ips, record.A)
		}
	}
	return ips, nil
}
