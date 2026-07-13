package ui

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/security"
	"github.com/krabiworld/sshm/internal/utils"
	"github.com/muesli/cancelreader"
	"golang.org/x/crypto/ssh"
)

type sshCmd struct {
	config  *ssh.ClientConfig
	address string
	port    string
}

func (s *sshCmd) Run() error {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", s.address, s.port), s.config)
	if err != nil {
		panic(err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		panic(err)
	}

	defer client.Close()
	defer session.Close()

	cancelableStdin, err := cancelreader.NewReader(os.Stdin)
	if err != nil {
		panic(err)
	}

	defer cancelableStdin.Cancel()

	session.Stdin = cancelableStdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	fd := os.Stdin.Fd()
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, oldState)

	w, h, err := term.GetSize(fd)
	if err != nil {
		w, h = 80, 24
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm-256color", h, w, modes); err != nil {
		return fmt.Errorf("request for pty failed: %w", err)
	}

	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	return session.Wait()
}

func (s *sshCmd) SetStdin(r io.Reader)  {}
func (s *sshCmd) SetStdout(w io.Writer) {}
func (s *sshCmd) SetStderr(w io.Writer) {}

func (m model) connectSsh(name string) tea.Cmd {
	server := m.config.Get(name)

	c := &sshCmd{
		config: &ssh.ClientConfig{
			User:            server.Username,
			HostKeyCallback: utils.GetKnownHosts(),
		},
		address: server.Address,
		port:    server.Port,
	}

	var auth ssh.AuthMethod

    switch server.AuthType {
    case config.AuthPassword:
        pw, err := security.GetPassword(name)
        if err != nil { panic(err) }
        auth = ssh.Password(pw)
    case config.AuthKey:
        var passphrase string
        if server.HasPassphrase {
            p, err := security.GetPassword(name)
            if err != nil { panic(err) }
            passphrase = p
        }
        auth = utils.GetAuthMethod(server.IdentityFile, passphrase)
    default:
        panic("Unknown auth type")
    }

    c.config.Auth = []ssh.AuthMethod{auth}

	return tea.Exec(c, nil)
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
