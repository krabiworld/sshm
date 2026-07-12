package security

import "github.com/zalando/go-keyring"

const service = "sshm"

func GetPassword(name string) (string, error) {
	return keyring.Get(service, name)
}

func SetPassword(name, password string) error {
	return keyring.Set(service, name, password)
}

func DeletePassword(name string) error {
	return keyring.Delete(service, name)
}
