package main

import (
	"bytes"
	"context"
	"log/slog"
	"net"
	"time"
)

type RoundTripper interface {
	RequestReply(ctx context.Context, address string, req []byte) ([]byte, error)
}

var _ RoundTripper = (*TCPRoundTripper)(nil)

var DefaultTCPRoundTripper = &TCPRoundTripper{
	DialTimeout:  30 * time.Second,
	WriteTimeout: 10 * time.Second,
	ReadTimeout:  10 * time.Second,
}

type TCPRoundTripper struct {
	DialTimeout  time.Duration
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

func NewTCPRoundTripper(dialTimeout time.Duration, writeTimeout time.Duration, readTimeout time.Duration) *TCPRoundTripper {
	return &TCPRoundTripper{
		DialTimeout:  dialTimeout,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
}

// RequestReply implements RoundTripper.
func (t *TCPRoundTripper) RequestReply(ctx context.Context, address string, req []byte) ([]byte, error) {
	_, span := tracer.Start(ctx, "TCPRoundTripper.RequestReply")
	defer span.End()

	slog.DebugContext(ctx, "Connecting to TCP server", slog.String("address", address))
	conn, err := net.DialTimeout("tcp", address, t.DialTimeout)
	if err != nil {
		slog.Error("Error connecting to TCP server", slog.String("address", address), slog.String("error", err.Error()))
		return nil, err
	}
	defer conn.Close()

	slog.DebugContext(ctx, "Sending request to TCP server", slog.String("address", address))
	conn.SetDeadline(time.Now().Add(t.WriteTimeout))
	_, err = conn.Write(req)
	if err != nil {
		slog.Error("Error writing to TCP server", slog.String("address", address), slog.String("error", err.Error()))
		return nil, err
	}

	buf := make([]byte, 64*1024)

	conn.SetDeadline(time.Now().Add(t.ReadTimeout))
	_, err = conn.Read(buf)
	if err != nil {
		slog.Error("Error reading from TCP server", slog.String("address", address), slog.String("error", err.Error()))
		return nil, err
	}

	trimmedData := bytes.TrimRight(buf, "\x00")

	slog.DebugContext(ctx, "Received response from TCP server", slog.String("address", address))

	return trimmedData, nil
}
