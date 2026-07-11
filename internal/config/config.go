package config

import (
	"encoding/json"
	"maps"
	"os"
	"os/user"
	"strings"
)

const (
	ThemeDark = "dark"
	ThemeLight = "light"
	ThemeTransparent = "transparent"
)

const (
	AuthMethodIdentityFile = "identity_file"
	AuthMethodPassword     = "password"
)

type Application struct {
	Theme string `json:"theme"`
}

type Defaults struct {
	Username     string `json:"username"`
	Port         string `json:"port"`
	AuthMethod   string `json:"auth_method"`
	IdentityFile string `json:"identity_file"`
}

type Server struct {
	Address      string `json:"address"`
	Username     string `json:"username,omitempty"`
	Port         string `json:"port,omitempty"`
	AuthMethod   string `json:"auth_method,omitempty"`
	IdentityFile string `json:"identity_file,omitempty"`
}

type Config struct {
	Application Application       `json:"application"`
	Defaults    Defaults          `json:"defaults"`
	Servers     map[string]Server `json:"servers"`
}

func (c *Config) SaveApplication(app Application, filePath string) error {
	c.Application = app

	return c.Write(filePath)
}

func (c *Config) Get(name string) Server {
	server := c.Servers[name]
	c.defaults(&server, applyDefaults)

	return server
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

	type ConfigMigration struct {
		Config
		Hosts map[string]Server `json:"hosts"`
	}

	var tmp ConfigMigration
	if err := json.Unmarshal(configFile, &tmp); err != nil {
		return err
	}

	c.Application = tmp.Application
	c.Defaults = tmp.Defaults

	if len(tmp.Hosts) > 0 {
		c.Servers = make(map[string]Server)
		maps.Copy(c.Servers, tmp.Hosts)
	} else {
		c.Servers = tmp.Servers
		if c.Servers == nil {
			c.Servers = make(map[string]Server)
		}
	}

	// Fill defaults
	applyDefaults(&c.Application.Theme, "dark")

	usr, _ := user.Current()
	applyDefaults(&c.Defaults.Username, usr.Username)
	applyDefaults(&c.Defaults.Port, "22")
	applyDefaults(&c.Defaults.AuthMethod, "identity_file")
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
	f(&s.AuthMethod, c.Defaults.AuthMethod)
	f(&s.IdentityFile, c.Defaults.IdentityFile)
}

func applyDefaults(val *string, def string) {
	if strings.TrimSpace(*val) == "" {
		*val = def
	}
}

func stripDefaults(val *string, def string) {
	if strings.TrimSpace(*val) == strings.TrimSpace(def) {
		*val = ""
	}
}
