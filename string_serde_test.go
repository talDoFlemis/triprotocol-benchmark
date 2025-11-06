package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			expectedString: "OP|token=sumtoken|operacao=soma|nums=1,2,3|FIM\n",
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
					Timestamp:  NonISO8601Time{time.Date(2025, 10, 30, 18, 16, 4, 585339000, time.UTC)},
					Enrollment: "538349",
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &AuthResponse{}},
		},
		{
			name:        "SUM Response",
			inputString: "OK|numeros_originais=1.0,2.0,3.0|quantidade=3|soma=6.0|media=2.0|maximo=3.0|minimo=1.0|timestamp_calculo=2025-11-01T16:04:55.257055|timestamp=2025-11-01T16:04:55.256385|FIM",
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &SumResponse{
					OriginalNumbers:      []float64{1, 2, 3},
					Amount:               3,
					Sum:                  6,
					Mean:                 2,
					Maximum:              3,
					Minimum:              1,
					CalculationTimestamp: NonISO8601Time{time.Date(2025, 11, 1, 16, 4, 55, 257055000, time.UTC)},
					Timestamp:            NonISO8601Time{time.Date(2025, 11, 1, 16, 4, 55, 256385000, time.UTC)},
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &SumResponse{}},
		},
		{
			name:        "Echo Response",
			inputString: "OK|mensagem_original=tubias|mensagem_eco=ECO: tubias|timestamp_servidor=2025-10-30T21:12:41.305529|tamanho_mensagem=6|hash_md5=929a27e9c93c793fb599ab483f3f720d|timestamp=2025-10-30T21:12:41.304798|FIM",
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &EchoResponse{
					OriginalMessage: "tubias",
					ServerTimestamp: NonISO8601Time{time.Date(2025, time.October, 30, 21, 12, 41, 305529000, time.UTC)},
					EchoMessage:     "ECO: tubias",
					MessageSize:     6,
					HashMD5:         "929a27e9c93c793fb599ab483f3f720d",
					Timestamp:       NonISO8601Time{time.Date(2025, 10, 30, 21, 12, 41, 304798000, time.UTC)},
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
					ISOTimestamp:      NonISO8601Time{time.Date(2025, 10, 30, 21, 22, 51, 687237000, time.UTC)},
					UnixTimestamp:     UnixTimestamp{time.Date(2025, 10, 30, 21, 22, 51, 687242269, time.UTC)},
					FormatedTimestamp: "30/10/2025 21:22:51",
					Year:              2025,
					Month:             10,
					Day:               30,
					Hour:              21,
					Minute:            22,
					Second:            51,
					Microsecond:       687237,
					Timestamp:         NonISO8601Time{time.Date(2025, 10, 30, 21, 22, 51, 686268000, time.UTC)},
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
					TimeActive:          UnixTimestamp{time.Date(2025, 10, 30, 21, 43, 40, 585550785, time.UTC)},
					Version:             "1.0.0",
					Metrics: StatusResponseMetrics{
						SimulatedCPU:     34.94,
						SimulatedMemory:  63.05,
						LatencySimulated: 2.66,
					},
					Timestamp: NonISO8601Time{time.Date(2025, 10, 30, 21, 43, 40, 584881000, time.UTC)},
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
					TimeActive:          UnixTimestamp{time.Date(2025, 10, 31, 00, 04, 11, 470656871, time.UTC)},
					Version:             "1.0.0",
					Metrics: StatusResponseMetrics{
						SimulatedCPU:     73.87,
						SimulatedMemory:  68.77,
						LatencySimulated: 8.66,
					},
					SessionDetails: &map[string]StatusResponseSessionDetails{
						"538349": {
							TimestampLogin: UnixTimestamp{time.Date(2025, 10, 31, 00, 04, 11, 0, time.UTC)},
							IPClient:       "191.6.14.5",
							Name:           "SAID CAVALCANTE RODRIGUES",
							Enrollment:     "538349",
						},
					},
					DatabaseStatistics: &StatusDatabaseStatistics{
						TotalSessions:     30,
						TotalOperations:   98,
						OperationsPerType: StatusDatabaseOperationType{Authentication: 35, Echo: 14, History: 16, Sum: 13, Status: 11, Timestamp: 9},
						UniqueStudents:    2,
					},
					Timestamp: NonISO8601Time{time.Date(2025, 10, 31, 00, 04, 11, 470210000, time.UTC)},
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &StatusResponse{}},
		},
		{
			name:        "History Response",
			inputString: "OK|aluno_id=538349|limite_solicitado=2|total_encontrado=2|historico={'operacao': 'status', 'parametros': {'detalhado': True}, 'resultado': {'status': 'ATIVO', 'timestamp_consulta': '2025-10-31T16:48:31.155271', 'operacoes_processadas': 113, 'sessoes_ativas': 1, 'tempo_ativo': 1761929311.155286, 'versao': '1.0.0', 'estatisticas_banco': {'total_sessoes': 35, 'total_operacoes': 152, 'operacoes_por_tipo': {'autenticacao': 40, 'echo': 39, 'historico': 21, 'soma': 22, 'status': 16, 'timestamp': 14}, 'alunos_unicos': 2}, 'sessoes_detalhes': {'538349': {'timestamp_login': 1761929311, 'ip_cliente': '191.6.14.5', 'nome': 'SAID CAVALCANTE RODRIGUES', 'matricula': '538349'}}, 'metricas': {'cpu_simulado': 59.03, 'memoria_simulada': 33.28, 'latencia_simulada': 9.73}}, 'timestamp': '2025-10-31T16:48:31.156806', 'sucesso': True},{'operacao': 'timestamp', 'parametros': {}, 'resultado': {'timestamp_unix': 1761929310.9852066, 'timestamp_iso': '2025-10-31T16:48:30.985204', 'timestamp_formatado': '31/10/2025 16:48:30', 'ano': 2025, 'mes': 10, 'dia': 31, 'hora': 16, 'minuto': 48, 'segundo': 30, 'microsegundo': 985204}, 'timestamp': '2025-10-31T16:48:30.985333', 'sucesso': True}|timestamp_consulta=2025-10-31T01:14:19.616416|estatisticas={'total_operacoes': 1, 'operacoes_sucesso': 1, 'operacoes_erro': 0, 'taxa_sucesso': 100.0}|operacoes_mais_usadas=('status', 1)|timestamp=2025-10-31T01:14:19.615292|FIM",
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &HistoryResponse{
					StudentID:      "538349",
					RequestedLimit: 2,
					TotalFound:     2,
					History: []HistoryOperationHistoryResponse{
						{
							Operation: "status",
							Params: map[string]any{
								"detalhado": true,
							},
							Result: map[string]any{
								"status":                "ATIVO",
								"timestamp_consulta":    "2025-10-31T16:48:31.155271",
								"operacoes_processadas": 113,
								"sessoes_ativas":        1,
								"tempo_ativo":           1.761929311155286e+09,
								"versao":                "1.0.0",
								"estatisticas_banco": map[string]any{
									"total_sessoes":   float64(35),
									"total_operacoes": float64(152),
									"operacoes_por_tipo": map[string]any{
										"autenticacao": float64(40),
										"echo":         float64(39),
										"historico":    float64(21),
										"soma":         float64(22),
										"status":       float64(16),
										"timestamp":    float64(14),
									},
									"alunos_unicos": float64(2),
								},
								"sessoes_detalhes": map[string]any{"538349": map[string]any{"timestamp_login": 1.761929311e+09, "ip_cliente": "191.6.14.5", "nome": "SAID CAVALCANTE RODRIGUES", "matricula": "538349"}},
								"metricas":         map[string]any{"cpu_simulado": 59.03, "memoria_simulada": 33.28, "latencia_simulada": 9.73},
							},
							Timestamp: NonISO8601Time{time.Date(2025, 10, 31, 16, 48, 31, 156806000, time.UTC)},
							Success:   true,
						},
						{
							Operation: "timestamp",
							Params:    map[string]any{},
							Result: map[string]any{
								"timestamp_unix":      1761929310.9852066,
								"timestamp_iso":       "2025-10-31T16:48:30.985204",
								"timestamp_formatado": "31/10/2025 16:48:30",
								"ano":                 2025,
								"mes":                 10,
								"dia":                 31,
								"hora":                16,
								"minuto":              48,
								"segundo":             30,
								"microsegundo":        985204,
							},
							Timestamp: NonISO8601Time{time.Date(2025, 10, 31, 16, 48, 30, 985333000, time.UTC)},
							Success:   true,
						},
					},
					MostUsedOperations: [][]any{
						{
							"status", 1,
						},
					},
					Timestamp:        NonISO8601Time{time.Date(2025, 10, 31, 01, 14, 19, 615292000, time.UTC)},
					ConsultTimestamp: NonISO8601Time{time.Date(2025, 10, 31, 01, 14, 19, 616416000, time.UTC)},
					Stats: HistoryResponseStats{
						SuccessRate:       100,
						SuccessOperations: 1,
						TotalOperations:   1,
						ErroOperations:    0,
					},
				},
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
					Timestamp: NonISO8601Time{time.Date(2025, 10, 30, 21, 32, 25, 38812000, time.UTC)},
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
