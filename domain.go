package main

import "time"

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
	Numbers []int `json:"numeros" validate:"required,min=1"`
}

type SumResponse struct {
	NumbersProcessed []int   `json:"numeros_processados"`
	Sum              int     `json:"soma"`
	Mean             float64 `json:"media"`
	Maximum          float64 `json:"maximo"`
	Minimum          float64 `json:"minimo"`
	Amount           float64 `json:"quantidade"`
}

type TimestampRequest struct {
}

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
