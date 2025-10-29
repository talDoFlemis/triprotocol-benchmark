package main

type HandlerRequest[T any] struct {
	Protocol string `query:"protocolo" validate:"required, oneof=json proto string"`
	Payload  T
}

type HandlerAuthRequest struct {
	HandlerRequest[AuthRequest]
}
