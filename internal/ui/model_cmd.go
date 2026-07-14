package ui

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/security"
	"github.com/krabiworld/sshm/internal/utils"
	"github.com/muesli/cancelreader"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

type hostKeyRequiredError struct {
	hostname string
	key      ssh.PublicKey
}

func (e *hostKeyRequiredError) Error() string {
	return "host key confirmation required"
}

type sshConnectedMsg struct {
	client       *ssh.Client
	session      *ssh.Session
	cancelReader cancelreader.CancelReader
}

func (s *sshConnectedMsg) Run() error {
	defer s.client.Close()
	defer s.session.Close()
	defer s.cancelReader.Cancel()

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
	if err := s.session.RequestPty("xterm-256color", h, w, modes); err != nil {
		return err
	}
	if err := s.session.Shell(); err != nil {
		return err
	}

	resizeChan, cleanup := utils.WatchTermResize(fd)
	defer cleanup()

	go func() {
		for range resizeChan {
			if w, h, err := term.GetSize(fd); err == nil {
				_ = s.session.WindowChange(h, w)
			}
		}
	}()

	return s.session.Wait()
}

func (s *sshConnectedMsg) SetStdin(r io.Reader)  {}
func (s *sshConnectedMsg) SetStdout(w io.Writer) {}
func (s *sshConnectedMsg) SetStderr(w io.Writer) {}

func (m model) dialSsh(name string) tea.Cmd {
	return func() tea.Msg {
		server := m.config.Get(name)

		clientConfig := &ssh.ClientConfig{
			User: server.Username,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				err := utils.GetKnownHosts()(hostname, remote, key)
				if err == nil {
					return nil
				}

				if keyErr, ok := errors.AsType[*knownhosts.KeyError](err); ok {
					if len(keyErr.Want) == 0 {
						return &hostKeyRequiredError{hostname, key}
					}

					return errors.New("Warning: remote host identification has changed!")
				}

				return err
			},
		}

		var auth ssh.AuthMethod

		switch server.AuthType {
		case config.AuthPassword:
			pw, err := security.GetPassword(name)
			if err != nil {
				return errMsg{fmt.Errorf("Password retrieval error: %w", err)}
			}
			auth = ssh.Password(pw)
		case config.AuthKey:
			var passphrase string
			if server.HasPassphrase {
				p, err := security.GetPassword(name)
				if err != nil {
					return errMsg{fmt.Errorf("Passphrase retrieval error: %w", err)}
				}
				passphrase = p
			}
			auth = utils.GetAuthMethod(server.IdentityFile, passphrase)
		case config.AuthAgent:
			agentDial, err := utils.GetAgentDial()
			if err != nil {
				return errMsg{fmt.Errorf("Cannot connect to ssh-agent: %w", err)}
			}

			agentClient := agent.NewClient(agentDial)

			auth = ssh.PublicKeysCallback(agentClient.Signers)
		default:
			return errMsg{fmt.Errorf("Unknown auth method: %s", server.AuthType)}
		}

		clientConfig.Auth = []ssh.AuthMethod{auth}

		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", server.Address, server.Port), clientConfig)
		if err != nil {
			return errMsg{err}
		}

		session, err := client.NewSession()
		if err != nil {
			client.Close()
			return errMsg{err}
		}

		cancelableStdin, err := cancelreader.NewReader(os.Stdin)
		if err != nil {
			client.Close()
			session.Close()
			return errMsg{err}
		}

		session.Stdin = cancelableStdin
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		return &sshConnectedMsg{
			client:       client,
			session:      session,
			cancelReader: cancelableStdin,
		}
	}
}

func (m model) runSshSession(msg *sshConnectedMsg) tea.Cmd {
	return tea.Exec(msg, func(err error) tea.Msg {
		if _, ok := errors.AsType[*ssh.ExitError](err); ok {
			return nil
		}
		if err != nil {
			return errMsg{err}
		}
		return nil
	})
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
