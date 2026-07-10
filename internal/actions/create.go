package actions

import (
	"github.com/krabiworld/sshm/internal/app"
	"github.com/krabiworld/sshm/internal/utils"
)

func CreateHost(ctx app.Context) {
	utils.WriteHost(ctx, "")
}
