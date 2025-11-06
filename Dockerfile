FROM golang:1.25.3-alpine AS builder
WORKDIR /app

RUN apk add --no-cache \
  gcc \
  musl-dev \
  ca-certificates

COPY go.mod go.sum /app/

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=cache,target="/root/.cache/go-build" \
  CGO_ENABLED=0 GOOS=linux go build -o tui  -ldflags '-s -w -extldflags "-static"' .

FROM ubuntu:oracular AS user
RUN useradd -u 10001 scratchuser

USER scratchuser
RUN touch /tmp/debug.log

FROM scratch
WORKDIR /app


COPY --from=builder /app/tui ./
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=user /etc/passwd /etc/passwd
COPY --from=user /tmp/debug.log ./


USER scratchuser
STOPSIGNAL SIGINT

CMD ["/app/tui"]
