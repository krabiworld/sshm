package app

import (
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/rivo/tview"
)

type Context struct {
	Config     *config.Config
	ConfigPath string
	App        *tview.Application
	Pages      *tview.Pages
	Table      *tview.Table
}

func (ctx *Context) UpdateTable(filter string) {
	ctx.Table.Clear()

	headers := []string{"Hostname", "Address", "Username", "Port", "Identity"}
	for col, text := range headers {
		cell := tview.NewTableCell(text).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetExpansion(1)

		ctx.Table.SetCell(0, col, cell)
	}

	var hostnames []string
	for host := range ctx.Config.Hosts {
		hostnames = append(hostnames, host)
	}
	sort.Strings(hostnames)

	rowIdx := 1
	for _, host := range hostnames {
		cfgHost := ctx.Config.Hosts[host]

		address := cfgHost.Address
		user := cfgHost.Username
		port := cfgHost.Port
		if port == "" {
			port = ctx.Config.Defaults.Port + " (d)"
		}
		authMethod := cfgHost.AuthMethod
		if authMethod == "" {
			authMethod = ctx.Config.Defaults.AuthMethod
		}
		identity := cfgHost.IdentityFile
		switch authMethod {
		case config.AuthMethodIdentityFile:
			if identity == "" {
				identity = ctx.Config.Defaults.IdentityFile + " (d)"
			}
		case config.AuthMethodPassword:
			identity = "Password (d)"
		}

		if filter != "" {
			if !strings.HasPrefix(host, filter) && !strings.HasPrefix(address, filter) {
				continue
			}
		}

		ctx.Table.SetCell(rowIdx, 0, tview.NewTableCell(host))
		ctx.Table.SetCell(rowIdx, 1, tview.NewTableCell(address))
		ctx.Table.SetCell(rowIdx, 2, tview.NewTableCell(user))
		ctx.Table.SetCell(rowIdx, 3, tview.NewTableCell(port))
		ctx.Table.SetCell(rowIdx, 4, tview.NewTableCell(identity))

		rowIdx++
	}
}
