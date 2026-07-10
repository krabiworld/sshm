package actions

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/krabiworld/sshm/internal/app"
)

func CopyID(ctx app.Context) {
	row, _ := ctx.Table.GetSelection()
	cell := ctx.Table.GetCell(row, 0)
	if cell == nil {
		return
	}

	binary, err := exec.LookPath("ssh-copy-id")
	if err != nil {
		panic(err)
	}

	c := ctx.Config.Hosts[cell.Text]

	var args []string

	if port := c.Port; port != "" {
		args = append(args, "-p", port)
	}

	args = append(args, fmt.Sprintf("%s@%s", c.Username, c.Address))

	ctx.App.Suspend(func() {
		cmd := exec.Command(binary, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	})
}
