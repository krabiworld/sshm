package utils

import (
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
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

func GetKnownHosts() ssh.HostKeyCallback {
	callback, err := knownhosts.New(ExpandPath("~/.ssh/known_hosts"))
	if err != nil {
		panic(err)
	}

	return callback
}
