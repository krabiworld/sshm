package forms

import (
	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/config"
)

const (
	SettingsUsername     = "username"
	SettingsPort         = "port"
	SettingsIdentityFile = "identity_file"
)

func NewSettings(cfg config.Config) *huh.Form {
	var (
		username     = cfg.Defaults.Username
		port         = cfg.Defaults.Port
		authType     = cfg.Defaults.AuthType
		identityFile = cfg.Defaults.IdentityFile
	)
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key(ServerUsername).
				Title("Username").
				Value(&username).
				Inline(true).
				Validate(validateIsNotEmpty("Username")),
			huh.NewInput().
				Key(ServerPort).
				Title("Port").
				Value(&port).
				Inline(true).
				Validate(validatePort),
			huh.NewSelect[config.AuthType]().
				Key(ServerAuthType).
				Title("Auth type").
				Options(huh.NewOption("key", config.AuthKey), huh.NewOption("password", config.AuthPassword)).
				Value(&authType).
				Inline(true),
			huh.NewInput().
				Key(ServerIdentityFile).
				Title("Identity file").
				Value(&identityFile).
				Inline(true).
				Validate(validateIsNotEmpty("Identity file")),
			huh.NewConfirm().Affirmative("Save").Negative("Discard").Inline(true),
		),
	).WithWidth(80).WithTheme(FormTheme{})
}
