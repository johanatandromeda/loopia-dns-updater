package dns

import (
	"andromeda.nu/loopia-ipv6-updater/pkg/config"
	"github.com/jonlil/loopia-go"
	"log"
)

func FindRecords(conf config.Config) {
	client, err := loopia.New(conf.Loopia.Username, conf.Loopia.Password)
	if err != nil {
		log.Fatal(err)
	}
	for _, domain := range conf.Domain {
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
			for _, record := range records {
				log.Printf("Found record %s", record)
			}
		}
	}
}
