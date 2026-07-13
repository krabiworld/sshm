package ui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/security"
	"github.com/krabiworld/sshm/internal/utils"
	"github.com/muesli/cancelreader"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type sshCmd struct {
	serverName string
	server     config.Server
}

func (s *sshCmd) Run() error {
	clientConfig := &ssh.ClientConfig{
		User: s.server.Username,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			err := utils.GetKnownHosts()(hostname, remote, key)
			if err == nil {
				return nil
			}

			if keyErr, ok := errors.AsType[*knownhosts.KeyError](err); ok {
				if len(keyErr.Want) == 0 {
					fingerprint := ssh.FingerprintSHA256(key)

					fmt.Printf("The authenticity of host '%s' can't be established.\n", hostname)
					fmt.Printf("%s key fingerprint is: %s\n", key.Type(), fingerprint)
					fmt.Println("This key is not known by any other names.")
					fmt.Print("Are you sure you want to continue connecting (yes/no)? ")

					reader := bufio.NewReader(os.Stdin)
					answer, _ := reader.ReadString('\n')
					answer = strings.TrimSpace(strings.ToLower(answer))

					if answer == "yes" {
						err := utils.AddHostKey(hostname, key)
						if err != nil {
							fmt.Printf("Warning: Failed to add host to the list of known hosts: %v\n", err)
						}
						return nil
					}

					return fmt.Errorf("Host key verification failed.")
				}

				fmt.Println("Warning: remote host identification has changed!")
				return err
			}

			return err
		},
	}

	var auth ssh.AuthMethod

	switch s.server.AuthType {
	case config.AuthPassword:
		pw, err := security.GetPassword(s.serverName)
		if err != nil {
			panic(err)
		}
		auth = ssh.Password(pw)
	case config.AuthKey:
		var passphrase string
		if s.server.HasPassphrase {
			p, err := security.GetPassword(s.serverName)
			if err != nil {
				panic(err)
			}
			passphrase = p
		}
		auth = utils.GetAuthMethod(s.server.IdentityFile, passphrase)
	default:
		panic("Unknown auth method")
	}

	clientConfig.Auth = []ssh.AuthMethod{auth}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", s.server.Address, s.server.Port), clientConfig)
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

	return tea.Exec(&sshCmd{serverName: name, server: server}, nil)
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
