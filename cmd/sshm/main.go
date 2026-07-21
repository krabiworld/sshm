package main

import (
	"flag"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/ui"
	"github.com/krabiworld/sshm/internal/utils"
)

var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	showVersion := flag.Bool("version", false, "")
	showDebug := flag.Bool("debug", false, "")
	configPath := flag.String("config", "~/.ssh/config.sshm.json", "")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\nCommit: %s\nBuild Date: %s\n", version, commit, date)
		return
	}

	if *showDebug {
		printDebug()
		return
	}

	cfg, err := config.New(utils.ExpandPath(*configPath))
	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(ui.NewModel(cfg))
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
