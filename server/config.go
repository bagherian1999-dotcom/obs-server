package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Domain      string `json:"domain"`
	EnableSSL   bool   `json:"enable_ssl"`
	SSLCertPath string `json:"ssl_cert_path"`
	SSLKeyPath  string `json:"ssl_key_path"`
	DataDir     string `json:"data_dir"`
}

var config Config

func loadConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		// Return default config if file doesn't exist
		config = Config{
			Host:      "0.0.0.0",
			Port:      "8080",
			Domain:    "",
			EnableSSL: false,
			DataDir:   "./data",
		}
		return nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&config)
}

func (c *Config) GetAddress() string {
	return c.Host + ":" + c.Port
}

func (c *Config) GetPublicURL() string {
	scheme := "http"
	if c.EnableSSL {
		scheme = "https"
	}
	
	if c.Domain != "" {
		if c.EnableSSL && c.Port == "443" {
			return scheme + "://" + c.Domain
		} else if !c.EnableSSL && c.Port == "80" {
			return scheme + "://" + c.Domain
		}
		return scheme + "://" + c.Domain + ":" + c.Port
	}
	
	return scheme + "://localhost:" + c.Port
}
