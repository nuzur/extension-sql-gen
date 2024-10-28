package main

import (
	extensionsdk "github.com/nuzur/extension-sdk"
	"github.com/nuzur/extension-sql-gen/server"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		extensionsdk.Module,
		fx.Provide(
			server.New,
		),
	).Run()
}
