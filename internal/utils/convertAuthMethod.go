package utils

import "github.com/krabiworld/sshm/internal/config"

var intToAuthMethod = map[int]string{
	0: config.AuthMethodIdentityFile,
	1: config.AuthMethodPassword,
}

var authMethodToInt = make(map[string]int)

func init() {
	for k, v := range intToAuthMethod {
		authMethodToInt[v] = k
	}
}

func ConvertAuthMethodToInt(authMethod string) int {
	return authMethodToInt[authMethod]
}

func ConvertIntToAuthMethod(key int) string {
	return intToAuthMethod[key]
}
