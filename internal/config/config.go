package config

import (
	"encoding/json"
	"os"
)

const (
	AuthMethodIdentityFile = "identity_file"
	AuthMethodPassword     = "password"
)

type ConfigDefaults struct {
	Port         string `json:"port"`
	AuthMethod   string `json:"auth_method"`
	IdentityFile string `json:"identity_file"`
}

type ConfigHost struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	ConfigDefaults
}

type Config struct {
	Defaults ConfigDefaults        `json:"defaults"`
	Hosts    map[string]ConfigHost `json:"hosts"`
}

func (c *Config) Read(filePath string) (err error) {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(configFile, c); err != nil {
		return err
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
