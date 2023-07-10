package net

import (
	"log"
	"net"
)

func GetGlobalIpv4Address(ifName string) (string, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range ifaces {
		log.Printf("Found interface %s", i.Name)
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			log.Printf("Found address %s", addr.String())
		}
	}

	return "", nil
}
