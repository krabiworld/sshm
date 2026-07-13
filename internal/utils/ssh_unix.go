//go:build !windows

package utils

import "net"

func dialNamedPipe(socket string) (net.Conn, error) {
	return net.Dial("unix", socket)
}
