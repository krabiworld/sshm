package config

import (
	"encoding/json"
	"os"
)

type ConfigSettings struct {
	CloseAfterConnection bool `json:"close_after_connection"`
}

type ConfigDefaults struct {
	Port         string `json:"port"`
	IdentityFile string `json:"identity_file"`
}

type ConfigHost struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	ConfigDefaults
}

type Config struct {
	Settings ConfigSettings        `json:"settings"`
	Defaults ConfigDefaults        `json:"defaults"`
	Hosts    map[string]ConfigHost `json:"hosts"`
}

func (c *Config) Get(hostname string) ConfigHost {
	return c.Hosts[hostname]
}

func (c *Config) Delete(hostname string) {
	delete(c.Hosts, hostname)
}

func (c *Config) Read(filePath string) (err error) {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(configFile, c); err != nil {
		return err
	}
	for name, host := range c.Hosts {
		if host.Port == "" {
			host.Port = c.Defaults.Port
		}
		if host.IdentityFile == "" {
			host.IdentityFile = c.Defaults.IdentityFile
		}
		c.Hosts[name] = host
	}
	return
}

func (c *Config) Write(filePath string) error {
	configBytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, configBytes, 0644); err != nil {
		return err
	}
	return nil
}
