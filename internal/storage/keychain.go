package storage

import "github.com/zalando/go-keyring"

const service = "sshm"

type keychain struct{}

func (keychain) GetPassword(name string) (string, error) {
	return keyring.Get(service, name)
}

func (keychain) SetPassword(name, password string) error {
	return keyring.Set(service, name, password)
}

func (keychain) DeletePassword(name string) error {
	return keyring.Delete(service, name)
}
