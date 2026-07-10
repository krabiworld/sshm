package actions

import (
	"fmt"

	"github.com/krabiworld/sshm/internal/app"
	"github.com/rivo/tview"
)

func DeleteHost(ctx app.Context) {
	modal := tview.NewModal().
		SetText("Are you sure you want to delete the host?").
		AddButtons([]string{"No", "Yes"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel != "Yes" {
				ctx.Pages.RemovePage("delete_modal")
				return
			}

			row, _ := ctx.Table.GetSelection()
			if row == 0 {
				return
			}
			cell := ctx.Table.GetCell(row, 0)
			if cell == nil {
				return
			}
			delete(ctx.Config.Hosts, cell.Text)
			if err := ctx.Config.Write(ctx.ConfigPath); err != nil {
				ctx.App.Stop()
				fmt.Printf("Error while initializing config: %v\n", err)
			}

			ctx.Table.RemoveRow(row)
			if row >= ctx.Table.GetRowCount() && ctx.Table.GetRowCount() > 1 {
				ctx.Table.Select(ctx.Table.GetRowCount()-1, 0)
			}

			ctx.Pages.RemovePage("delete_modal")
		})
	ctx.Pages.AddPage("delete_modal", modal, true, true)
}
