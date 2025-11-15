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

type PresentationLayerResponse[T OperationResponse] struct {
	Body       T
	Err        *PresentationLayerErrorResponse
	StatusCode int
}

type OperationRequest interface {
	IsOperation() bool
	CommandOrOperationName() string
}

type OperationResponse interface {
	OperationResponseName() string
}

type NonISO8601Time struct {
	time.Time
}

var (
	_ json.Marshaler   = (*NonISO8601Time)(nil)
	_ json.Unmarshaler = (*NonISO8601Time)(nil)
)

// MarshalJSON implements the json.Marshaler interface for ISO8601Time.
func (t NonISO8601Time) MarshalJSON() ([]byte, error) {
	s := t.Format(time.RFC3339)
	return json.Marshal(s)
}

// UnmarshalJSON implements the json.Unmarshaler interface for ISO8601Time.
func (t *NonISO8601Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return t.Parse(s)
}

func (t *NonISO8601Time) Parse(s string) error {
	parsedTime, err := time.Parse("2006-01-02T15:04:05.000000", s)
	if err != nil {
		// Try without seconds, because some responses are dumb now
		parsedTime, err = time.Parse("2006-01-02T15:04", s)
		if err != nil {
			return err
		}
	}
	*t = NonISO8601Time{parsedTime}
	return nil
}

type UnixTimestamp struct {
	time.Time
}

var (
	_ json.Marshaler   = (*UnixTimestamp)(nil)
	_ json.Unmarshaler = (*UnixTimestamp)(nil)
)

// MarshalJSON implements the json.Marshaler interface for UnixTimestamp.
func (t UnixTimestamp) MarshalJSON() ([]byte, error) {
	// Convert time.Time to unix timestamp with nanoseconds precision as a float
	seconds := t.Unix()
	nanos := t.Nanosecond()
	floatTimestamp := float64(seconds) + float64(nanos)/1e9
	return json.Marshal(floatTimestamp)
}

// UnmarshalJSON implements the json.Unmarshaler interface for UnixTimestamp.
func (t *UnixTimestamp) UnmarshalJSON(data []byte) error {
	var floatValue float64
	if err := json.Unmarshal(data, &floatValue); err != nil {
		return err
	}

	// Convert unix timestamp to time.Time
	seconds := int64(floatValue)
	nanos := int64((floatValue - float64(seconds)) * 1e9)
	*t = UnixTimestamp{time.Unix(seconds, nanos).UTC()}
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
	Token      string         `json:"token"`
	Name       string         `strings:"nome"`
	Enrollment string         `strings:"matricula"`
	Timestamp  NonISO8601Time `json:"timestamp"`
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
	OriginalMessage string         `json:"mensagem_original"`
	EchoMessage     string         `json:"mensagem_eco"`
	ServerTimestamp NonISO8601Time `json:"timestamp_servidor"`
	MessageSize     int            `json:"tamanho_mensagem"`
	HashMD5         string         `json:"hash_md5"`
	Timestamp       NonISO8601Time `json:"timestamp"`
}

// OperationResponseName implements OperationResponse.
func (e EchoResponse) OperationResponseName() string {
	return "ECHO_RESPONSE"
}

type SumRequest struct {
	Numbers []int `json:"numeros" strings:"nums" validate:"required,min=1,max=1000"`
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
	OriginalNumbers      []float64      `json:"numeros_originais"`
	Sum                  float64        `json:"soma"`
	Mean                 float64        `json:"media"`
	Maximum              float64        `json:"maximo"`
	Minimum              float64        `json:"minimo"`
	Amount               float64        `json:"quantidade"`
	Timestamp            NonISO8601Time `json:"timestamp"`
	CalculationTimestamp NonISO8601Time `json:"timestamp_calculo"`
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
	FormatedTimestamp string         `json:"timestamp_formatado"`
	ISOTimestamp      NonISO8601Time `json:"timestamp_iso"`
	UnixTimestamp     UnixTimestamp  `json:"timestamp_unix"`
	Year              int            `json:"ano"`
	Month             int            `json:"mes"`
	Day               int            `json:"dia"`
	Hour              int            `json:"hora"`
	Minute            int            `json:"minuto"`
	Second            int            `json:"segundo"`
	Microsecond       int            `json:"microsegundo"`
	Timestamp         NonISO8601Time `json:"timestamp"`
}

// OperationResponseName implements OperationResponse.
func (t TimestampResponse) OperationResponseName() string {
	return "TIMESTAMP_RESPONSE"
}

type StatusRequest struct {
	Detailed bool `json:"detalhado"`
}

// IsOperation implements OperationRequest.
func (s StatusRequest) IsOperation() bool {
	return true
}

// CommandOrOperationName implements OperationRequest.
func (s StatusRequest) CommandOrOperationName() string {
	return "status"
}

type StatusResponseMetrics struct {
	SimulatedCPU     float64 `json:"cpu_simulado"`
	SimulatedMemory  float64 `json:"memoria_simulada"`
	LatencySimulated float64 `json:"latencia_simulada"`
}

type StatusResponseSessionDetails struct {
	TimestampLogin UnixTimestamp `json:"timestamp_login"`
	IPClient       string        `json:"ip_cliente"`
	Name           string        `json:"nome"`
	Enrollment     string        `json:"matricula"`
}

type StatusDatabaseOperationType struct {
	Authentication int `json:"autenticacao"`
	Echo           int `json:"echo"`
	History        int `json:"historico"`
	Sum            int `json:"soma"`
	Status         int `json:"status"`
	Timestamp      int `json:"timestamp"`
}

type StatusDatabaseStatistics struct {
	TotalSessions     int                         `json:"total_sessoes"`
	TotalOperations   int                         `json:"total_operacoes"`
	OperationsPerType StatusDatabaseOperationType `json:"operacoes_por_tipo"`
	UniqueStudents    int                         `json:"alunos_unicos"`
}

type StatusResponse struct {
	Status              string                                   `json:"status"`
	OperationsProcessed int                                      `json:"operacoes_processadas"`
	TimeActive          UnixTimestamp                            `json:"tempo_ativo"`
	Version             string                                   `json:"versao"`
	ActiveSessions      int                                      `json:"sessoes_ativas,omitempty"`
	Timestamp           NonISO8601Time                           `json:"timestamp"`
	DatabaseStatistics  *StatusDatabaseStatistics                `json:"estatisticas_banco,omitempty"`
	SessionDetails      *map[string]StatusResponseSessionDetails `json:"sessoes_detalhes,omitempty"`
	Metrics             StatusResponseMetrics                    `json:"metricas"`
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
	return "historico"
}

type HistoryResponseStats struct {
	TotalOperations   int     `json:"total_operacoes"`
	SuccessOperations int     `json:"operacoes_sucesso"`
	ErroOperations    int     `json:"operacoes_erro"`
	SuccessRate       float64 `json:"taxa_sucesso"`
}

type HistoryOperationHistoryResponse struct {
	Operation string         `json:"operacao"`
	Params    map[string]any `json:"parametros,omitempty"`
	Result    map[string]any `json:"resultado,omitempty"`
	Timestamp NonISO8601Time `json:"timestamp"`
	Success   bool           `json:"sucesso"`
}

type HistoryResponse struct {
	StudentID          string                            `json:"aluno_id"`
	RequestedLimit     int                               `json:"limite_solicitado"`
	TotalFound         int                               `json:"total_encontrado"`
	History            []HistoryOperationHistoryResponse `json:"historico"`
	ConsultTimestamp   NonISO8601Time                    `json:"timestamp_consulta"`
	Stats              HistoryResponseStats              `json:"estatisticas"`
	MostUsedOperations [][]any                           `json:"operacoes_mais_usadas,omitempty"`
	Timestamp          NonISO8601Time                    `json:"timestamp"`
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
	Message   string         `json:"mensagem" strings:"msg"`
	Timestamp NonISO8601Time `json:"timestamp"`
}

// OperationResponseName implements OperationResponse.
func (l LogoutResponse) OperationResponseName() string {
	return "LOGOUT_RESPONSE"
}
