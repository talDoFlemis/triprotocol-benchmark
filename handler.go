package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	healthgo "github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("")

type MainHandler struct {
	e            *echo.Echo
	httpSettings *HTTPSettings
	appSettings  *AppSettings
	health       *healthgo.Health
}

func NewMainHandler(httpSettings *HTTPSettings, appSettings *AppSettings, health *healthgo.Health) *MainHandler {
	logger := slog.Default()

	e := echo.New()

	e.HideBanner = true

	e.Validator = &CustomValidator{validator: GetValidator()}
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

	h := &MainHandler{
		e:            e,
		httpSettings: httpSettings,
		appSettings:  appSettings,
	}

	e.POST("/auth", h.Auth)
	e.GET("/healthz", h.Health)

	return h
}

func (h *MainHandler) Health(c echo.Context) error {
	check := h.health.Measure(c.Request().Context())

	statusCode := http.StatusOK
	if check.Status != healthgo.StatusOK {
		statusCode = http.StatusServiceUnavailable
	}

	return c.JSON(statusCode, check)
}

func (h *MainHandler) Auth(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "MainHandler.Auth")
	defer span.End()

	var req HandlerAuthRequest

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	err = c.Validate(&req)
	if err != nil {
		return err
	}

	appRequest := AuthRequest{
		StudentID: req.Payload.StudentID,
	}

	serde, err := h.getProtocolSerde(req.Protocol)
	if err != nil {
		return err
	}

	appClient := NewAppLayerClient[AuthRequest, AuthResponse](serde, DefaultTCPRoundTripper, h.appSettings)

	resp, err := appClient.Auth(ctx, "", &appRequest)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &resp)
}

func (h *MainHandler) getProtocolSerde(protocol string) (Serde, error) {
	switch protocol {
	case "json":
		return jsonserde, nil
	case "string":
		return strserde, nil
	case "proto":
		return protoserde, nil
	default:
		slog.Error("Error getting protocol marshaller", slog.String("protocol", protocol))
		return nil, errors.New("protocol not supported")
	}
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
