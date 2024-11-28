package middlewares

import (
	"net/http"
	"strings"

	authLib "github.com/GeorgiyGusev/auth-library/provider"
	"github.com/labstack/echo/v4"
	"log/slog"
)

func NewAuthMiddleware(provider authLib.AuthProvider, logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			path := c.Path()

			if !provider.IsEndpointSecure(path) {
				return next(c)
			}

			// Extract the authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logger.Error("Missing authorization token")
				return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"message": "missing authorization token"})
			}

			// Extract the token from the header
			splitedHeader := strings.Split(authHeader, " ")
			if len(splitedHeader) != 2 {
				logger.Error("Invalid authorization token")
				return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"message": "invalid authorization token"})
			}

			tokenString := splitedHeader[1]

			// Authorize the token and the route
			userDetails, err := provider.Authorize(c.Request().Context(), path, tokenString)
			if err != nil {
				logger.Error("Authorization failed", "error", err)
				return echo.NewHTTPError(http.StatusForbidden, "authorization failed: "+err.Error())
			}

			logger.Info("Authorization succeeded", "user", userDetails.Username)

			// Add user details to the context
			c.Set(authLib.UserDetailsKey, userDetails)

			// Proceed to the next handler
			if err := next(c); err != nil {
				return err
			}

			return nil
		}
	}
}
