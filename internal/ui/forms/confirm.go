package forms

import (
	"strings"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/mattn/go-runewidth"
)

const FormConfirmed = "form_confirmed"

func NewConfirm(title string) *huh.Form {
	firstLine := title
	if before, _, ok := strings.Cut(title, "\n"); ok {
		firstLine = before
	}

	width := runewidth.StringWidth(firstLine) + 4

	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(FormConfirmed).
				Title(title).
				WithButtonAlignment(lipgloss.Left),
		),
	).WithWidth(width).WithTheme(FormTheme{})
}
