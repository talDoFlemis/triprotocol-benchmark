package main

import (
	"context"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AppLayerClient[T OperationRequest, R OperationResponse] struct {
	Presentation Serde
	AppSettings  *AppSettings
	RoundTripper RoundTripper
}

func NewAppLayerClient[T OperationRequest, R OperationResponse](presentation Serde,
	roundTripper RoundTripper,
	appSettings *AppSettings,
) *AppLayerClient[T, R] {
	return &AppLayerClient[T, R]{
		Presentation: presentation,
		RoundTripper: roundTripper,
		AppSettings:  appSettings,
	}
}

func (c *AppLayerClient[T, R]) Auth(ctx context.Context, address string, req *AuthRequest) (*AuthResponse, error) {
	ctx, span := tracer.Start(ctx, "AppLayerClient.Auth", trace.WithAttributes(
		attribute.String("applayer.student_id", req.StudentID),
		attribute.String("transportlayer.address", address),
	))
	defer span.End()

	logger := slog.With(
		slog.String("applayer.student_id", req.StudentID),
		slog.String("address", address),
	)
	logger.DebugContext(ctx, "Performing auth operation")

	var authResponse AuthResponse

	err := internalDo(ctx, address, *req, &authResponse, "", c.Presentation, c.RoundTripper)
	if err != nil {
		logger.ErrorContext(ctx, "Auth failed", slog.String("error", err.Error()))
		return nil, err
	}

	logger.InfoContext(ctx, "Login successful", slog.String("applayer.token", authResponse.Token))

	return &authResponse, nil
}

func (c *AppLayerClient[T, R]) Do(ctx context.Context, address string, req T, resp R, token string) error {
	ctx, span := tracer.Start(ctx, "AppLayerClient.Do", trace.WithAttributes(
		attribute.String("applayer.token", token),
		attribute.String("transportlayer.address", address),
		attribute.String("applayer.operation_name", req.CommandOrOperationName()),
	))
	defer span.End()

	logger := slog.With(
		slog.String("applayer.operation_name", req.CommandOrOperationName()),
		slog.String("applayer.token", token),
		slog.String("address", address),
	)

	err := internalDo(ctx, address, req, resp, token, c.Presentation, c.RoundTripper)
	if err != nil {
		logger.ErrorContext(ctx, "Operation failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (c *AppLayerClient[T, R]) Logout(ctx context.Context, address string, req *LogoutRequest, token string) (*LogoutResponse, error) {
	ctx, span := tracer.Start(ctx, "AppLayerClient.Logout", trace.WithAttributes(
		attribute.String("applayer.token", token),
		attribute.String("transportlayer.address", address),
	))
	defer span.End()

	logger := slog.With(
		slog.String("applayer.operation_name", req.CommandOrOperationName()),
		slog.String("applayer.token", token),
		slog.String("address", address),
	)

	var resp LogoutResponse

	err := internalDo(ctx, address, *req, &resp, token, c.Presentation, c.RoundTripper)
	if err != nil {
		logger.ErrorContext(ctx, "Logout failed", slog.String("error", err.Error()))
		return nil, err
	}

	logger.InfoContext(ctx, "Logout successful")

	return &resp, nil
}

func internalDo[T OperationRequest, R OperationResponse](ctx context.Context, address string, req T, resp R, token string, serde Serde, roundTripper RoundTripper) error {
	ctx, span := tracer.Start(ctx, "AppLayerClient.internalDo", trace.WithAttributes(
		attribute.String("applayer.token", token),
		attribute.String("transportlayer.address", address),
		attribute.String("applayer.operation_name", req.CommandOrOperationName()),
	))
	defer span.End()

	logger := slog.With(
		slog.String("applayer.operation_name", req.CommandOrOperationName()),
		slog.String("applayer.token", token),
		slog.String("address", address),
	)
	logger.DebugContext(ctx, "Performing operation")

	presentationLayerReq := PresentationLayerRequest{
		Body: req,
	}

	if token != "" {
		presentationLayerReq.Token = token
	}

	rawRequest, err := serde.Marshal(presentationLayerReq)
	if err != nil {
		logger.ErrorContext(ctx, "Error serializing request", slog.String("error", err.Error()))
		return err
	}

	logger.DebugContext(ctx, "Sending request", slog.String("request", string(rawRequest)))

	rawResponse, err := roundTripper.RequestReply(ctx, address, rawRequest)
	if err != nil {
		logger.ErrorContext(ctx, "Error performing request", slog.String("error", err.Error()))
		return err
	}

	logger.DebugContext(ctx, "Received response", slog.String("response", string(rawResponse)))

	appLayerResp := PresentationLayerResponse[R]{
		Body: resp,
	}

	err = serde.Unmarshal(rawResponse, &appLayerResp)
	if err != nil {
		logger.ErrorContext(ctx, "Error deserializing response", slog.String("error", err.Error()))
		return err
	}

	if appLayerResp.StatusCode >= http.StatusBadRequest {
		logger.ErrorContext(ctx, "Operation returned error", slog.Int("status_code", appLayerResp.StatusCode))
		return appLayerResp.Err
	}

	logger.InfoContext(ctx, "Operation successful")
	return nil
}
