package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SOMA quebrado, diz que numeros nao eh uma lista
// Historico quebrado, faltando informacao na stream

func TestStringSerialization(t *testing.T) {
	tests := []struct {
		name           string
		inputStruct    any
		expectedString string
	}{
		{
			name: "AUTH request",
			inputStruct: PresentationLayerRequest{
				Body: AuthRequest{
					StudentID: "538349",
					Timestamp: time.Date(2025, 10, 10, 14, 30, 0, 0, time.UTC),
				},
			},
			expectedString: "AUTH|aluno_id=538349|timestamp=2025-10-10T14:30:00Z|FIM\n",
		},
		{
			name: "LOGOUT request",
			inputStruct: PresentationLayerRequest{
				Body:  LogoutRequest{},
				Token: "123",
			},
			expectedString: "LOGOUT|token=123|FIM\n",
		},
		{
			name: "ECHO request",
			inputStruct: PresentationLayerRequest{
				Body: EchoRequest{
					Message: "Hello, world!",
				},
				Token: "abcd",
			},
			expectedString: "OP|token=abcd|operacao=echo|mensagem=Hello, world!|FIM\n",
		},
		{
			name: "SUM request",
			inputStruct: PresentationLayerRequest{
				Body: SumRequest{
					Numbers: []int{1, 2, 3},
				},
				Token: "sumtoken",
			},
			expectedString: "OP|token=sumtoken|operacao=soma|numeros=1,2,3|FIM\n",
		},
		{
			name: "Timestamp request",
			inputStruct: PresentationLayerRequest{
				Body:  TimestampRequest{},
				Token: "timestamptoken",
			},
			expectedString: "OP|token=timestamptoken|operacao=timestamp|FIM\n",
		},
		{
			name: "Status request",
			inputStruct: PresentationLayerRequest{
				Body: StatusRequest{
					Detailed: true,
				},
				Token: "statustoken",
			},
			expectedString: "OP|token=statustoken|operacao=status|detalhado=true|FIM\n",
		},
		{
			name: "History request",
			inputStruct: PresentationLayerRequest{
				Body: HistoryRequest{
					Limit: 1,
				},
				Token: "historytoken",
			},
			expectedString: "OP|token=historytoken|operacao=historico|limite=1|FIM\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" - Marshal", func(t *testing.T) {
			serde := StringSerde{}
			result, err := serde.Marshal(tt.inputStruct)

			require.NoError(t, err, "Marshal should not return an error")
			assert.Equal(t, tt.expectedString, string(result), "Marshaled string should match expected output")
		})
	}
}

func TestStringDeserialization(t *testing.T) {
	tests := []struct {
		name           string
		inputString    string
		bindStruct     PresentationLayerResponse[OperationResponse]
		expectedStruct PresentationLayerResponse[OperationResponse]
	}{
		{
			name:        "AUTH Response",
			inputString: "OK|token=tokenauth|nome=SAID CAVALCANTE RODRIGUES|matricula=538349|timestamp=2025-10-30T18:16:04.585339|FIM",
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &AuthResponse{
					Token:      "tokenauth",
					Name:       "SAID CAVALCANTE RODRIGUES",
					Timestamp:  time.Date(2025, 10, 30, 18, 16, 4, 585339000, time.UTC),
					Enrollment: "538349",
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &AuthResponse{}},
		},
		{
			name:        "Echo Response",
			inputString: "OK|mensagem_original=tubias|mensagem_eco=ECO: tubias|timestamp_servidor=2025-10-30T21:12:41.305529|tamanho_mensagem=6|hash_md5=929a27e9c93c793fb599ab483f3f720d|timestamp=2025-10-30T21:12:41.304798|FIM",
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &EchoResponse{
					OriginalMessage: "tubias",
					ServerTimestamp: time.Date(2025, time.October, 30, 21, 12, 41, 305529000, time.UTC),
					EchoMessage:     "ECO: tubias",
					MessageSize:     6,
					HashMD5:         "929a27e9c93c793fb599ab483f3f720d",
					Timestamp:       time.Date(2025, 10, 30, 21, 12, 41, 304798000, time.UTC),
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &EchoResponse{}},
		},
		{
			name:        "Timestamp Response",
			inputString: "OK|timestamp_unix=1761859371.6872423|timestamp_iso=2025-10-30T21:22:51.687237|timestamp_formatado=30/10/2025 21:22:51|ano=2025|mes=10|dia=30|hora=21|minuto=22|segundo=51|microsegundo=687237|timestamp=2025-10-30T21:22:51.686268|FIM",
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &TimestampResponse{
					ISOTimestamp:      time.Date(2025, 10, 30, 21, 22, 51, 687237000, time.UTC),
					UnixTimestamp:     "1761859371.6872423",
					FormatedTimestamp: "30/10/2025 21:22:51",
					Year:              2025,
					Month:             10,
					Day:               30,
					Hour:              21,
					Minute:            22,
					Second:            51,
					Microsecond:       687237,
					Timestamp:         time.Date(2025, 10, 30, 21, 22, 51, 686268000, time.UTC),
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &TimestampResponse{}},
		},
		{
			name:        "Status Not Detailed Response",
			inputString: "OK|status=ATIVO|timestamp_consulta=2025-10-30T21:43:40.585539|operacoes_processadas=33|sessoes_ativas=1|tempo_ativo=1761860620.5855508|versao=1.0.0|metricas={'cpu_simulado': 34.94, 'memoria_simulada': 63.05, 'latencia_simulada': 2.66}|timestamp=2025-10-30T21:43:40.584881|FIM",
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &StatusResponse{
					Status:              "ATIVO",
					OperationsProcessed: 33,
					ActiveSessions:      1,
					TimeActive:          "1761860620.5855508",
					Version:             "1.0.0",
					Metrics: StatusResponseMetrics{
						SimulatedCPU:     34.94,
						SimulatedMemory:  63.05,
						LatencySimulated: 2.66,
					},
					Timestamp: time.Date(2025, 10, 30, 21, 43, 40, 584881000, time.UTC),
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &StatusResponse{}},
		},
		{
			name:        "Status Detailed Response",
			inputString: "OK|status=ATIVO|timestamp_consulta=2025-10-31T00:04:11.470644|operacoes_processadas=64|sessoes_ativas=1|tempo_ativo=1761869051.4706569|versao=1.0.0|estatisticas_banco={'total_sessoes': 30, 'total_operacoes': 98, 'operacoes_por_tipo': {'autenticacao': 35, 'echo': 14, 'historico': 16, 'soma': 13, 'status': 11, 'timestamp': 9}, 'alunos_unicos': 2}|sessoes_detalhes={'538349': {'timestamp_login': 1761869051, 'ip_cliente': '191.6.14.5', 'nome': 'SAID CAVALCANTE RODRIGUES', 'matricula': '538349'}}|metricas={'cpu_simulado': 73.87, 'memoria_simulada': 68.77, 'latencia_simulada': 8.66}|timestamp=2025-10-31T00:04:11.470210|FIM",
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &StatusResponse{
					Status:              "ATIVO",
					OperationsProcessed: 64,
					ActiveSessions:      1,
					TimeActive:          "1761869051.4706569",
					Version:             "1.0.0",
					Metrics: StatusResponseMetrics{
						SimulatedCPU:     73.87,
						SimulatedMemory:  68.77,
						LatencySimulated: 8.66,
					},
					SessionDetails: &map[string]StatusResponseSessionDetails{
						"538349": {
							TimestampLogin: "1761869051",
							IPClient:       "191.6.14.5",
							Name:           "SAID CAVALCANTE RODRIGUES",
							Enrollment:     "538349",
						},
					},
					// {'total_sessoes': 30, 'total_operacoes': 98, 'operacoes_por_tipo': {'autenticacao': 35, 'echo': 14, 'historico': 16, 'soma': 13, 'status': 11, 'timestamp': 9}, 'alunos_unicos': 2}
					DatabaseStatistics: &StatusDatabaseStatistics{
						TotalSessions:     30,
						TotalOperations:   98,
						OperationsPerType: StatusDatabaseOperationType{Authentication: 35, Echo: 14, History: 16, Sum: 13, Status: 11, Timestamp: 9},
						UniqueStudents:    2,
					},
					Timestamp: time.Date(2025, 10, 31, 00, 04, 11, 470210000, time.UTC),
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &StatusResponse{}},
		},
		{
			name:        "History Response",
			inputString: "OK|aluno_id=538349|limite_solicitado=1|total_encontrado=1|historico={'operacao': 'status', 'parametros': {'detalhado': True}, 'resultado': {'status': 'ATIVO', 'timestamp_consulta': '2025-10-31T01:14:19.431377', 'operacoes_processadas': 69, 'sessoes_ativas': 1, 'tempo_ativo': 1761873259.4313905, 'versao': '1.0.0', 'estatisticas_banco': {'total_sessoes': 31, 'total_operacoes': 104, 'operacoes_por_tipo': {'autenticacao': 36, 'echo': 15, 'historico': 17, 'soma': 14, 'status': 12, 'timestamp': 10}, 'alunos_unicos': 2}, 'sessoes_detalhes': {'538349': {'timestamp_login': 1761873259, 'ip_cliente': '191.6.14.5', 'nome': 'SAID CAVALCANTE RODRIGUES', 'matricula': '538349'}}, 'metricas': {'cpu_simulado': 22.68, 'memoria_simulada': 59.08, 'latencia_simulada': 4.28}}, 'timestamp': '2025-10-31T01:14:19.432029', 'sucesso': True}|timestamp_consulta=2025-10-31T01:14:19.616416|estatisticas={'total_operacoes': 1, 'operacoes_sucesso': 1, 'operacoes_erro': 0, 'taxa_sucesso': 100.0}|operacoes_mais_usadas=('status', 1)|timestamp=2025-10-31T01:14:19.615292|FIM",
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body:       &HistoryResponse{},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &HistoryResponse{}},
		},
		{
			name:        "Logout Response",
			inputString: "OK|msg=Logout realizado com sucesso|timestamp=2025-10-30T21:32:25.038812|FIM",
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &LogoutResponse{
					Message:   "Logout realizado com sucesso",
					Timestamp: time.Date(2025, 10, 30, 21, 32, 25, 38812000, time.UTC),
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &LogoutResponse{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" - Unmarshall", func(t *testing.T) {
			serde := StringSerde{}

			err := serde.Unmarshal([]byte(tt.inputString), &tt.bindStruct)

			require.NoError(t, err, "Unmarshall should not return an error")
			assert.Equal(t, tt.expectedStruct.StatusCode, tt.bindStruct.StatusCode, "Status codes should match")
			assert.Equal(t, tt.expectedStruct.Body, tt.bindStruct.Body, "Body should match")
			assert.Nil(t, tt.bindStruct.Err, "Error should be nil")
		})
	}
}
