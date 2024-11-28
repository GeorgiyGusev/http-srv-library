package http_srv_library

import (
	"github.com/GeorgiyGusev/http-srv-library/core"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http_server",
	fx.Provide(
		core.LoadConfig,
		core.NewHttpServer,
	),
	fx.Invoke(core.RunHttpServer),
)
var ModuleWithAuth = fx.Module(
	"http_server",
	fx.Provide(
		core.LoadConfig,
		core.NewHttpServerWithAuth,
	),
	fx.Invoke(core.RunHttpServer),
)
