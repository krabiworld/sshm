package actions

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/krabiworld/sshm/internal/app"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/zalando/go-keyring"
)

func Connect(ctx app.Context) {
	row, _ := ctx.Table.GetSelection()
	cell := ctx.Table.GetCell(row, 0)
	if cell == nil {
		return
	}

	host := ctx.Config.Hosts[cell.Text]

	var (
		binary string
		args   = []string{}
		err    error
	)

	authMethod := host.AuthMethod
	if authMethod == "" {
		authMethod = ctx.Config.Defaults.AuthMethod
	}

	switch authMethod {
	case config.AuthMethodIdentityFile:
		binary, err = exec.LookPath("ssh")
	case config.AuthMethodPassword:
		binary, err = exec.LookPath("sshpass")
		args = append(args, "-d", "3", "ssh")
	}
	if err != nil {
		ctx.App.Stop()
		panic(err)
	}

	args = append(args, fmt.Sprintf("%s@%s", host.Username, host.Address))

	if port := host.Port; port != "" {
		args = append(args, "-p", port)
	}
	if identityFile := host.IdentityFile; identityFile != "" {
		args = append(args, "-i", identityFile)
	}

	ctx.App.Suspend(func() {
		r, w, err := os.Pipe()
		if err != nil {
			panic(err)
		}

		cmd := exec.Command(binary, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.ExtraFiles = []*os.File{r}
		cmd.Start()

		if authMethod == config.AuthMethodPassword {
			password, err := keyring.Get("sshm", fmt.Sprintf("%s-%s", host.Address, host.Username))
			if err != nil {
				r.Close()
				w.Close()
				panic(err)
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
