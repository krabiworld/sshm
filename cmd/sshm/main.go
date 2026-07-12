package main

import (
	"flag"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/ui"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configPath := flag.String("config", filepath.Join(homeDir, ".ssh", "config.sshm.json"), "")
	flag.Parse()

	var cfg config.Config

	_, err = os.Stat(*configPath)
	if os.IsNotExist(err) {
		cfg.Write(*configPath)
	}

	cfg.Read(*configPath)

	p := tea.NewProgram(ui.NewModel(cfg, *configPath))
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
