package app

import (
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/rivo/tview"
)

type Context struct {
	config     *config.Config
	configPath string
	App        *tview.Application
	Pages      *tview.Pages
	Table      *tview.Table
}

func NewContext(configPath string, app *tview.Application, pages *tview.Pages, table *tview.Table) Context {
	return Context{
		config:     &config.Config{},
		configPath: configPath,
		App:        app,
		Pages:      pages,
		Table:      table,
	}
}

func (ctx *Context) GetApplication() config.Application {
	return ctx.config.Application
}

func (ctx *Context) SaveApplication(app config.Application) error {
	err := ctx.config.SaveApplication(app, ctx.configPath)
	if err != nil {
		return err
	}

	ctx.ApplyTheme()
	ctx.App.ForceDraw()

	return nil
}

func (ctx *Context) ApplyTheme() {
	switch ctx.GetApplication().Theme {
	case config.ThemeDark:
		tview.Styles = tview.Theme{
			PrimitiveBackgroundColor:    tcell.ColorBlack,
			ContrastBackgroundColor:     tcell.ColorDarkSlateGray,
			MoreContrastBackgroundColor: tcell.ColorDarkSlateGray,
			BorderColor:                 tcell.ColorWhite,
			TitleColor:                  tcell.ColorWhite,
			GraphicsColor:               tcell.ColorWhite,
			PrimaryTextColor:            tcell.ColorWhite,
			SecondaryTextColor:          tcell.ColorYellow,
			TertiaryTextColor:           tcell.ColorGreen,
			InverseTextColor:            tcell.ColorBlue,
			ContrastSecondaryTextColor:  tcell.ColorNavy,
		}
	case config.ThemeLight:
		tview.Styles = tview.Theme{
			PrimitiveBackgroundColor:    tcell.ColorWhite,
			ContrastBackgroundColor:     tcell.ColorLightGray,
			MoreContrastBackgroundColor: tcell.ColorDarkGray,
			BorderColor:                 tcell.ColorBlack,
			TitleColor:                  tcell.ColorBlue,
			GraphicsColor:               tcell.ColorBlack,
			PrimaryTextColor:            tcell.ColorBlack,
			SecondaryTextColor:          tcell.ColorDarkBlue,
			TertiaryTextColor:           tcell.ColorDarkGreen,
			InverseTextColor:            tcell.ColorWhite,
			ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
		}
	case config.ThemeTransparent:
		tview.Styles = tview.Theme{
			PrimitiveBackgroundColor:    tcell.ColorDefault,
			ContrastBackgroundColor:     tcell.ColorDarkSlateGray,
			MoreContrastBackgroundColor: tcell.ColorDarkSlateGray,
			BorderColor:                 tcell.ColorWhite,
			TitleColor:                  tcell.ColorWhite,
			GraphicsColor:               tcell.ColorWhite,
			PrimaryTextColor:            tcell.ColorWhite,
			SecondaryTextColor:          tcell.ColorYellow,
			TertiaryTextColor:           tcell.ColorGreen,
			InverseTextColor:            tcell.ColorBlue,
			ContrastSecondaryTextColor:  tcell.ColorNavy,
		}
	}
	ctx.Table.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	ctx.Table.SetSelectedStyle(tcell.StyleDefault.Background(tview.Styles.MoreContrastBackgroundColor))
}

func (ctx *Context) GetServer(name string) config.Server {
	return ctx.config.Get(name)
}

func (ctx *Context) SaveServer(name string, server config.Server) error {
	err := ctx.config.Save(name, server, ctx.configPath)
	if err != nil {
		return err
	}

	ctx.UpdateTable("")

	return nil
}

func (ctx *Context) DeleteServer(name string) error {
	err := ctx.config.Delete(name, ctx.configPath)
	if err != nil {
		return err
	}

	ctx.UpdateTable("")

	return nil
}

func (ctx *Context) GetDefaults() config.Defaults {
	return ctx.config.Defaults
}

func (ctx *Context) SaveDefaults(def config.Defaults) error {
	return ctx.config.SaveDefaults(def, ctx.configPath)
}

func (ctx *Context) ReadConfig() error {
	return ctx.config.Read(ctx.configPath)
}

func (ctx *Context) WriteConfig() error {
	return ctx.config.Write(ctx.configPath)
}

func (ctx *Context) UpdateTable(filter string) {
	ctx.Table.Clear()

	headers := []string{"Name", "Address", "Username", "Port", "Identity"}
	for col, text := range headers {
		cell := tview.NewTableCell(text).
			SetTextColor(tview.Styles.ContrastSecondaryTextColor).
			SetSelectable(false).
			SetExpansion(1)

		ctx.Table.SetCell(0, col, cell)
	}

	var servers []string
	for server := range ctx.config.Servers {
		servers = append(servers, server)
	}
	sort.Strings(servers)

	rowIdx := 1
	for _, server := range servers {
		serverCfg := ctx.config.Servers[server]

		address := serverCfg.Address
		user := serverCfg.Username
		if user == "" {
			user = ctx.GetDefaults().Username + " (d)"
		}
		port := serverCfg.Port
		if port == "" {
			port = ctx.GetDefaults().Port + " (d)"
		}
		authMethod := serverCfg.AuthMethod
		if authMethod == "" {
			authMethod = ctx.GetDefaults().AuthMethod
		}
		identity := serverCfg.IdentityFile
		switch authMethod {
		case config.AuthMethodIdentityFile:
			if identity == "" {
				identity = ctx.GetDefaults().IdentityFile
			}
		case config.AuthMethodPassword:
			identity = "Password"
		}

		if serverCfg.AuthMethod == "" {
			identity += " (d)"
		}

		if filter != "" {
			if !strings.Contains(server, filter) && !strings.Contains(address, filter) {
				continue
			}
		}

		ctx.Table.SetCell(rowIdx, 0, tview.NewTableCell(server))
		ctx.Table.SetCell(rowIdx, 1, tview.NewTableCell(address))
		ctx.Table.SetCell(rowIdx, 2, tview.NewTableCell(user))
		ctx.Table.SetCell(rowIdx, 3, tview.NewTableCell(port))
		ctx.Table.SetCell(rowIdx, 4, tview.NewTableCell(identity))

		rowIdx++
	}
}
