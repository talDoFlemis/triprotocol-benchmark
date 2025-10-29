package main

import (
	"encoding/json"
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
	return a.Code + ": " + a.Message
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
	IsOperation() bool
	CommandOrOperationName() string
}

type OperationResponse interface {
	OperationResponseName() string
}

type ISO8601Time struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface for ISO8601Time.
func (t ISO8601Time) MarshalJSON() ([]byte, error) {
	// Format the time as an ISO 8601 string (RFC3339 is a common ISO 8601 variant).
	// time.RFC3339 provides "YYYY-MM-DDTHH:MM:SSZ" or "YYYY-MM-DDTHH:MM:SS-ZZ:ZZ" format.
	s := t.Format(time.RFC3339)
	return json.Marshal(s)
}

// UnmarshalJSON implements the json.Unmarshaler interface for ISO8601Time.
func (t *ISO8601Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	// Parse the ISO 8601 string back into a time.Time.
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*t = ISO8601Time{parsedTime}
	return nil
}

var (
	_ OperationRequest  = (*AuthRequest)(nil)
	_ OperationResponse = (*AuthResponse)(nil)
	_ OperationRequest  = (*EchoRequest)(nil)
	_ OperationResponse = (*EchoResponse)(nil)
	_ OperationRequest  = (*SumRequest)(nil)
	_ OperationResponse = (*SumResponse)(nil)
	_ OperationRequest  = (*TimestampRequest)(nil)
	_ OperationResponse = (*TimestampResponse)(nil)
	_ OperationRequest  = (*StatusRequest)(nil)
	_ OperationResponse = (*StatusResponse)(nil)
	_ OperationRequest  = (*HistoryRequest)(nil)
	_ OperationResponse = (*HistoryResponse)(nil)
	_ OperationRequest  = (*LogoutRequest)(nil)
	_ OperationResponse = (*LogoutResponse)(nil)
)

type AuthRequest struct {
	StudentID string    `json:"aluno_id" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
}

// IsOperation implements OperationRequest.
func (a AuthRequest) IsOperation() bool {
	return false
}

// CommandOrOperationName implements OperationRequest.
func (a AuthRequest) CommandOrOperationName() string {
	return "AUTH"
}

type AuthResponse struct {
	Token      string `json:"token"`
	Name       string `json:"nome"`
	Enrollment string `json:"matricula"`
}

// OperationResponseName implements OperationResponse.
func (a AuthResponse) OperationResponseName() string {
	return "AUTH_RESPONSE"
}

type EchoRequest struct {
	Message string `json:"mensagem"`
}

// CommandOrOperationName implements OperationRequest.
func (e EchoRequest) CommandOrOperationName() string {
	return "echo"
}

// IsOperation implements OperationRequest.
func (e EchoRequest) IsOperation() bool {
	return true
}

type EchoResponse struct {
	OriginalMessage string    `json:"mensagem_original"`
	EchoMessage     string    `json:"mensagem_echo"`
	ServerTimestamp time.Time `json:"timestamp_servidor"`
	MessageSize     string    `json:"tamanho_mensagem"`
	HashMD5         string    `json:"hash_md5"`
}

// OperationResponseName implements OperationResponse.
func (e EchoResponse) OperationResponseName() string {
	return "ECHO_RESPONSE"
}

type SumRequest struct {
	Numbers []int `json:"numeros" validate:"required,min=1,max=1000"`
}

// IsOperation implements OperationRequest.
func (s SumRequest) IsOperation() bool {
	return true
}

// CommandOrOperationName implements OperationRequest.
func (s SumRequest) CommandOrOperationName() string {
	return "soma"
}

type SumResponse struct {
	NumbersProcessed []int   `json:"numeros_processados"`
	Sum              int     `json:"soma"`
	Mean             float64 `json:"media"`
	Maximum          float64 `json:"maximo"`
	Minimum          float64 `json:"minimo"`
	Amount           float64 `json:"quantidade"`
}

// OperationResponseName implements OperationResponse.
func (s SumResponse) OperationResponseName() string {
	return "SUM_RESPONSE"
}

type TimestampRequest struct{}

// IsOperation implements OperationRequest.
func (t TimestampRequest) IsOperation() bool {
	return true
}

// CommandOrOperationName implements OperationRequest.
func (t TimestampRequest) CommandOrOperationName() string {
	return "timestamp"
}

type TimestampResponse struct {
	FormatedTimestamp string `json:"timestamp_formatado"`
	UnixTimestamp     string `json:"timestamp_unix"`
	Timezone          string `json:"timezone"`
	DayOfWeek         string `json:"dia_semana"`
	AdditionalInfo    string `json:"informacao_adicional"`
}

// OperationResponseName implements OperationResponse.
func (t TimestampResponse) OperationResponseName() string {
	return "TIMESTAMP_RESPONSE"
}

type StatusRequest struct {
	Detailed bool `json:"detailed"`
}

// IsOperation implements OperationRequest.
func (s StatusRequest) IsOperation() bool {
	return true
}

// CommandOrOperationName implements OperationRequest.
func (s StatusRequest) CommandOrOperationName() string {
	return "STATUS"
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

// OperationResponseName implements OperationResponse.
func (s StatusResponse) OperationResponseName() string {
	return "STATUS_RESPONSE"
}

type HistoryRequest struct {
	Limit int `json:"limite" validate:"required,min=1,max=100"`
}

// IsOperation implements OperationRequest.
func (h HistoryRequest) IsOperation() bool {
	return true
}

// CommandOrOperationName implements OperationRequest.
func (h HistoryRequest) CommandOrOperationName() string {
	return "HISTORY"
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

// OperationResponseName implements OperationResponse.
func (h HistoryResponse) OperationResponseName() string {
	return "HISTORY_RESPONSE"
}

type LogoutRequest struct{}

// IsOperation implements OperationRequest.
func (l LogoutRequest) IsOperation() bool {
	return false
}

// CommandOrOperationName implements OperationRequest.
func (l LogoutRequest) CommandOrOperationName() string {
	return "LOGOUT"
}

type LogoutResponse struct {
	Message         string `json:"mensagem"`
	FinishedSession string `json:"sessao_encerrada"`
}

// OperationResponseName implements OperationResponse.
func (l LogoutResponse) OperationResponseName() string {
	return "LOGOUT_RESPONSE"
}
