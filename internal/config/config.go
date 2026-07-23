package config

import (
	"encoding/json"
	"maps"
	"os"

	"github.com/krabiworld/sshm/internal/utils"
)

type (
	AuthType    string
	StorageType string
)

const (
	AuthKey      AuthType = "key"
	AuthPassword AuthType = "password"
	AuthAgent    AuthType = "agent"
)

type Server struct {
	Address      string   `json:"address"`
	Username     string   `json:"username"`
	Port         string   `json:"port"`
	AuthType     AuthType `json:"auth_type"`
	IdentityFile string   `json:"identity_file"`
	Password     string   `json:"password"`
}

type Config struct {
	filePath string
	data     struct {
		Servers map[string]Server `json:"servers"`
	}
}

func New(filePath string) (*Config, error) {
	c := &Config{
		filePath: filePath,
	}
	c.data.Servers = make(map[string]Server)

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
	return c.data.Servers[name]
}

func (c *Config) GetAll() map[string]Server {
	servers := make(map[string]Server, len(c.data.Servers))
	maps.Copy(servers, c.data.Servers)
	return servers
}

func (c *Config) Set(name string, server Server) error {
	c.data.Servers[name] = server

	return c.write()
}

func (c *Config) Delete(name string) error {
	delete(c.data.Servers, name)

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
