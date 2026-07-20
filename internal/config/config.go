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
	AuthAgent    AuthType = "agent"
)

type Defaults struct {
	Username       string   `json:"username,omitempty"`
	Port           string   `json:"port"`
	AuthType       AuthType `json:"auth_type"`
	IdentityFile   string   `json:"identity_file"`
	KnownHostsFile string   `json:"known_hosts_file"`
}

type Server struct {
	Address        string   `json:"address"`
	Username       string   `json:"username,omitempty"`
	Port           string   `json:"port,omitempty"`
	AuthType       AuthType `json:"auth_type,omitempty"`
	IdentityFile   string   `json:"identity_file,omitempty"`
	HasPassphrase  bool     `json:"has_passphrase,omitempty"`
	KnownHostsFile string   `json:"known_hosts_file,omitempty"`
}

type Config struct {
	filePath string
	data     struct {
		Defaults Defaults          `json:"defaults"`
		Servers  map[string]Server `json:"servers"`
	}
}

func New(filePath string) (*Config, error) {
	c := &Config{
		filePath: filePath,
	}
	c.data.Servers = make(map[string]Server)

	c.ensureBaseDefaults()

	configFile, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return c, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(configFile, &c.data); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Get(name string) Server {
	server := c.data.Servers[name]
	c.defaults(&server, false)

	return server
}

func (c *Config) GetAll() map[string]Server {
	servers := make(map[string]Server, len(c.data.Servers))
	for name, server := range c.data.Servers {
		c.defaults(&server, false)
		servers[name] = server
	}
	return servers
}

func (c *Config) GetRaw(name string) Server {
	return c.data.Servers[name]
}

func (c *Config) Set(name string, server Server) error {
	c.defaults(&server, true)
	c.data.Servers[name] = server

	return c.write()
}

func (c *Config) Delete(name string) error {
	delete(c.data.Servers, name)

	return c.write()
}

func (c *Config) GetDefaults() Defaults {
	return c.data.Defaults
}

func (c *Config) SetDefaults(def Defaults) error {
	c.data.Defaults = def
	for name, server := range c.data.Servers {
		c.defaults(&server, true)
		c.data.Servers[name] = server
	}

	return c.write()
}

func (c *Config) write() error {
	bytes, err := json.MarshalIndent(c.data, "", "\t")
	if err != nil {
		return err
	}

	if err := utils.CreateSshDir(); err != nil {
		return err
	}

	return os.WriteFile(c.filePath, bytes, 0600)
}

func (c *Config) defaults(s *Server, strip bool) {
	if !strip && s.Username == "" && c.data.Defaults.Username == "" {
		s.Username = utils.GetCurrentUsername()
	}

	if strip {
		stripDefaults(&s.Username, c.data.Defaults.Username)
		stripDefaults(&s.Port, c.data.Defaults.Port)
		stripDefaults(&s.AuthType, c.data.Defaults.AuthType)
		stripDefaults(&s.IdentityFile, c.data.Defaults.IdentityFile)
		stripDefaults(&s.KnownHostsFile, c.data.Defaults.KnownHostsFile)
	} else {
		applyDefaults(&s.Username, c.data.Defaults.Username)
		applyDefaults(&s.Port, c.data.Defaults.Port)
		applyDefaults(&s.AuthType, c.data.Defaults.AuthType)
		applyDefaults(&s.IdentityFile, c.data.Defaults.IdentityFile)
		applyDefaults(&s.KnownHostsFile, c.data.Defaults.KnownHostsFile)
	}
}

func (c *Config) ensureBaseDefaults() {
	applyDefaults(&c.data.Defaults.Port, "22")
	applyDefaults(&c.data.Defaults.AuthType, AuthKey)
	applyDefaults(&c.data.Defaults.IdentityFile, "~/.ssh/id_rsa")
	applyDefaults(&c.data.Defaults.KnownHostsFile, "~/.ssh/known_hosts")
}
