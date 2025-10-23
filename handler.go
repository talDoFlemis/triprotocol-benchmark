package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	healthgo "github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"go.opentelemetry.io/otel/attribute"
)

type MainHandler struct {
	e            *echo.Echo
	httpSettings *HTTPSettings
	health       *healthgo.Health
}

func NewMainHandler(httpSettings *HTTPSettings, appSettings *AppSettings, health *healthgo.Health) *MainHandler {
	logger := slog.Default()

	e := echo.New()

	e.HideBanner = true

	// Set custom error handler
	e.HTTPErrorHandler = GlobalErrorHandler

	e.Use(slogecho.New(logger))
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(httpSettings.Timeout) * time.Second,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: httpSettings.CORS.Origins,
		AllowMethods: httpSettings.CORS.Methods,
		AllowHeaders: httpSettings.CORS.Headers,
	}))

	e.Use(otelecho.Middleware(appSettings.Name,
		otelecho.WithMetricAttributeFn(func(r *http.Request) []attribute.KeyValue {
			return []attribute.KeyValue{
				attribute.String("client.ip", r.RemoteAddr),
				attribute.String("user.agent", r.UserAgent()),
			}
		}),
		otelecho.WithEchoMetricAttributeFn(func(c echo.Context) []attribute.KeyValue {
			return []attribute.KeyValue{
				attribute.String("handler.path", c.Path()),
				attribute.String("handler.method", c.Request().Method),
			}
		}),
	))

	return &MainHandler{
		e:            e,
		httpSettings: httpSettings,
	}
}

// HealthCheck godoc
//
// @Summary Check the health of the service
// @Tags health
// @Produce json
// @Success 200 {object} healthgo.Check
// @Failure 503 {object} healthgo.Check
// @Router /healthz [get]
func (h *MainHandler) Health(c echo.Context) error {
	check := h.health.Measure(c.Request().Context())

	statusCode := http.StatusOK
	if check.Status != healthgo.StatusOK {
		statusCode = http.StatusServiceUnavailable
	}

	return c.JSON(statusCode, check)
}

func (h *MainHandler) Start() error {
	slog.Info("listening for requests", slog.String("ip", h.httpSettings.IP), slog.String("port", h.httpSettings.Port))
	address := h.httpSettings.IP + ":" + h.httpSettings.Port
	return h.e.Start(address)
}

func (h *MainHandler) Shutdown(ctx context.Context) error {
	slog.Info("shutting down http server")
	return h.e.Shutdown(ctx)
}
