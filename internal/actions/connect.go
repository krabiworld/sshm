package actions

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/krabiworld/sshm/internal/app"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/utils"
	"github.com/zalando/go-keyring"
)

func Connect(ctx app.Context) {
	row, _ := ctx.Table.GetSelection()
	cell := ctx.Table.GetCell(row, 0)
	if cell == nil {
		return
	}

	var (
		server     = ctx.GetServer(cell.Text)
		binary     string
		args       = []string{}
		username   = utils.CheckDefault(server.Username, ctx.GetDefaults().Username)
		authMethod = utils.CheckDefault(server.AuthMethod, ctx.GetDefaults().AuthMethod)
		err        error
	)

	switch authMethod {
	case config.AuthMethodIdentityFile:
		binary, err = exec.LookPath("ssh")
	case config.AuthMethodPassword:
		binary, err = exec.LookPath("sshpass")
		args = append(args, "-d", "3", "ssh")
	}
	utils.CheckError(&ctx, err)

	args = append(args, fmt.Sprintf("%s@%s", username, server.Address))

	if port := server.Port; port != "" {
		args = append(args, "-p", port)
	}
	if identityFile := server.IdentityFile; identityFile != "" {
		args = append(args, "-i", identityFile)
	}

	ctx.App.Suspend(func() {
		r, w, err := os.Pipe()
		utils.CheckError(&ctx, err)

		cmd := exec.Command(binary, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.ExtraFiles = []*os.File{r}
		cmd.Start()

		if authMethod == config.AuthMethodPassword {
			password, err := keyring.Get("sshm", fmt.Sprintf("%s-%s", server.Address, server.Username))
			if err != nil {
				r.Close()
				w.Close()
				utils.CheckError(&ctx, err)
			}
			r.Close()
			_, err = w.Write([]byte(password + "\n"))
			w.Close()
		} else {
			r.Close()
			w.Close()
		}

		cmd.Wait()
	})
}
