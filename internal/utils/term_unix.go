//go:build !windows

package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func WatchTermResize(_ uintptr) (chan struct{}, func()) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGWINCH)

	outChan := make(chan struct{}, 1)

	go func() {
		for range sigch {
			outChan <- struct{}{}
		}
	}()

	cleanup := func() {
		signal.Stop(sigch)
		close(outChan)
	}

	return outChan, cleanup
}
