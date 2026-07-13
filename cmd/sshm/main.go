package main

import (
	"flag"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/ui"
	"github.com/krabiworld/sshm/internal/utils"
)

func main() {
	configPath := flag.String("config", "~/.ssh/config.sshm.json", "")
	flag.Parse()

	*configPath = utils.ExpandPath(*configPath)

	var cfg config.Config

	_, err := os.Stat(*configPath)
	if os.IsNotExist(err) {
		cfg.Write(*configPath)
	}

	cfg.Read(*configPath)

	p := tea.NewProgram(ui.NewModel(cfg, *configPath))
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
