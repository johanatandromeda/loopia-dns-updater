package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Domains []struct {
		Name       string `yaml:"name"`
		Interfaces []struct {
			IfName        string   `yaml:"ifName"`
			MatchUnknown4 bool     `yaml:"matchUnknown4"`
			MatchUnknown6 bool     `yaml:"matchUnknown6"`
			Match4        []string `yaml:"match4"`
			Match6        []string `yaml:"match6"`
		} `yaml:"interfaces"`
	} `yaml:"domains"`
	Loopia struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"loopia"`
}

func ReadConfig(fileName string) (Config, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	cfg.SanityCheck()
	return cfg, nil
}

func (c *Config) SanityCheck() {
	for _, domain := range c.Domains {
		foundMatchUnknown4 := false
		foundMatchUnknown6 := false
		for _, iface := range domain.Interfaces {
			if iface.MatchUnknown4 {
				if foundMatchUnknown4 {
					log.Fatal("Duplicate matchUnknown4 in domain %s", domain.Name)
				}
				foundMatchUnknown4 = true
			}
			if iface.MatchUnknown6 {
				if foundMatchUnknown6 {
					log.Fatal("Duplicate matchUnknown6 in domain %s", domain.Name)
				}
				foundMatchUnknown6 = true
			}
		}
	}
}
