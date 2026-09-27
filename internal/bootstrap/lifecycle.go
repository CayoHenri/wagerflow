package bootstrap

import (
	"context"
	"log"

	"github.com/CayoHenri/wagerflow/internal/config"
	"go.uber.org/fx"
)

func registerLifecycle(lifecycle fx.Lifecycle, cfg *config.Config) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Printf("%s started | environment=%s | http_port=%d", cfg.App.Name, cfg.App.Env, cfg.HTTP.Port)
				return nil
			},

			OnStop: func(ctx context.Context) error {
				log.Printf("%s stopped", cfg.App.Name)
				return nil
			},
		},
	)
}
