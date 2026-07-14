package main

import (
	"flag"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/ui"
	"github.com/krabiworld/sshm/internal/utils"
)

func main() {
	configPath := flag.String("config", "~/.ssh/config.sshm.json", "")
	flag.Parse()

	cfg, err := config.New(utils.ExpandPath(*configPath))
	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(ui.NewModel(cfg))
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
