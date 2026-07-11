package actions

import (
	"strings"

	"github.com/krabiworld/sshm/internal/app"
	"github.com/krabiworld/sshm/internal/utils"
	"github.com/rivo/tview"
)

func Settings(ctx app.Context) {
	var (
		theme = ctx.GetApplication().Theme
		port         = ctx.GetDefaults().Port
		authMethod   = ctx.GetDefaults().AuthMethod
		identityFile = ctx.GetDefaults().IdentityFile
	)
	form := tview.NewForm()
	form.AddDropDown("Theme", []string{"Dark", "Light", "Transparent"}, utils.GetIndex(utils.ThemeOrder, theme), func(_ string, optionIndex int) { theme = utils.GetByIndex(utils.ThemeOrder, optionIndex) })
	form.AddInputField("Default port", port, 0, nil, func(text string) { port = text })
	form.AddDropDown("Default auth method", []string{"Identity file", "Password"}, utils.GetIndex(utils.AuthMethodOrder, authMethod), func(_ string, optionIndex int) { authMethod = utils.GetByIndex(utils.AuthMethodOrder, optionIndex) })
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

		app := ctx.GetApplication()
		app.Theme = theme
		utils.CheckError(&ctx, ctx.SaveApplication(app))

		def := ctx.GetDefaults()
		def.Port = port
		def.AuthMethod = authMethod
		def.IdentityFile = identityFile
		utils.CheckError(&ctx, ctx.SaveDefaults(def))

		ctx.Pages.RemovePage("settings")
		ctx.App.SetFocus(ctx.Table)
	})
	form.SetBorder(true).SetTitle("Settings")

	ctx.Pages.AddPage("settings", tview.NewGrid().
		SetRows(0, 13, 0).
		SetColumns(0, 60, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true), true, true)
	ctx.App.SetFocus(form)
}
