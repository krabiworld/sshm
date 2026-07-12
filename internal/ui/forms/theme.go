package forms

import "charm.land/huh/v2"

type FormTheme struct{}

func (FormTheme) Theme(isDark bool) *huh.Styles {
	return huh.ThemeDracula(isDark)
}
