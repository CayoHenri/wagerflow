package main

import (
	"github.com/CayoHenri/wagerflow/internal/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(bootstrap.Module).Run()
}
