package main

import (
	"andromeda.nu/loopia-ipv6-updater/pkg/config"
	"andromeda.nu/loopia-ipv6-updater/pkg/dns"
	"andromeda.nu/loopia-ipv6-updater/pkg/net"
	"flag"
	"fmt"
	"golang.org/x/exp/slog"
	"log"
	"os"
)

var version = ""

func main() {

	var programLevel = new(slog.LevelVar) // Info by default
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

	fmt.Printf("Starting loopia-ipv6-updater V %s\n", version)

	configFile := flag.String("c", "/etc/loopia-ipv6-updater.yaml", "Config file")
	help := flag.Bool("h", false, "Show help")
	debug := flag.Bool("d", false, "Debug")
	quiet := flag.Bool("q", false, "Quiet")
	dry := flag.Bool("n", false, "Dry run")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *debug {
		programLevel.Set(slog.LevelDebug)
	}

	if *quiet {
		programLevel.Set(slog.LevelWarn)
	}

	slog.Info("Using config file " + *configFile)

	config, err := config.ReadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	addresses, err := net.GetGlobalAddresses(config)
	if err != nil {
		log.Fatal(err)
	}

	dns.UpdateRecords(config, addresses, *dry)
}
