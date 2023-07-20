package dns

import (
	"andromeda.nu/loopia-ipv6-updater/pkg/config"
	"andromeda.nu/loopia-ipv6-updater/pkg/net"
	"github.com/jonlil/loopia-go"
	"log"
)

func FindRecords(conf config.Config, addresses map[string]net.Address) {
	client, err := loopia.New(conf.Loopia.Username, conf.Loopia.Password)
	if err != nil {
		log.Fatal(err)
	}
	for _, domain := range conf.Domain {
		aByName := make(map[string]string)
		aaaaByName := make(map[string]string)
		log.Printf("Processing domain %s", domain.Name)
		subdomains, err := client.GetSubdomains(domain.Name)
		if err != nil {
			log.Fatal(err)
		}
		for _, subdomain := range subdomains {
			log.Printf("Processing subdomain %s", subdomain.Name)
			records, err := client.GetZoneRecords(domain.Name, subdomain.Name)
			if err != nil {
				log.Fatal(err)
			}
			fqdn := subdomain.Name + "." + domain.Name
			for _, record := range records {
				if record.Type == "A" {
					log.Printf("Found A record %s = %s", fqdn, record.Value)
					aByName[subdomain.Name] = record.Value
				}
				if record.Type == "AAAA" {
					log.Printf("Found AAAA record %s = %s", fqdn, record.Value)
					aaaaByName[subdomain.Name] = record.Value
				}
			}
		}
	}
}
