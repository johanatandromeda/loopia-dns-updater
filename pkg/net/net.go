package net

import (
	"fmt"
	"github.com/johanatandromeda/loopia-dns-updater/pkg/config"
	"log/slog"
	"net"
	"strings"
)

type Address struct {
	Ipv4 net.Addr
	Ipv6 net.Addr
}

func GetGlobalAddresses(config config.Config) (map[string]Address, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ifInConf := map[string]struct{}{}

	for _, domain := range config.Domains {
		for _, iface := range domain.Interfaces {
			ifInConf[iface.IfName] = struct{}{}
		}
	}

	addresses := make(map[string]Address)

	for _, i := range ifaces {
		slog.Debug("Found interface " + i.Name)
		if _, ok := ifInConf[i.Name]; ok {
			addrs, err := i.Addrs()
			if err != nil {
				return nil, err
			}
			for _, addr := range addrs {
				ip, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					return nil, err
				}
				slog.Debug(fmt.Sprintf("Found address %s", ip))
				if ip.IsGlobalUnicast() && !ip.IsPrivate() {
					if strings.Contains(ip.String(), ":") {
						slog.Debug(fmt.Sprintf("Adding global IPv6 address %s for %s", addr, i.Name))
						a := getAddressEntry(i.Name, addresses)
						a.Ipv6 = addr
						addresses[i.Name] = a
					} else if strings.Contains(ip.String(), ".") {
						slog.Debug(fmt.Sprintf("Adding global IPv4 address %s for %s", addr, i.Name))
						a := getAddressEntry(i.Name, addresses)
						a.Ipv4 = addr
						addresses[i.Name] = a
					}
				}
			}
		}
	}

	// Sanity check
	for ifName, addr := range addresses {
		if addr.Ipv6 == nil && addr.Ipv4 == nil {
			return nil, fmt.Errorf("No global IPv4 or IPv6 address found for interface %s", ifName)
		}
	}

	return addresses, nil
}

func getAddressEntry(name string, addresses map[string]Address) Address {
	if addr, ok := addresses[name]; ok {
		return addr
	} else {
		addr = Address{}
		return addr
	}
}
