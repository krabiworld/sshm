package actions

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/krabiworld/sshm/internal/app"
)

func Connect(ctx app.Context) {
	row, _ := ctx.Table.GetSelection()
	cell := ctx.Table.GetCell(row, 0)
	if cell == nil {
		return
	}

	binary, err := exec.LookPath("ssh")
	if err != nil {
		panic(err)
	}

	c := ctx.Config.Hosts[cell.Text]

	var args []string
	if ctx.Config.Settings.CloseAfterConnection {
		args = []string{"ssh", fmt.Sprintf("%s@%s", c.Username, c.Address)}
	} else {
		args = []string{fmt.Sprintf("%s@%s", c.Username, c.Address)}
	}

	if port := c.Port; port != "" {
		args = append(args, "-p", port)
	}
	if identityFile := c.IdentityFile; identityFile != "" {
		args = append(args, "-i", identityFile)
	}

	if ctx.Config.Settings.CloseAfterConnection {
		ctx.App.Stop()
		env := os.Environ()
		if err := syscall.Exec(binary, args, env); err != nil {
			panic(err)
		}
	} else {
		ctx.App.Suspend(func() {
			cmd := exec.Command(binary, args...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		})
	}
}
