package storage

import "github.com/krabiworld/sshm/internal/config"

type plaintext struct{
	cfg *config.Config
}

func (c *plaintext) GetPassword(name string) (string, error) {
	return c.cfg.Get(name).Password, nil
}

func (c *plaintext) SetPassword(name, password string) error {
	server := c.cfg.Get(name)
	server.Password = password
	return c.cfg.Set(name, server)
}

func (c *plaintext) DeletePassword(name string) error {
	server := c.cfg.Get(name)
	server.Password = ""
	return c.cfg.Set(name, server)
}
