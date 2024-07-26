package http_support_library

import (
	"github.com/neiasit/http-support-library/http"
	"github.com/neiasit/http-support-library/http_gateway"
	"go.uber.org/fx"
)

var HttpServerModule = fx.Module(
	"http_server",
	fx.Provide(
		http.LoadConfig,
		http.NewHttpServer,
	),
	fx.Invoke(http.RunHttpServer),
)

var GrpcGatewayModule = fx.Module(
	"grpc_gateway",
	fx.Provide(
		http_gateway.NewGatewayServer,
	),
)

var Module = fx.Module(
	"http support",
	GrpcGatewayModule,
	HttpServerModule,
)
