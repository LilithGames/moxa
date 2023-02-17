package cluster

import (
	"fmt"
	"testing"
	"time"

	"github.com/miekg/dns"

	assert "github.com/stretchr/testify/require"
)

func Test_Dns(t *testing.T) {
	config, err := GetDefaultDnsClientConfig()
	// config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	assert.Nil(t, err)
	client := &dns.Client{Timeout: time.Second}
	ips, err := DnsResolveRecordA(client, config, "www.baidu.com.")
	assert.Nil(t, err)
	fmt.Printf("%+v\n", ips)
}
