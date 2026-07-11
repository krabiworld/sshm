package utils

import "github.com/krabiworld/sshm/internal/app"

func CheckError(ctx *app.Context, err error) {
	if err != nil {
		if ctx != nil {
			ctx.App.Stop()
		}
		panic(err)
	}
}
