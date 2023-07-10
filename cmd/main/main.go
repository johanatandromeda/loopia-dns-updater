package main

import (
	"andromeda.nu/loopia-ipv6-updater/pkg/net"
	"fmt"
	"log"
)

var version = ""

func main() {
	fmt.Printf("Starting loopia-ipv6-updater V %s\n", version)

	ifName := "test"
	ipv4, err := net.GetGlobalIpv4Address(ifName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found global IPv6 %s for interface %s\n", ipv4, ifName)

}
