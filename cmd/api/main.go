package main

import (
	"github.com/CayoHenri/wagerflow/internal/bootstrap"
	"github.com/CayoHenri/wagerflow/internal/config"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func init() {
	godotenv.Load()
}

func main() {
	fx.New(
		config.Module,
		bootstrap.Module,
	).Run()
}
