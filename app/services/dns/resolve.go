package dns

import (
	"strings"

	"github.com/miekg/dns"
)

func (s *Service) resolveARecord(dnsServer string, hostname string) ([]string, error) {
	dnsServer = strings.TrimSpace(dnsServer)
	if !strings.Contains(dnsServer, ":") {
		dnsServer += ":53"
	}

	dClient := dns.Client{}
	dMsg := dns.Msg{}
	dMsg.SetQuestion(dns.Fqdn(hostname), dns.TypeA)

	resp, _, err := dClient.Exchange(&dMsg, dnsServer)
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, ans := range resp.Answer {
		if aRecord, ok := ans.(*dns.A); ok {
			ips = append(ips, aRecord.A.String())
		}
	}

	return ips, nil
}
