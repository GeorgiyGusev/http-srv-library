package http_srv_library

import (
	"github.com/GeorgiyGusev/http-srv-library/http"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http_server",
	fx.Provide(
		http.LoadConfig,
		http.NewHttpServer,
	),
	fx.Invoke(http.RunHttpServer),
)
