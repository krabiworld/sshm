package config

import (
	"encoding/json"
	"os"

	"github.com/krabiworld/sshm/internal/utils"
)

type AuthType string

const (
	AuthKey      AuthType = "key"
	AuthPassword AuthType = "password"
)

type Defaults struct {
	Username     string   `json:"username,omitempty"`
	Port         string   `json:"port"`
	AuthType     AuthType `json:"auth_type"`
	IdentityFile string   `json:"identity_file"`
}

type Server struct {
	Address       string   `json:"address"`
	Username      string   `json:"username,omitempty"`
	Port          string   `json:"port,omitempty"`
	AuthType      AuthType `json:"auth_type,omitempty"`
	IdentityFile  string   `json:"identity_file,omitempty"`
	HasPassphrase bool     `json:"has_passphrase,omitempty"`
}

type Config struct {
	Defaults Defaults          `json:"defaults"`
	Servers  map[string]Server `json:"servers"`
}

func (c *Config) Get(name string) Server {
	server := c.Servers[name]
	c.defaults(&server, false)

	return server
}

func (c *Config) GetAll() map[string]Server {
	servers := make(map[string]Server, len(c.Servers))
	for name, server := range c.Servers {
		c.defaults(&server, false)
		servers[name] = server
	}
	return servers
}

func (c *Config) GetOriginal(name string) Server {
	return c.Servers[name]
}

func (c *Config) Save(name string, server Server, filePath string) error {
	c.defaults(&server, true)
	c.Servers[name] = server

	return c.Write(filePath)
}

func (c *Config) Delete(name, filePath string) error {
	delete(c.Servers, name)

	return c.Write(filePath)
}

func (c *Config) SaveDefaults(def Defaults, filePath string) error {
	c.Defaults = def
	for name, server := range c.Servers {
		c.defaults(&server, true)
		c.Servers[name] = server
	}

	return c.Write(filePath)
}

func (c *Config) Read(filePath string) (err error) {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(configFile, &c); err != nil {
		return err
	}

	// Fill defaults
	applyDefaults(&c.Defaults.Port, "22")
	applyDefaults(&c.Defaults.AuthType, AuthKey)
	applyDefaults(&c.Defaults.IdentityFile, "~/.ssh/id_rsa")

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

func (c *Config) defaults(s *Server, strip bool) {
	if !strip && s.Username == "" && c.Defaults.Username == "" {
		s.Username = utils.GetCurrentUsername()
	}

	if strip {
		stripDefaults(&s.Username, c.Defaults.Username)
		stripDefaults(&s.Port, c.Defaults.Port)
		stripDefaults(&s.AuthType, c.Defaults.AuthType)
		stripDefaults(&s.IdentityFile, c.Defaults.IdentityFile)
	} else {
		applyDefaults(&s.Username, c.Defaults.Username)
		applyDefaults(&s.Port, c.Defaults.Port)
		applyDefaults(&s.AuthType, c.Defaults.AuthType)
		applyDefaults(&s.IdentityFile, c.Defaults.IdentityFile)
	}
}

func applyDefaults[T comparable](val *T, def T) {
	if *val == *new(T) {
		*val = def
	}
}

func stripDefaults[T comparable](val *T, def T) {
	if *val == def {
		*val = *new(T)
	}
}
