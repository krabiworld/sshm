package ui

import (
	"fmt"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"github.com/krabiworld/sshm/internal/security"
)

func (m model) connectSSH(name string) tea.Cmd {
	var (
		server = m.config.Get(name)
		binary string
		args   = []string{}
		err    error
	)

	if server.HasPassword {
		binary, err = exec.LookPath("sshpass")
		args = append(args, "-d", "3", "ssh", "-o", "StrictHostKeyChecking=no")
	} else {
		binary, err = exec.LookPath("ssh")
	}
	if err != nil {
		panic(err)
	}

	args = append(args, fmt.Sprintf("%s@%s", server.Username, server.Address))

	if port := server.Port; port != "" {
		args = append(args, "-p", port)
	}
	if identityFile := server.IdentityFile; !server.HasPassword && identityFile != "" {
		args = append(args, "-i", identityFile)
	}

	cmd := exec.Command(binary, args...)

	if server.HasPassword {
		r, w, err := os.Pipe()
		if err != nil {
			panic(err)
		}

		cmd.ExtraFiles = []*os.File{r}

		password, err := security.GetPassword(name)
		if err != nil {
			r.Close()
			w.Close()
			panic(err)
		}

		go func() {
			defer w.Close()
			_, _ = w.Write([]byte(password))
		}()
	}

	return tea.ExecProcess(cmd, nil)
}

func (m model) copyId(name string) tea.Cmd {
	binary, err := exec.LookPath("ssh-copy-id")
	if err != nil {
		panic(err)
	}

	server := m.config.Get(name)

	var args []string

	if port := server.Port; port != "" {
		args = append(args, "-p", port)
	}

	args = append(args, fmt.Sprintf("%s@%s", server.Username, server.Address))

	return tea.ExecProcess(exec.Command(binary, args...), nil)
}
