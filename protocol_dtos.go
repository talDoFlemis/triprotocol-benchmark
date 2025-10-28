package main

import "time"

type AuthRequestToServerStringProtocol struct {
	Operation string `strings:"id=AUTH"`
	UserID    string `strings:"aluno_id"`
	Timestamp string `strings:"timestamp"`
}

type AuthResponseToServerStringProcotol struct {
	Token      string `strings:"token"`
	Name       string `strings:"nome"`
	Enrollment string `strings:"matricula"`
}

type LogoutRequestToServerStringProtocol struct {
}

type ErrorResponseStringProtocol struct {
	Message   string    `strings:"msg"`
	Timestamp time.Time `strings:"timestamp"`
}

type OperationRequestStringProtocol struct {
	Operation string ``
}
