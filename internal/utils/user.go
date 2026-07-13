package utils

import "os/user"

func GetCurrentUsername() string {
	usr, _ := user.Current()
	return usr.Username
}
