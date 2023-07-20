package main

import (
	"andromeda.nu/loopia-ipv6-updater/pkg/config"
	"andromeda.nu/loopia-ipv6-updater/pkg/dns"
	"andromeda.nu/loopia-ipv6-updater/pkg/net"
	"flag"
	"fmt"
	"log"
	"os"
)

var version = ""

func main() {
	fmt.Printf("Starting loopia-ipv6-updater V %s\n", version)

	configFile := flag.String("c", "/etc/loopia-ipv6-updater.yaml", "Config file")
	help := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	log.Printf("Using config file %s", *configFile)

	config, err := config.ReadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	addresses, err := net.GetGlobalAddresses(config)
	if err != nil {
		log.Fatal(err)
	}

	dns.FindRecords(config, addresses)
}
