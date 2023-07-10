package net

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func GetGlobalIpv4Address(ifName string) (string, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range ifaces {
		log.Printf("Found interface %s", i.Name)
		if i.Name == ifName {
			addrs, err := i.Addrs()
			if err != nil {
				return "", err
			}
			for _, addr := range addrs {
				ip, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					return "", err
				}
				log.Printf("Found address %s", ip)
				if ip.IsGlobalUnicast() && !ip.IsPrivate() && strings.Contains(ip.String(), ":") {
					return ip.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("No global IPv6 address found for interface %s", ifName)
}
