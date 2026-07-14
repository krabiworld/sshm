//go:build windows

package utils

import (
	"time"

	"github.com/charmbracelet/x/term"
)

func WatchTermResize(fd uintptr) (chan struct{}, func()) {
	outChan := make(chan struct{}, 1)
	done := make(chan struct{})

	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()

		lastW, lastH, _ := term.GetSize(fd)

		for {
			select {
			case <-ticker.C:
				w, h, err := term.GetSize(fd)
				if err == nil && (w != lastW || h != lastH) {
					lastW, lastH = w, h
					select {
					case outChan <- struct{}{}:
					default:
					}
				}
			case <-done:
				return
			}
		}
	}()

	cleanup := func() {
		close(done)
		close(outChan)
	}

	return outChan, cleanup
}
