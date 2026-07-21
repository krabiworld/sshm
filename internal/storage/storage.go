package storage

import "github.com/krabiworld/sshm/internal/config"

type Storage interface {
	GetPassword(name string) (string, error)
	SetPassword(name, password string) error
	DeletePassword(name string) error
}

type storage struct {
	cfg       *config.Config
	keychain  Storage
	plaintext Storage
}

func NewStorage(cfg *config.Config) Storage {
	return &storage{cfg, keychain{}, &plaintext{cfg}}
}

func (s *storage) GetPassword(name string) (string, error) {
	return s.impl(name).GetPassword(name)
}

func (s *storage) SetPassword(name, password string) error {
	return s.impl(name).SetPassword(name, password)
}

func (s *storage) DeletePassword(name string) error {
	return s.impl(name).DeletePassword(name)
}

func (s *storage) impl(name string) Storage {
	storageType := s.cfg.Get(name).PasswordStorageType

	switch storageType {
	case config.StorageKeychain:
		return s.keychain
	case config.StoragePlaintext:
		return s.plaintext
	default:
		panic("unknown storage type")
	}
}
