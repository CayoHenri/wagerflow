package bootstrap

import (
	"context"
	"log"

	"go.uber.org/fx"
)

func registerLifecycle(lifecycle fx.Lifecycle) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Println("wagerflow started")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Println("wagerflow stopped")
				return nil
			},
		},
	)
}
