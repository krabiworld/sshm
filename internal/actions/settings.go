package actions

import (
	"fmt"
	"strings"

	"github.com/krabiworld/sshm/internal/app"
	"github.com/rivo/tview"
)

func Settings(ctx app.Context) {
	var (
		closeAfterConnection = ctx.Config.Settings.CloseAfterConnection
		port                 = ctx.Config.Defaults.Port
		identityFile         = ctx.Config.Defaults.IdentityFile
	)
	form := tview.NewForm()
	form.AddCheckbox("Close sshm after connecting?", closeAfterConnection, func(checked bool) { closeAfterConnection = checked })
	form.AddInputField("Default port", port, 0, nil, func(text string) { port = text })
	form.AddInputField("Default identity file", identityFile, 0, nil, func(text string) { identityFile = text })
	form.AddButton("Save", func() {
		port = strings.TrimSpace(port)
		identityFile = strings.TrimSpace(identityFile)

		if port == "" || identityFile == "" {
			modal := tview.NewModal().
				SetText("Please fill in all fields!").
				AddButtons([]string{"ОК"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					ctx.Pages.RemovePage("add_error_modal")
					ctx.App.SetFocus(form)
				})
			ctx.Pages.AddPage("add_error_modal", modal, true, true)
			return
		}

		ctx.Config.Settings.CloseAfterConnection = closeAfterConnection
		ctx.Config.Defaults.Port = port
		ctx.Config.Defaults.IdentityFile = identityFile
		if err := ctx.Config.Write(ctx.ConfigPath); err != nil {
			ctx.App.Stop()
			fmt.Printf("Error while initializing config: %v\n", err)
			return
		}

		ctx.UpdateTable("")

		ctx.Pages.RemovePage("settings")
		ctx.App.SetFocus(ctx.Table)
	})
	form.SetBorder(true).SetTitle("Settings")

	ctx.Pages.AddPage("settings", tview.NewGrid().
		SetRows(0, 11, 0).
		SetColumns(0, 60, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true), true, true)
	ctx.App.SetFocus(form)
}
