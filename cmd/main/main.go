package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/johanatandromeda/loopia-dns-updater/pkg/config"
	"github.com/johanatandromeda/loopia-dns-updater/pkg/dns"
	"github.com/johanatandromeda/loopia-dns-updater/pkg/net"
	"log"
	"log/slog"
	"os"
	"path"
	"sort"
)

var version = ""

func main() {

	var programLevel = new(slog.LevelVar) // Info by default
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

	fmt.Printf("Starting loopia-ipv6-updater V %s\n", version)

	configFile := flag.String("c", "/etc/loopia-dns-updater.yaml", "Config file")
	dataDir := flag.String("s", "/var/lib/loopia-dns-updater", "Data directory")
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

	changed, err := checkIfChanged(addresses, *dataDir)
	if err != nil {
		log.Fatal(err)
	}
	if changed {
		dns.UpdateRecords(config, addresses, *dry)
	} else {
		slog.Info("Interfaces have not changed IP")
	}
	if !*dry {
		writeIpState(addresses, *dataDir)
	}
}

func checkIfChanged(addresses map[string]net.Address, dataDir string) (bool, error) {
	if _, err := os.Stat(dataDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dataDir, 0o700)
		if err != nil {
			return false, err
		}
	}
	dataFile := path.Join(dataDir, "ifstate")
	_, exists := os.Stat(dataFile)
	if os.IsNotExist(exists) {
		slog.Debug("No previous run detected")
		return true, nil
	}
	oldIpsBytes, err := os.ReadFile(dataFile)
	if err != nil {
		_ = os.Remove(dataFile)
		slog.Warn(fmt.Sprintf("Old interface IP file %s corrupted. Deleteing it", dataFile))
		return true, nil
	}
	oldIps := string(oldIpsBytes)
	newIps := calculateIpState(addresses)
	return oldIps != newIps, nil
}

func writeIpState(addresses map[string]net.Address, dataDir string) {
	newIps := calculateIpState(addresses)
	dataFile := path.Join(dataDir, "ifstate")
	_ = os.WriteFile(dataFile, []byte(newIps), 0o600)
}

func calculateIpState(addresses map[string]net.Address) string {
	ifLines := make([]string, 0, 10)
	for ifName, addr := range addresses {
		var ipv4 string
		var ipv6 string
		if addr.Ipv4 != nil {
			ipv4 = addr.Ipv4.String()
		}
		if addr.Ipv6 != nil {
			ipv6 = addr.Ipv6.String()
		}
		ifLines = append(ifLines, ifName+"-"+ipv4+"-"+ipv6)
	}
	sort.Strings(ifLines)
	var b bytes.Buffer
	for _, l := range ifLines {
		b.WriteString(l + "\n")
	}
	newIps := b.String()
	return newIps
}
