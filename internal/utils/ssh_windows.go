//go:build windows

package utils

import (
	"net"

	"github.com/Microsoft/go-winio"
)

func dialNamedPipe(socket string) (net.Conn, error) {
	return winio.DialPipe(socket, nil)
}
