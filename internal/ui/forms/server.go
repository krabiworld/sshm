package forms

import (
	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/config"
)

func NewServer(cfg *config.Config, currentName string) *huh.Form {
	var (
		title          = "Add new server"
		name           = currentName
		address        string
		username       string
		port           string
		authType       = cfg.GetDefaults().AuthType
		identity       string
		knownHostsFile string
		password       string
	)

	currentServer := cfg.GetRaw(currentName)
	if currentServer != (config.Server{}) {
		title = "Modify server"
		address = currentServer.Address
		username = currentServer.Username
		port = currentServer.Port
		authType = currentServer.AuthType
		if authType == "" {
			authType = cfg.GetDefaults().AuthType
		}
		identity = currentServer.IdentityFile
		knownHostsFile = currentServer.KnownHostsFile
		if currentServer.HasPassphrase || currentServer.AuthType == config.AuthPassword || cfg.GetDefaults().AuthType == config.AuthPassword {
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
				Placeholder(cfg.GetDefaults().Username).
				Value(&username).
				Inline(true),
			huh.NewInput().
				Key(ServerPort).
				Title("Port").
				Placeholder(cfg.GetDefaults().Port).
				Value(&port).
				Inline(true).
				Validate(validatePort),
			huh.NewSelect[config.AuthType]().
				Key(ServerAuthType).
				Title("Auth type").
				Options(
					huh.NewOption("key", config.AuthKey),
					huh.NewOption("password", config.AuthPassword),
					huh.NewOption("agent", config.AuthAgent),
				).
				Value(&authType).
				Inline(true),
			huh.NewInput().
				Key(ServerIdentityFile).
				Title("Identity file").
				Placeholder(cfg.GetDefaults().IdentityFile).
				Value(&identity).
				Inline(true),
			huh.NewInput().
				Key(ServerKnownHostsFile).
				Title("Known hosts file").
				Placeholder(cfg.GetDefaults().KnownHostsFile).
				Value(&knownHostsFile).
				Inline(true),
			huh.NewInput().
				Key(ServerPassword).
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Placeholder(password).
				Inline(true),
			huh.NewConfirm().Affirmative("Save").Negative("Discard").Inline(true),
		).Title(title),
	).WithWidth(80).WithTheme(FormTheme{})
}
