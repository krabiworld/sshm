package forms

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func NewMasterPassword() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key(MasterPassword).
				Title("Master Password").
				EchoMode(huh.EchoModePassword).
				Inline(true),
			huh.NewConfirm().
				Key(Confirmed).
				Title("Save to keychain?").
				WithButtonAlignment(lipgloss.Left),
		),
	).WithWidth(40).WithTheme(FormTheme{})
}
