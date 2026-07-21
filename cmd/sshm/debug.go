package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/x/term"
)

func getKerboardProtocol() (string, error) {
	fd := os.Stdout.Fd()

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, oldState)

	// \x1b[?u  - Kitty Keyboard Protocol
	// \x1b[>4m - xterm modifyOtherKeys
	os.Stdout.WriteString("\x1b[?u\x1b[>4m")

	responseChan := make(chan string)
	go func() {
		buf := make([]byte, 1024)
		n, err := os.Stdin.Read(buf)
		if err == nil && n > 0 {
			responseChan <- string(buf[:n])
		}
	}()

	var response string
	select {
	case response = <-responseChan:
	case <-time.After(100 * time.Millisecond):
		response = "timeout"
	}

	term.Restore(fd, oldState)

	if response == "timeout" {
		return "Legacy (timeout)", nil
	}

	if strings.Contains(response, "\x1b[?") && strings.Contains(response, "u") {
		return "Kitty Keyboard Protocol", nil
	} else if strings.Contains(response, "\x1b[>") {
		return "xterm modifyOtherKeys Protocol", nil
	} else {
		return "Legacy", nil
	}
}

func getTerminalSize() (width, height int, err error) {
	fd := os.Stdout.Fd()

	if !term.IsTerminal(fd) {
		return 0, 0, fmt.Errorf("stdout is not a terminal")
	}

	width, height, err = term.GetSize(fd)
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func printDebug() {
	envs := []string{
		"TERM",
		"TERM_PROGRAM",
		"TERM_PROGRAM_VERSION",
		"COLORTERM",
		"SHELL",
		"LANG",
		"LC_ALL",
		"LC_CTYPE",
		"TMUX",
		"STY",
	}

	builder := strings.Builder{}

	for _, key := range envs {
		val, ok := os.LookupEnv(key)
		if !ok {
			val = "<unset>"
		}
		fmt.Fprintf(&builder, "%s: %s\n", key, val)
	}

	keyboardProtocol, err := getKerboardProtocol()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(&builder, "Keyboard protocol: %s\n", keyboardProtocol)

	w, h, err := getTerminalSize()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(&builder, "Terminal size: %dx%d\n", w, h)

	fmt.Print(builder.String())
}
