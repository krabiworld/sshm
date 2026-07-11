package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/krabiworld/sshm/internal/actions"
	"github.com/krabiworld/sshm/internal/app"
	"github.com/krabiworld/sshm/internal/utils"
	"github.com/rivo/tview"
)

func main() {
	homeDir, err := os.UserHomeDir()
	utils.CheckError(nil, err)

	configPath := flag.String("config", filepath.Join(homeDir, ".ssh", "config.sshm.json"), "")
	flag.Parse()

	ctx := app.NewContext(
		*configPath,
		tview.NewApplication(),
		tview.NewPages(),
		tview.NewTable().SetSelectable(true, false),
	)

	_, err = os.Stat(*configPath)
	if os.IsNotExist(err) {
		utils.CheckError(&ctx, ctx.WriteConfig())
	}

	utils.CheckError(&ctx, ctx.ReadConfig())

	ctx.ApplyTheme()

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
			actions.Create(ctx)
			return nil
		}

		if event.Key() == tcell.KeyCtrlM && ctx.App.GetFocus() == ctx.Table {
			actions.Modify(ctx)
			return nil
		}

		if event.Key() == tcell.KeyCtrlD && ctx.App.GetFocus() == ctx.Table {
			actions.Delete(ctx)
			return nil
		}

		if event.Key() == tcell.KeyCtrlI && ctx.App.GetFocus() == ctx.Table {
			actions.CopyID(ctx)
			return nil
		}

		return event
	})

	utils.CheckError(&ctx, ctx.App.SetRoot(ctx.Pages, true).Run())
}
