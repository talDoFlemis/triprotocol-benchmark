package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
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

	h := &MainHandler{
		e:            e,
		httpSettings: httpSettings,
		appSettings:  appSettings,
	}

	e.POST("/auth", h.Auth)
	e.GET("/healthz", h.Health)

	return h
}

// Health godoc
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

// Auth godoc
//
// @Summary Authenticate the user
// @Tags alan
// @Accept       json
// @Produce json
// @Param        enumstring    query     protocol  true  "protocol to use"  Enums(string,json,proto)
// @Param request body AuthRequest true "AuthRequest"
// @Success 200 {object} AuthResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 422 {object} ValidationErrorResponse
// @Router /auth [post]
func (h *MainHandler) Auth(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "MainHandler.Auth")
	defer span.End()

	var req AuthRequest

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	conn, err := net.DialTimeout("tcp", h.appSettings.ProtocolTestServerAddress, time.Duration(h.appSettings.TCPTimeoutInSeconds)*time.Second)
	if err != nil {
		slog.Error("Error connecting to TCP server", slog.String("address", h.appSettings.ProtocolTestServerAddress), slog.String("error", err.Error()))
		return err
	}
	defer conn.Close()

	serde, err := h.getProtocolSerde(req.Protocol)
	if err != nil {
		return err
	}

	data, err := serde.Marshal(req)
	if err != nil {
		slog.Error("Error serializing request", slog.String("protocol", req.Protocol), slog.String("error", err.Error()))
		return err
	}

	rawBytes, err := h.RequestReply(ctx, conn, data)
	if err != nil {
		return err
	}

	resp := AuthResponse{}

	err = serde.Unmarshal(rawBytes, &resp)
	if err != nil {
		slog.Error("Error unmarshaling response", slog.String("protocol", req.Protocol), slog.String("error", err.Error()))
		return err
	}

	return c.JSON(http.StatusOK, &resp)
}

func (h *MainHandler) RequestReply(ctx context.Context, conn net.Conn, data []byte) ([]byte, error) {
	ctx, span := tracer.Start(ctx, "MainHandler.RequestReply")
	defer span.End()

	err := h.WriteIntoConnection(ctx, conn, data)
	if err != nil {
		return nil, err
	}

	rawResponse, err := h.ReadFromConnection(ctx, conn)
	if err != nil {
		return nil, err
	}

	return rawResponse, nil
}

func (h *MainHandler) WriteIntoConnection(ctx context.Context, conn net.Conn, data []byte) error {
	_, span := tracer.Start(ctx, "MainHandler.WriteIntoConnection")
	defer span.End()

	conn.SetDeadline(time.Now().Add(time.Duration(h.appSettings.TCPTimeoutInSeconds) * time.Second))
	_, err := conn.Write(data)
	if err != nil {
		slog.Error("Error writing to TCP server", slog.String("address", h.appSettings.ProtocolTestServerAddress), slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (h *MainHandler) ReadFromConnection(ctx context.Context, conn net.Conn) ([]byte, error) {
	_, span := tracer.Start(ctx, "MainHandler.ReadFromConnection")
	defer span.End()
	buf := make([]byte, 4096)

	conn.SetDeadline(time.Now().Add(time.Duration(h.appSettings.TCPTimeoutInSeconds) * time.Second))
	_, err := conn.Read(buf)
	if err != nil {
		slog.Error("Error reading from TCP server", slog.String("address", h.appSettings.ProtocolTestServerAddress), slog.String("error", err.Error()))
		return nil, err
	}

	return buf, nil
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
