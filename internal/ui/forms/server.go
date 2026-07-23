package forms

import (
	"fmt"

	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/utils"
)

func NewServer(cfg *config.Config, currentName string) *huh.Form {
	var (
		title    = "Add new server"
		name     = currentName
		address  string
		username = utils.GetCurrentUsername()
		port     = "22"
		authType = config.AuthKey
		identity = "~/.ssh/id_rsa"
		password string
	)

	currentServer := cfg.Get(currentName)
	if currentServer != (config.Server{}) {
		title = "Modify server"
		address = currentServer.Address
		username = currentServer.Username
		port = currentServer.Port
		authType = currentServer.AuthType
		identity = currentServer.IdentityFile
		if len(currentServer.Password) > 0 {
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
				Validate(func(s string) error {
					if err := validateIsNotEmpty("Name")(s); err != nil {
						return err
					}
					if server := cfg.Get(s); server != (config.Server{}) && currentName != s {
						return fmt.Errorf("A server named %s already exists.", s)
					}
					return nil
				}),
			huh.NewInput().
				Key(ServerAddress).
				Title("Address").
				Value(&address).
				Inline(true).
				Validate(validateIsNotEmpty("Address")),
			huh.NewInput().
				Key(ServerUsername).
				Title("Username").
				Value(&username).
				Inline(true),
			huh.NewInput().
				Key(ServerPort).
				Title("Port").
				Value(&port).
				Inline(true).
				Validate(validatePort),
			huh.NewSelect[config.AuthType]().
				Key(ServerAuthType).
				Title("Auth type").
				Options(
					huh.NewOption("Key", config.AuthKey),
					huh.NewOption("Password", config.AuthPassword),
					huh.NewOption("Agent", config.AuthAgent),
				).
				Value(&authType).
				Inline(true),
			huh.NewInput().
				Key(ServerIdentityFile).
				Title("Identity file").
				Value(&identity).
				Inline(true),
			huh.NewInput().
				Key(ServerPassword).
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Placeholder(password).
				Inline(true),
			huh.NewConfirm().
				Key(Confirmed).
				Affirmative("Save").
				Negative("Discard").
				Inline(true),
		).Title(title),
	).WithWidth(80).WithTheme(FormTheme{})
}
