package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Domain []struct {
		Name       string `yaml:"name"`
		Interfaces []struct {
			IfName string `yaml:"ifName"`
			Prefix int    `yaml:"prefix"`
		} `yaml:"interfaces"`
	} `yaml:"domain"`
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
	return cfg, nil
}
