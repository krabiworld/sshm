package actions

import (
	"github.com/krabiworld/sshm/internal/app"
	"github.com/krabiworld/sshm/internal/utils"
)

func ModifyHost(ctx app.Context) {
	row, _ := ctx.Table.GetSelection()
	cell := ctx.Table.GetCell(row, 0)
	if cell == nil {
		return
	}
	utils.WriteHost(ctx, cell.Text)
}
