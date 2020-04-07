package console

import (
	"context"
	"github.com/urfave/cli"
)

type App struct {
	cli.App

	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewApp creates a new console Application with some reasonable default settings
func NewApp() *App {
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &App{
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
}
