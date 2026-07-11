package utils

import (
	"fmt"
	"strings"

	"github.com/krabiworld/sshm/internal/app"
	"github.com/rivo/tview"
	"github.com/zalando/go-keyring"
)

func WriteHost(ctx app.Context, oldHostname string) {
	isModify := oldHostname != ""
	title := "Create Host"

	var (
		hostname     = oldHostname
		address      string
		username     string
		port         string
		authMethod   = ConvertAuthMethodToInt(ctx.Config.Defaults.AuthMethod)
		identityFile string
		password     string
	)

	if isModify {
		title = "Modify Host"
		if cfgHost, ok := ctx.Config.Hosts[oldHostname]; ok {
			address = cfgHost.Address
			username = cfgHost.Username
			port = cfgHost.Port
			if cfgHost.AuthMethod != "" {
				authMethod = ConvertAuthMethodToInt(cfgHost.AuthMethod)
			}
			identityFile = cfgHost.IdentityFile
			password, _ = keyring.Get("sshm", fmt.Sprintf("%s-%s", cfgHost.Address, cfgHost.Username))
		}
	}

	form := tview.NewForm()
	form.AddInputField("Hostname *", hostname, 0, nil, func(text string) { hostname = text })
	form.AddInputField("Address *", address, 0, nil, func(text string) { address = text })
	form.AddInputField("Username *", username, 0, nil, func(text string) { username = text })
	form.AddDropDown("Auth method *", []string{"Identity file", "Password"}, authMethod, func(_ string, optionIndex int) { authMethod = optionIndex })

	if isModify {
		form.AddInputField("Port", port, 0, nil, func(text string) { port = text })
		form.AddInputField("Identity file", identityFile, 0, nil, func(text string) { identityFile = text })
		form.AddPasswordField("Password", password, 0, '*', func(text string) { password = text })
	} else {
		form.AddInputField(fmt.Sprintf("Port (%s)", ctx.Config.Defaults.Port), "", 0, nil, func(text string) { port = text })
		form.AddInputField(fmt.Sprintf("Identity file (%s)", ctx.Config.Defaults.IdentityFile), "", 0, nil, func(text string) { identityFile = text })
		form.AddPasswordField("Password", "", 0, '*', func(text string) { password = text })
	}

	form.AddButton("Save", func() {
		hostname = strings.TrimSpace(hostname)
		address = strings.TrimSpace(address)
		username = strings.TrimSpace(username)
		port = strings.TrimSpace(port)
		identityFile = strings.TrimSpace(identityFile)
		password = strings.TrimSpace(password)

		if hostname == "" || address == "" || username == "" || (authMethod != 0 && authMethod != 1) {
			ShowErrorModal(ctx, "Please fill in all required fields (*)!", form)
			return
		}

		if isModify && hostname != oldHostname {
			delete(ctx.Config.Hosts, oldHostname)
		}

		hostConfig := ctx.Config.Hosts[hostname]
		hostConfig.Address = address
		hostConfig.Username = username
		if port != "" {
			hostConfig.Port = port
		}
		if identityFile != "" {
			hostConfig.IdentityFile = identityFile
		}
		if password != "" {
			if err := keyring.Set("sshm", fmt.Sprintf("%s-%s", hostConfig.Address, hostConfig.Username), password); err != nil {
				ctx.App.Stop()
				panic(err)
			}
		}
		ctx.Config.Hosts[hostname] = hostConfig

		if err := ctx.Config.Write(ctx.ConfigPath); err != nil {
			ctx.App.Stop()
			fmt.Printf("Error while initializing config: %v\n", err)
			return
		}

		ctx.UpdateTable("")
		ctx.Pages.RemovePage("write_host")
		ctx.App.SetFocus(ctx.Table)
	})

	form.SetBorder(true).SetTitle(title)

	ctx.Pages.AddPage("write_host", tview.NewGrid().
		SetRows(0, 17, 0).
		SetColumns(0, 65, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true), true, true)
	ctx.App.SetFocus(form)
}
