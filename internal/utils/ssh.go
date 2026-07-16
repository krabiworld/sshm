package utils

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const (
	sshDirPerm        os.FileMode = 0700
	sshKnownHostsPerm os.FileMode = 0600

	sshAgentEnv = "SSH_AUTH_SOCK"
)

func GetAuthMethod(keyPath string, password string) ssh.AuthMethod {
	keyBytes, err := os.ReadFile(ExpandPath(keyPath))
	if err != nil {
		panic(err)
	}

	var signer ssh.Signer
	if password != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(password))
	} else {
		signer, err = ssh.ParsePrivateKey(keyBytes)
	}

	if err != nil {
		panic(err)
	}

	return ssh.PublicKeys(signer)
}

func GetKnownHosts(path string) ssh.HostKeyCallback {
	path = ExpandPath(path)

	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, sshKnownHostsPerm)
	if err != nil {
		panic(err)
	}
	f.Close()

	callback, err := knownhosts.New(path)
	if err != nil {
		panic(err)
	}

	return callback
}

func AddHostKey(path, hostname string, key ssh.PublicKey) error {
	f, err := os.OpenFile(ExpandPath(path), os.O_APPEND|os.O_CREATE|os.O_WRONLY, sshKnownHostsPerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(knownhosts.Line([]string{hostname}, key) + "\n")
	return err
}

func CreateSshDir() error {
	return os.MkdirAll(filepath.Dir(ExpandPath("~/.ssh")), sshDirPerm)
}

func GetAgentDial() (net.Conn, error) {
	socket := os.Getenv(sshAgentEnv)

	if socket == "" && runtime.GOOS == "windows" {
		socket = `\\.\pipe\openssh-ssh-agent`
	}

	if socket == "" {
		return nil, fmt.Errorf(sshAgentEnv + " is not defined")
	}

	return dialNamedPipe(socket)
}
