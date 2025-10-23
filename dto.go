package main

import "time"

type AuthRequest struct {
	Protocol   string    `query:"protocol" validate:"required, oneof=json proto string"`
	Enrollment string    `json:"enrollment" validate:"required"`
	Timestamp  time.Time `json:"timestamp" validate:"required"`
}

type ServerAuthRequest struct{}

type AuthResponse struct {
	Token string `json:"token"`
}
