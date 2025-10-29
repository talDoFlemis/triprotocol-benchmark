package main

import (
	"time"
)

var _ error = (*PresentationLayerErrorResponse)(nil)

type PresentationLayerErrorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

// Error implements error.
func (a *PresentationLayerErrorResponse) Error() string {
	panic("unimplemented")
}

type PresentationLayerRequest struct {
	Token string
	Body  OperationRequest
}

type PresentationLayerResponse struct {
	Body       OperationResponse
	Err        PresentationLayerErrorResponse
	StatusCode int
}

type OperationRequest interface {
	OperationName() string
}

type OperationResponse interface {
	OperationResponseName() string
}

var (
	_ OperationRequest  = (*AuthRequest)(nil)
	_ OperationResponse = (*AuthResponse)(nil)
	_ OperationRequest  = (*LogoutRequest)(nil)
	_ OperationResponse = (*LogoutResponse)(nil)
)

type AuthRequest struct {
	StudentID string `json:"aluno_id" validate:"required"`
}

// OperationName implements OperationRequest.
func (a AuthRequest) OperationName() string {
	return "AUTH"
}

type AuthResponse struct {
	Token      string `json:"token"`
	Name       string `json:"nome"`
	Enrollment string `json:"matricula"`
}

// OperationResponseName implements OperationResponse.
func (a AuthResponse) OperationResponseName() string {
	panic("unimplemented")
}

type EchoRequest struct {
	Message string
}

type EchoResponse struct {
	OriginalMessage string    `json:"mensagem_original"`
	EchoMessage     string    `json:"mensagem_echo"`
	ServerTimestamp time.Time `json:"timestamp_servidor"`
	MessageSize     string    `json:"tamanho_mensagem"`
	HashMD5         string    `json:"hash_md5"`
}

type SumRequest struct {
	Numbers []int `json:"numeros" validate:"required,min=1,max=1000"`
}

type SumResponse struct {
	NumbersProcessed []int   `json:"numeros_processados"`
	Sum              int     `json:"soma"`
	Mean             float64 `json:"media"`
	Maximum          float64 `json:"maximo"`
	Minimum          float64 `json:"minimo"`
	Amount           float64 `json:"quantidade"`
}

type TimestampRequest struct{}

type TimestampResponse struct {
	FormatedTimestamp string `json:"timestamp_formatado"`
	UnixTimestamp     string `json:"timestamp_unix"`
	Timezone          string `json:"timezone"`
	DayOfWeek         string `json:"dia_semana"`
	AdditionalInfo    string `json:"informacao_adicional"`
}

type StatusRequest struct {
	Detailed bool `json:"detailed"`
}

type StatusResponse struct {
	Status              string    `json:"status"`
	OperationsProcessed int       `json:"operacoes_processadas"`
	TimeActive          time.Time `json:"tempo_ativo"`
	ActiveSessions      int       `json:"sessoes_ativas,omitempty"`
	DatabaseStatistics  string    `json:"estatisticas_banco,omitempty"`
	InUseMemory         string    `json:"memoria_uso,omitempty"`
	RecentConnections   string    `json:"conexoes_recentes,omitempty"`
}

type HistoryRequest struct {
	Limit int `json:"limite" validate:"required,min=1,max=100"`
}

type OperationRecord struct {
	ID         string         `json:"id"`
	Operation  string         `json:"operacao"`
	Timestamp  time.Time      `json:"timestamp"`
	Success    bool           `json:"sucesso"`
	Parameters map[string]any `json:"parametros"`
	Result     map[string]any `json:"resultado"`
}

type HistoryResponse struct {
	Operations []OperationRecord `json:"operacoes"`
	TotalFound int               `json:"total_encontrado"`
	Statistics string            `json:"estatisticas"`
}

type LogoutRequest struct{}

// OperationName implements OperationRequest.
func (l *LogoutRequest) OperationName() string {
	panic("unimplemented")
}

type LogoutResponse struct {
	Message         string `json:"mensagem"`
	FinishedSession string `json:"sessao_encerrada"`
}

// OperationResponseName implements OperationResponse.
func (l LogoutResponse) OperationResponseName() string {
	panic("unimplemented")
}
