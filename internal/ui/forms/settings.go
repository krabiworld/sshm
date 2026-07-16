package forms

import (
	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/utils"
)

func NewSettings(cfg *config.Config) *huh.Form {
	var (
		username       = cfg.GetDefaults().Username
		port           = cfg.GetDefaults().Port
		authType       = cfg.GetDefaults().AuthType
		identityFile   = cfg.GetDefaults().IdentityFile
		knownHostsFile = cfg.GetDefaults().KnownHostsFile
	)
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key(ServerUsername).
				Title("Username").
				Placeholder(utils.GetCurrentUsername()).
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
					huh.NewOption("key", config.AuthKey),
					huh.NewOption("password", config.AuthPassword),
					huh.NewOption("agent", config.AuthAgent),
				).
				Value(&authType).
				Inline(true),
			huh.NewInput().
				Key(ServerIdentityFile).
				Title("Identity file").
				Value(&identityFile).
				Inline(true).
				Validate(validateIsNotEmpty("Identity file")),
			huh.NewInput().
				Key(ServerKnownHostsFile).
				Title("Known hosts file").
				Value(&knownHostsFile).
				Inline(true).
				Validate(validateIsNotEmpty("Known hosts file")),
			huh.NewConfirm().Affirmative("Save").Negative("Discard").Inline(true),
		).Title("Settings"),
	).WithWidth(80).WithTheme(FormTheme{})
}
