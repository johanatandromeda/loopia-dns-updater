package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	IfName string `yaml:"ifName"`
	Loopia struct {
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
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
	return cfg, nil
}
