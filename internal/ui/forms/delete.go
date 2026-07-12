package forms

import (
	"fmt"

	"charm.land/huh/v2"
)

const DeleteConfirmed = "delete_confirmed"

func NewDelete(name string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(DeleteConfirmed).
				Title(fmt.Sprintf("Are you sure you want to delete %s?", name)),
		),
	).WithWidth(50).WithTheme(FormTheme{})
}
