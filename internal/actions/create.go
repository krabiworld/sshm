package actions

import (
	"github.com/krabiworld/sshm/internal/app"
	"github.com/krabiworld/sshm/internal/utils"
)

func Create(ctx app.Context) {
	utils.WriteServer(ctx, "")
}
