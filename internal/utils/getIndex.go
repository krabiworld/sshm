package utils

import (
	"github.com/krabiworld/sshm/internal/config"
)

var (
	ThemeOrder = []string{
		config.ThemeDark,
		config.ThemeLight,
		config.ThemeTransparent,
	}
	AuthMethodOrder = []string{
		config.AuthMethodIdentityFile,
		config.AuthMethodPassword,
	}
)

func GetIndex[T comparable](slice []T, target T) int {
	for i, item := range slice {
		if item == target {
			return i
		}
	}
	return -1
}

func GetByIndex[T any](slice []T, index int) T {
	return slice[index]
}
