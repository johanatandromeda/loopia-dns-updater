package dns

import (
	"andromeda.nu/loopia-ipv6-updater/pkg/config"
	"andromeda.nu/loopia-ipv6-updater/pkg/net"
	"fmt"
	"github.com/jonlil/loopia-go"
	"golang.org/x/exp/slog"
	"log"
	gonet "net"
)

func UpdateRecords(conf config.Config, addresses map[string]net.Address, bry bool) {
	client, err := loopia.New(conf.Loopia.Username, conf.Loopia.Password)
	if err != nil {
		log.Fatal(err)
	}
	for _, domain := range conf.Domains {
		aByName := make(map[string]loopia.Record)
		aaaaByName := make(map[string]loopia.Record)
		ifByFqdn4 := make(map[string]string)
		ifByFqdn6 := make(map[string]string)
		var matchUnknown4 gonet.Addr
		var matchUnknown6 gonet.Addr
		for _, iface := range domain.Interfaces {
			if iface.MatchUnknown4 {
				if addr, ok := addresses[iface.IfName]; ok && addr.Ipv4 != nil {
					matchUnknown4 = addr.Ipv4
				} else {
					log.Fatalf("Interface %s is set to match unknown IPv4 but has no public address", iface.IfName)
				}
			}
			if iface.MatchUnknown6 {
				if addr, ok := addresses[iface.IfName]; ok && addr.Ipv6 != nil {
					matchUnknown6 = addr.Ipv6
				} else {
					log.Fatalf("Interface %s is set to match unknown IPv4 but has no public address", iface.IfName)
				}
			}
			for _, fqdn := range iface.Match4 {
				ifByFqdn4[fqdn] = iface.IfName
			}
			for _, fqdn := range iface.Match6 {
				ifByFqdn6[fqdn] = iface.IfName
			}
		}
		slog.Info(fmt.Sprintf("Processing domain %s", domain.Name))
		subdomains, err := client.GetSubdomains(domain.Name)
		if err != nil {
			log.Fatal(err)
		}
		for _, subdomain := range subdomains {
			slog.Debug(fmt.Sprintf("Processing subdomain %s", subdomain.Name))
			records, err := client.GetZoneRecords(domain.Name, subdomain.Name)
			if err != nil {
				log.Fatal(err)
			}
			fqdn := subdomain.Name + "." + domain.Name
			for _, record := range records {
				if record.Type == "A" {
					slog.Debug(fmt.Sprintf("Found A record %s = %s", fqdn, record.Value))
					aByName[fqdn] = record
				}
				if record.Type == "AAAA" {
					slog.Debug(fmt.Sprintf("Found AAAA record %s = %s", fqdn, record.Value))
					aaaaByName[fqdn] = record
				}
			}
		}

		// Update the records
		for fqdn, record := range aByName {
			if ifName, ok := ifByFqdn4[fqdn]; ok {
				addr, ok := addresses[ifName]
				if !ok || addr.Ipv4 == nil {
					log.Fatalf("No IPv4 address found for %s", ifName)
				}
				ifIp, _, err := gonet.ParseCIDR(addr.Ipv4.String())
				if err != nil {
					log.Fatal(err)
				}
				slog.Debug(fmt.Sprintf("Checking whether match known IPv4 %s is set for %s (current value %s)", ifIp, fqdn, record.Value))
			} else if matchUnknown4 != nil {
				ifIp, _, err := gonet.ParseCIDR(matchUnknown4.String())
				if err != nil {
					log.Fatal(err)
				}
				slog.Debug(fmt.Sprintf("Checking whether match unknown IPv4 %s is set for %s (current value %s)", ifIp, fqdn, record.Value))
			}
		}
		for fqdn, record := range aaaaByName {
			if ifName, ok := ifByFqdn6[fqdn]; ok {
				addr, ok := addresses[ifName]
				if !ok || addr.Ipv6 == nil {
					log.Fatalf("No IPv4 address found for %s", ifName)
				}
				newIp, err := applyNet(record.Value, addr.Ipv6)
				if err != nil {
					log.Fatal(err)
				}
				slog.Debug(fmt.Sprintf("Checking whether match known IPv6 %s is set for %s (current value %s)", newIp, fqdn, record.Value))
			} else if matchUnknown6 != nil {
				newIp, err := applyNet(record.Value, matchUnknown6)
				if err != nil {
					log.Fatal(err)
				}
				slog.Debug(fmt.Sprintf("Checking whether match unknown IPv6 %s is set for %s (current value %s)", newIp, fqdn, record.Value))
			}
		}
	}
}

func applyNet(addr string, ifAddr gonet.Addr) (string, error) {
	ifIp, ifNet, err := gonet.ParseCIDR(ifAddr.String())
	if err != nil {
		return "", err
	}
	ifIpBytes := ifIp.To16()
	ifNetBytes := ifNet.Mask
	ip := gonet.ParseIP(addr)
	if ip == nil {
		return "", fmt.Errorf("Invalid IP in DNS: %s", addr)
	}
	ipBytes := ip.To16()
	ipWithNewNet := make([]byte, 16)
	for i := 0; i < 16; i++ {
		ipWithNewNet[i] = ifNetBytes[i]&ifIpBytes[i] + ^ifNetBytes[i]&ipBytes[i]
	}
	newIp := gonet.IP(ipWithNewNet)
	return newIp.String(), nil
}
