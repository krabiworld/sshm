package forms

import (
	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/config"
)

const (
	ServerName         = "name"
	ServerAddress      = "address"
	ServerUsername     = "username"
	ServerPort         = "port"
	ServerIdentityFile = "identity_file"
	ServerPassword     = "password"
)

func NewServer(cfg config.Config, currentName string) *huh.Form {
	var (
		name     = currentName
		address  string
		username string
		port     string
		identity string
		password string
	)

	currentServer := cfg.GetOriginal(currentName)
	if currentServer != (config.Server{}) {
		address = currentServer.Address
		username = currentServer.Username
		port = currentServer.Port
		identity = currentServer.IdentityFile
		if currentServer.HasPassword {
			password = "********"
		}
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key(ServerName).
				Title("Name").
				Value(&name).
				Inline(true).
				Validate(validateIsNotEmpty("Name")),
			huh.NewInput().
				Key(ServerAddress).
				Title("Address").
				Value(&address).
				Inline(true).
				Validate(validateIsNotEmpty("Address")),
			huh.NewInput().
				Key(ServerUsername).
				Title("Username").
				Placeholder(cfg.Defaults.Username).
				Value(&username).
				Inline(true),
			huh.NewInput().
				Key(ServerPort).
				Title("Port").
				Placeholder(cfg.Defaults.Port).
				Value(&port).
				Inline(true).
				Validate(validatePort),
			huh.NewInput().
				Key(ServerIdentityFile).
				Title("Identity file").
				Placeholder(cfg.Defaults.IdentityFile).
				Value(&identity).
				Inline(true),
			huh.NewInput().
				Key(ServerPassword).
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&password).
				Inline(true),
			huh.NewConfirm().Affirmative("Save").Negative("Discard").Inline(true),
		),
	).WithWidth(80).WithTheme(FormTheme{})
}
