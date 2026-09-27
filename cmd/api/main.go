package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

type App struct {
	Name string
}

func NewApp() *App {
	return &App{
		Name: "wagerflow",
	}
}

func RegisterLifecycle(lifecycle fx.Lifecycle, app *App) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Printf("%s started\n", app.Name)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				fmt.Printf("%s stopped\n", app.Name)
				return nil
			},
		},
	)
}

func main() {
	fx.New(
		fx.Provide(
			NewApp,
		),
		fx.Invoke(
			RegisterLifecycle,
		),
	).Run()
}
