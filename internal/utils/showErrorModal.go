package utils

import (
	"github.com/krabiworld/sshm/internal/app"
	"github.com/rivo/tview"
)

func ShowErrorModal(ctx app.Context, text string, returnFocus tview.Primitive) {
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"ОК"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ctx.Pages.RemovePage("error_modal")
			ctx.App.SetFocus(returnFocus)
		})
	ctx.Pages.AddPage("error_modal", modal, true, true)
}
