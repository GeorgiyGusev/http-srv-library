package core

import (
	"context"
	"errors"
	"github.com/GeorgiyGusev/http-srv-library/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	authLib "github.com/neiasit/auth-library/provider"
	"go.uber.org/fx"
	"log/slog"
	"net"
	"net/http"
	"time"
)

var httpServerTag = slog.String("server", "http_server")

func NewHttpServer(logger *slog.Logger, cfg *Config) *echo.Echo {
	// Http srv setup
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError:      true,
		LogLatency:       true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       true,
		LogRoutePath:     true,
		LogRequestID:     true,
		LogReferer:       true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				"uri", v.URI,
				"status", v.Status,
				"method", v.Method,
				"host", v.Host,
				"user_agent", v.UserAgent,
				"error", v.Error,
				"latency", v.Latency,
				"request_id", v.RequestID,
				"protocol", v.Protocol,
			)
			return nil
		},
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: cfg.AllowedHeaders,
		AllowMethods: cfg.AllowedMethods,
	}))
	return e
}

func NewHttpServerWithAuth(logger *slog.Logger, cfg *Config, authProvider authLib.AuthProvider) *echo.Echo {
	// Http srv setup
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError:      true,
		LogLatency:       true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       true,
		LogRoutePath:     true,
		LogRequestID:     true,
		LogReferer:       true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				"uri", v.URI,
				"status", v.Status,
				"method", v.Method,
				"host", v.Host,
				"user_agent", v.UserAgent,
				"error", v.Error,
				"latency", v.Latency,
				"request_id", v.RequestID,
				"protocol", v.Protocol,
			)
			return nil
		},
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: cfg.AllowedHeaders,
		AllowMethods: cfg.AllowedMethods,
	}))
	e.Use(middlewares.NewAuthMiddleware(authProvider, logger))
	return e
}

func RunHttpServer(lc fx.Lifecycle, e *echo.Echo, logger *slog.Logger, cfg *Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", cfg.Address())
			if err != nil {
				logger.Error("cannot start server", "error", err.Error(), httpServerTag)
				return err
			}
			e.Listener = listener
			logger.Info("starting server", httpServerTag, "address", cfg.Address())
			go func() {
				err := e.Start("")
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("cannot start server, force exit", "error", err.Error(), httpServerTag)
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down", httpServerTag)
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			return e.Shutdown(ctx)
		},
	})
}