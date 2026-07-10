package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/krabiworld/sshm/internal/actions"
	"github.com/krabiworld/sshm/internal/app"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/rivo/tview"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configPath := flag.String("config", filepath.Join(homeDir, ".ssh", "config.sshm.json"), "")
	flag.Parse()

	ctx := app.Context{
		Config:     &config.Config{},
		ConfigPath: *configPath,
		App:        tview.NewApplication(),
		Pages:      tview.NewPages(),
		Table:      tview.NewTable().SetSelectable(true, false),
	}

	_, err = os.Stat(*configPath)
	if os.IsNotExist(err) {
		if err := ctx.Config.Write(ctx.ConfigPath); err != nil {
			ctx.App.Stop()
			fmt.Printf("Error while initializing config: %v\n", err)
			return
		}
	}

	if err := ctx.Config.Read(ctx.ConfigPath); err != nil {
		panic(err)
	}

	// Footer
	footer := tview.NewTextView().SetText("^F Search | ^A Add | ^M Modify | ^D Delete | ^I Copy ID | ^X Settings")
	searchField := tview.NewInputField().SetLabel("Search: ")

	// Layouts
	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ctx.Table, 0, 1, true).
		AddItem(footer, 1, 0, false)

	layoutWithPadding := tview.NewGrid().
		SetColumns(1, 0, 1).
		SetRows(0).
		AddItem(mainLayout, 0, 1, 1, 1, 0, 0, true)

	openSearch := func() {
		mainLayout.RemoveItem(footer)
		mainLayout.AddItem(searchField, 1, 0, true)

		ctx.App.SetFocus(searchField)
	}
	closeSearch := func() {
		mainLayout.RemoveItem(searchField)
		mainLayout.AddItem(footer, 1, 0, false)

		ctx.App.SetFocus(ctx.Table)
	}

	searchField.SetChangedFunc(func(text string) {
		ctx.UpdateTable(text)
	})
	searchField.SetDoneFunc(func(key tcell.Key) {
		closeSearch()
	})

	ctx.Pages.AddPage("main", layoutWithPadding, true, true)

	ctx.UpdateTable("")

	ctx.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			if ctx.App.GetFocus() == searchField {
				closeSearch()
				return nil
			}

			frontPage, _ := ctx.Pages.GetFrontPage()
			if frontPage != "main" && frontPage != "" {
				ctx.Pages.RemovePage(frontPage)
				ctx.App.SetFocus(ctx.Table)
				return nil
			}

			return nil
		}

		if event.Key() == tcell.KeyEnter && ctx.App.GetFocus() == ctx.Table {
			actions.Connect(ctx)
			return nil
		}

		if event.Key() == tcell.KeyCtrlX && ctx.App.GetFocus() == ctx.Table {
			actions.Settings(ctx)
			return nil
		}

		if event.Key() == tcell.KeyCtrlF && ctx.App.GetFocus() == ctx.Table {
			openSearch()
			return nil
		}

		if event.Key() == tcell.KeyCtrlA && ctx.App.GetFocus() == ctx.Table {
			actions.CreateHost(ctx)
			return nil
		}

		if event.Key() == tcell.KeyCtrlM && ctx.App.GetFocus() == ctx.Table {
			actions.ModifyHost(ctx)
			return nil
		}

		if event.Key() == tcell.KeyCtrlD && ctx.App.GetFocus() == ctx.Table {
			actions.DeleteHost(ctx)
			return nil
		}

		if event.Key() == tcell.KeyCtrlI && ctx.App.GetFocus() == ctx.Table {
			actions.CopyID(ctx)
			return nil
		}

		return event
	})

	if err := ctx.App.SetRoot(ctx.Pages, true).Run(); err != nil {
		panic(err)
	}
}
