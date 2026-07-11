package utils

import (
	"fmt"
	"strings"

	"github.com/krabiworld/sshm/internal/app"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/rivo/tview"
	"github.com/zalando/go-keyring"
)

func WriteServer(ctx app.Context, oldName string) {
	isModify := oldName != ""
	title := "Create Server"

	var (
		name         = oldName
		address      string
		username     string
		port         string
		authMethod   = ctx.GetDefaults().AuthMethod
		identityFile string
		password     string
	)

	if isModify {
		title = "Modify Server"
		serverCfg := ctx.GetServer(oldName)
		address = serverCfg.Address
		username = serverCfg.Username
		port = serverCfg.Port
		if serverCfg.AuthMethod != "" {
			authMethod = serverCfg.AuthMethod
		}
		identityFile = serverCfg.IdentityFile
		password, _ = keyring.Get("sshm", fmt.Sprintf("%s-%s", serverCfg.Address, serverCfg.Username))
	}

	def := ctx.GetDefaults()

	form := tview.NewForm()
	form.AddInputField("Name *", name, 0, nil, func(text string) { name = text })
	form.AddInputField("Address *", address, 0, nil, func(text string) { address = text })
	form.AddDropDown("Auth method *", []string{"Identity file", "Password"}, GetIndex(AuthMethodOrder, authMethod), func(_ string, optionIndex int) { authMethod = GetByIndex(AuthMethodOrder, optionIndex) })
	form.AddFormItem(tview.NewInputField().
		SetLabel("Username").
		SetText(username).
		SetPlaceholder(def.Username).
		SetChangedFunc(func(text string) { username = text }))
	form.AddFormItem(tview.NewInputField().
		SetLabel("Port").
		SetText(port).
		SetPlaceholder(def.Port).
		SetChangedFunc(func(text string) { port = text }))
	form.AddFormItem(tview.NewInputField().
		SetLabel("Identity file").
		SetText(identityFile).
		SetPlaceholder(def.IdentityFile).
		SetChangedFunc(func(text string) { identityFile = text }))
	form.AddPasswordField("Password", password, 0, '*', func(text string) { password = text })

	form.AddButton("Save", func() {
		name = strings.TrimSpace(name)
		address = strings.TrimSpace(address)
		username = strings.TrimSpace(username)
		port = strings.TrimSpace(port)
		identityFile = strings.TrimSpace(identityFile)
		password = strings.TrimSpace(password)

		if name == "" || address == "" {
			ShowErrorModal(ctx, "Please fill in all required fields (*)!", form)
			return
		}

		if isModify && name != oldName {
			ctx.DeleteServer(oldName)
		}

		serverCfg := ctx.GetServer(name)
		serverCfg.Address = address
		serverCfg.Username = username
		serverCfg.Port = port
		if config.AuthMethodIdentityFile == authMethod {
			serverCfg.AuthMethod = authMethod
		}
		serverCfg.IdentityFile = identityFile
		if password != "" {
			err := keyring.Set("sshm", fmt.Sprintf("%s-%s", serverCfg.Address, serverCfg.Username), password)
			CheckError(&ctx, err)
		}

		CheckError(&ctx, ctx.SaveServer(name, serverCfg))

		ctx.UpdateTable("")
		ctx.Pages.RemovePage("write_server")
		ctx.App.SetFocus(ctx.Table)
	})

	form.SetBorder(true).SetTitle(title)

	ctx.Pages.AddPage("write_server", tview.NewGrid().
		SetRows(0, 19, 0).
		SetColumns(0, 75, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true), true, true)
	ctx.App.SetFocus(form)
}

func checkDefault(value, defaultValue string) string {
	if defaultValue == value {
		return ""
	}
	return value
}
