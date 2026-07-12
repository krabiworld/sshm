package config

import (
	"encoding/json"
	"os"
	"os/user"
)

type Defaults struct {
	Username     string `json:"username"`
	Port         string `json:"port"`
	IdentityFile string `json:"identity_file"`
}

type Server struct {
	Address      string `json:"address"`
	Username     string `json:"username,omitempty"`
	Port         string `json:"port,omitempty"`
	IdentityFile string `json:"identity_file,omitempty"`
	HasPassword  bool   `json:"has_password,omitempty"`
}

type Config struct {
	Defaults Defaults          `json:"defaults"`
	Servers  map[string]Server `json:"servers"`
}

func (c *Config) Get(name string) Server {
	server := c.Servers[name]
	c.defaults(&server, applyDefaults)

	return server
}

func (c *Config) GetAll() map[string]Server {
	servers := make(map[string]Server, len(c.Servers))
	for name, server := range c.Servers {
		c.defaults(&server, applyDefaults)
		servers[name] = server
	}
	return servers
}

func (c *Config) GetOriginal(name string) Server {
	return c.Servers[name]
}

func (c *Config) Save(name string, server Server, filePath string) error {
	c.defaults(&server, stripDefaults)
	c.Servers[name] = server

	return c.Write(filePath)
}

func (c *Config) Delete(name, filePath string) error {
	delete(c.Servers, name)

	return c.Write(filePath)
}

func (c *Config) SaveDefaults(def Defaults, filePath string) error {
	c.Defaults = def

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
	usr, _ := user.Current()
	applyDefaults(&c.Defaults.Username, usr.Username)
	applyDefaults(&c.Defaults.Port, "22")
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

func (c *Config) defaults(s *Server, f func(*string, string)) {
	f(&s.Username, c.Defaults.Username)
	f(&s.Port, c.Defaults.Port)
	f(&s.IdentityFile, c.Defaults.IdentityFile)
}

func applyDefaults(val *string, def string) {
	if *val == "" {
		*val = def
	}
}

func stripDefaults(val *string, def string) {
	if *val == def {
		*val = ""
	}
}
