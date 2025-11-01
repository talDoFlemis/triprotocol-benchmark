package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONSerialization(t *testing.T) {
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
			expectedString: "{\"tipo\":\"autenticar\",\"aluno_id\":\"538349\"}",
		},
		{
			name: "LOGOUT request",
			inputStruct: PresentationLayerRequest{
				Body:  LogoutRequest{},
				Token: "123",
			},
			expectedString: "{\"tipo\":\"logout\",\"token\":\"123\"}",
		},
		{
			name: "ECHO request",
			inputStruct: PresentationLayerRequest{
				Body: EchoRequest{
					Message: "ola mundo",
				},
				Token: "abcd",
			},
			expectedString: "{\"tipo\":\"operacao\",\"operacao\":\"echo\",\"token\":\"abcd\",\"parametros\":{\"mensagem\":\"ola mundo\"}}",
		},
		{
			name: "SUM request",
			inputStruct: PresentationLayerRequest{
				Body: SumRequest{
					Numbers: []int{1, 2, 3},
				},
				Token: "sumtoken",
			},
			expectedString: "{\"tipo\":\"operacao\",\"operacao\":\"soma\",\"token\":\"sumtoken\",\"parametros\":{\"numeros\":[1,2,3]}}",
		},
		// {\"tipo\":\"operacao\",\"operacao\":\"timestamp\",\"token\":\"538349:191.6.14.5:1762003084:f98696e78fa3409e32c2866749b13bd888077fb271ab94b9fe31c0f1eb856efe\",\"parametros\":{}}
		{
			name: "Timestamp request",
			inputStruct: PresentationLayerRequest{
				Body:  TimestampRequest{},
				Token: "timestamptoken",
			},
			expectedString: "{\"tipo\":\"operacao\",\"operacao\":\"timestamp\",\"token\":\"timestamptoken\",\"parametros\":{}}",
		},
		{
			name: "Status request",
			inputStruct: PresentationLayerRequest{
				Body: StatusRequest{
					Detailed: true,
				},
				Token: "statustoken",
			},
			expectedString: "{\"tipo\":\"operacao\",\"operacao\":\"status\",\"token\":\"statustoken\",\"parametros\":{\"detalhado\":true}}",
		},
		{
			name: "History request",
			inputStruct: PresentationLayerRequest{
				Body: HistoryRequest{
					Limit: 1,
				},
				Token: "historytoken",
			},
			expectedString: "{\"tipo\":\"operacao\",\"operacao\":\"historico\",\"token\":\"historytoken\",\"parametros\":{\"limite\":1}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" - Marshal", func(t *testing.T) {
			serde := JSONSerde{}
			result, err := serde.Marshal(tt.inputStruct)

			require.NoError(t, err, "Marshal should not return an error")
			assert.Equal(t, tt.expectedString, string(result), "Marshaled string should match expected output")
		})
	}
}

func TestJSONDeserialization(t *testing.T) {
	tests := []struct {
		name           string
		inputString    string
		bindStruct     PresentationLayerResponse[OperationResponse]
		expectedStruct PresentationLayerResponse[OperationResponse]
	}{
		{
			name:        "AUTH Response",
			inputString: `{"sucesso": true, "mensagem": "Autenticação realizada com sucesso", "timestamp": "2025-11-01T13:18:04.381480", "token": "tokenauth", "dados_aluno": {"nome": "SAID CAVALCANTE RODRIGUES"}, "sessao_id": 15}`,
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &AuthResponse{
					Name:      "SAID CAVALCANTE RODRIGUES",
					Token:     "tokenauth",
					Timestamp: NonISO8601Time{time.Date(2025, 11, 1, 13, 18, 4, 381480000, time.UTC)},
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &AuthResponse{}},
		},
		{
			name:        "Logout Response",
			inputString: `{"sucesso": true, "mensagem": "Logout realizado com sucesso", "timestamp": "2025-11-01T15:00:29.804192"}`,
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &LogoutResponse{
					Message:   "Logout realizado com sucesso",
					Timestamp: NonISO8601Time{time.Date(2025, 11, 1, 15, 0, 29, 804192000, time.UTC)},
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &LogoutResponse{}},
		},
		{
			name:        "Echo Response",
			inputString: `{"sucesso": true, "mensagem": "Operação realizada com sucesso", "timestamp": "2025-11-01T14:59:47.312256", "resultado": {"mensagem_original": "ola mundo", "mensagem_eco": "ECO: ola mundo", "timestamp_servidor": "2025-11-01T14:59:47.312723", "tamanho_mensagem": 9, "hash_md5": "3b2613ff007c695c2d560d0e9c9ccbcf"}} `,
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &EchoResponse{
					OriginalMessage: "ola mundo",
					ServerTimestamp: NonISO8601Time{time.Date(2025, 11, 1, 14, 59, 47, 312723000, time.UTC)},
					EchoMessage:     "ECO: ola mundo",
					MessageSize:     9,
					HashMD5:         "3b2613ff007c695c2d560d0e9c9ccbcf",
					Timestamp:       NonISO8601Time{time.Date(2025, 11, 1, 14, 59, 47, 312256000, time.UTC)},
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &EchoResponse{}},
		},
		{
			name: "Sum Response",
			inputString: `
{
  "sucesso": true,
  "mensagem": "Operação realizada com sucesso",
  "timestamp": "2025-11-01T16:04:20.875755",
  "resultado": {
    "numeros_originais": [1, 2, 3],
    "quantidade": 3,
    "soma": 6.0,
    "media": 2.0,
    "maximo": 3.0,
    "minimo": 1.0,
    "timestamp_calculo": "2025-11-01T16:04:20.876219"
  }
}
			`,
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &SumResponse{
					OriginalNumbers:      []float64{1, 2, 3},
					Amount:               3,
					Sum:                  6,
					Mean:                 2,
					Maximum:              3,
					Minimum:              1,
					CalculationTimestamp: NonISO8601Time{time.Date(2025, 11, 1, 16, 04, 20, 876219000, time.UTC)},
					Timestamp:            NonISO8601Time{time.Date(2025, 11, 1, 16, 04, 20, 875755000, time.UTC)},
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &SumResponse{}},
		},
		{
			name: "Timestamp Response",
			inputString: `{"sucesso": true, "mensagem": "Operação realizada com sucesso",
			"timestamp": "2025-11-01T16:04:21.048770",
			"resultado": {"timestamp_unix": 1762013061.0492296, "timestamp_iso": "2025-11-01T16:04:21.049227", "timestamp_formatado": "01/11/2025 16:04:21", "ano": 2025, "mes": 11, "dia": 1, "hora": 16, "minuto": 4, "segundo": 21, "microsegundo": 49227}
			}`,
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &TimestampResponse{
					ISOTimestamp:      NonISO8601Time{time.Date(2025, 11, 1, 16, 04, 21, 49227000, time.UTC)},
					UnixTimestamp:     UnixTimestamp{time.Date(2025, 11, 1, 16, 04, 21, 49229621, time.UTC)},
					FormatedTimestamp: "01/11/2025 16:04:21",
					Year:              2025,
					Month:             11,
					Day:               01,
					Hour:              16,
					Minute:            4,
					Second:            21,
					Microsecond:       49227,
					Timestamp:         NonISO8601Time{time.Date(2025, 11, 1, 16, 04, 21, 48770000, time.UTC)},
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &TimestampResponse{}},
		},
		{
			name: "Status Not Detailed Response",
			inputString: `
{
  "sucesso": true,
  "mensagem": "Operação realizada com sucesso",
  "timestamp": "2025-11-01T16:04:20.875755",
  "resultado": {
    "status": "ATIVO",
    "timestamp_consulta": "2025-11-01T16:04:21.221182",
    "operacoes_processadas": 76,
    "sessoes_ativas": 1,
    "tempo_ativo": 1762013061.221196,
    "versao": "1.0.0",
    "metricas": {
      "cpu_simulado": 56.12,
      "memoria_simulada": 38.47,
      "latencia_simulada": 1.06
    }
	}
}
`,
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &StatusResponse{
					Status:              "ATIVO",
					OperationsProcessed: 76,
					ActiveSessions:      1,
					TimeActive:          UnixTimestamp{time.Date(2025, 11, 1, 16, 4, 21, 221195936, time.UTC)},
					Version:             "1.0.0",
					Metrics: StatusResponseMetrics{
						SimulatedCPU:     56.12,
						SimulatedMemory:  38.47,
						LatencySimulated: 1.06,
					},
					Timestamp: NonISO8601Time{time.Date(2025, 11, 1, 16, 4, 20, 875755000, time.UTC)},
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &StatusResponse{}},
		},
		{
			name: "Status Detailed Response",
			inputString: `
{
  "sucesso": true,
  "mensagem": "Operação realizada com sucesso",
  "timestamp": "2025-11-01T16:04:21.220755",
  "resultado": {
    "status": "ATIVO",
    "timestamp_consulta": "2025-11-01T16:04:21.221182",
    "operacoes_processadas": 76,
    "sessoes_ativas": 1,
    "tempo_ativo": 1762013061.221196,
    "versao": "1.0.0",
    "estatisticas_banco": {
      "total_sessoes": 18,
      "total_operacoes": 94,
      "operacoes_por_tipo": {
        "autenticacao": 19,
        "echo": 14,
        "historico": 19,
        "soma": 13,
        "status": 15,
        "timestamp": 14
      },
      "alunos_unicos": 2
    },
    "sessoes_detalhes": {
      "538349": {
        "timestamp_login": 1762013061,
        "ip_cliente": "191.6.14.5",
        "nome": "SAID CAVALCANTE RODRIGUES",
        "matricula": "538349"
      }
    },
    "metricas": {
      "cpu_simulado": 56.12,
      "memoria_simulada": 38.47,
      "latencia_simulada": 1.06
    }
  }
}
`,
			expectedStruct: PresentationLayerResponse[OperationResponse]{
				Body: &StatusResponse{
					Status:              "ATIVO",
					OperationsProcessed: 76,
					ActiveSessions:      1,
					TimeActive:          UnixTimestamp{time.Date(2025, 11, 01, 16, 04, 21, 221195936, time.UTC)},
					Version:             "1.0.0",
					Metrics: StatusResponseMetrics{
						SimulatedCPU:     56.12,
						SimulatedMemory:  38.47,
						LatencySimulated: 1.06,
					},
					SessionDetails: &map[string]StatusResponseSessionDetails{
						"538349": {
							TimestampLogin: UnixTimestamp{time.Date(2025, 11, 01, 16, 04, 21, 0, time.UTC)},
							IPClient:       "191.6.14.5",
							Name:           "SAID CAVALCANTE RODRIGUES",
							Enrollment:     "538349",
						},
					},
					DatabaseStatistics: &StatusDatabaseStatistics{
						TotalSessions:     18,
						TotalOperations:   94,
						OperationsPerType: StatusDatabaseOperationType{Authentication: 19, Echo: 14, History: 19, Sum: 13, Status: 15, Timestamp: 14},
						UniqueStudents:    2,
					},
					Timestamp: NonISO8601Time{time.Date(2025, 11, 1, 16, 4, 21, 220755000, time.UTC)},
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &StatusResponse{}},
		},
		{
			name: "History Response",
			inputString: `
{
  "sucesso": true,
  "mensagem": "Operação realizada com sucesso",
  "timestamp": "2025-11-01T16:04:21.392632",
  "resultado": {
    "aluno_id": "538349",
    "limite_solicitado": 2,
    "total_encontrado": 2,
    "historico": [
      {
        "operacao": "status",
        "parametros": { "detalhado": true },
        "resultado": {
          "status": "ATIVO",
          "timestamp_consulta": "2025-11-01T16:04:21.221182",
          "operacoes_processadas": 76,
          "sessoes_ativas": 1,
          "tempo_ativo": 1762013061.221196,
          "versao": "1.0.0",
          "estatisticas_banco": {
            "total_sessoes": 18,
            "total_operacoes": 94,
            "operacoes_por_tipo": {
              "autenticacao": 19,
              "echo": 14,
              "historico": 19,
              "soma": 13,
              "status": 15,
              "timestamp": 14
            },
            "alunos_unicos": 2
          },
          "sessoes_detalhes": {
            "538349": {
              "timestamp_login": 1762013061,
              "ip_cliente": "191.6.14.5",
              "nome": "SAID CAVALCANTE RODRIGUES",
              "matricula": "538349"
            }
          },
          "metricas": {
            "cpu_simulado": 56.12,
            "memoria_simulada": 38.47,
            "latencia_simulada": 1.06
          }
        },
        "timestamp": "2025-11-01T16:04:21.221762",
        "sucesso": true
      },
      {
        "operacao": "timestamp",
        "parametros": {},
        "resultado": {
          "timestamp_unix": 1762013061.0492296,
          "timestamp_iso": "2025-11-01T16:04:21.049227",
          "timestamp_formatado": "01/11/2025 16:04:21",
          "ano": 2025,
          "mes": 11,
          "dia": 1,
          "hora": 16,
          "minuto": 4,
          "segundo": 21,
          "microsegundo": 49227
        },
        "timestamp": "2025-11-01T16:04:21.049359",
        "sucesso": true
      }
    ],
    "timestamp_consulta": "2025-11-01T16:04:21.393995",
    "estatisticas": {
      "total_operacoes": 2,
      "operacoes_sucesso": 2,
      "operacoes_erro": 0,
      "taxa_sucesso": 100.0
    },
    "operacoes_mais_usadas": [
      ["status", 1],
      ["timestamp", 1]
    ]
  }
}
`,
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
								"timestamp_consulta":    "2025-11-01T16:04:21.221182",
								"operacoes_processadas": float64(76),
								"sessoes_ativas":        float64(1),
								"tempo_ativo":           1.762013061221196e+09,
								"versao":                "1.0.0",
								"estatisticas_banco": map[string]any{
									"total_sessoes":   float64(18),
									"total_operacoes": float64(94),
									"operacoes_por_tipo": map[string]any{
										"autenticacao": float64(19),
										"echo":         float64(14),
										"historico":    float64(19),
										"soma":         float64(13),
										"status":       float64(15),
										"timestamp":    float64(14),
									},
									"alunos_unicos": float64(2),
								},
								"sessoes_detalhes": map[string]any{"538349": map[string]any{"timestamp_login": 1.762013061e+09, "ip_cliente": "191.6.14.5", "nome": "SAID CAVALCANTE RODRIGUES", "matricula": "538349"}},
								"metricas":         map[string]any{"cpu_simulado": 56.12, "memoria_simulada": 38.47, "latencia_simulada": 1.06},
							},
							Timestamp: NonISO8601Time{time.Date(2025, 11, 1, 16, 4, 21, 221762000, time.UTC)},
							Success:   true,
						},
						{
							Operation: "timestamp",
							Params:    map[string]any{},
							Result: map[string]any{
								"timestamp_unix":      1.7620130610492296e+09,
								"timestamp_iso":       "2025-11-01T16:04:21.049227",
								"timestamp_formatado": "01/11/2025 16:04:21",
								"ano":                 float64(2025),
								"mes":                 float64(11),
								"dia":                 float64(1),
								"hora":                float64(16),
								"minuto":              float64(4),
								"segundo":             float64(21),
								"microsegundo":        float64(49227),
							},
							Timestamp: NonISO8601Time{time.Date(2025, 11, 1, 16, 4, 21, 49359000, time.UTC)},
							Success:   true,
						},
					},
					MostUsedOperations: [][]any{
						{"status", float64(1)},
						{"timestamp", float64(1)},
					},
					Timestamp:        NonISO8601Time{time.Date(2025, 11, 1, 16, 4, 21, 392632000, time.UTC)},
					ConsultTimestamp: NonISO8601Time{time.Date(2025, 11, 1, 16, 4, 21, 393995000, time.UTC)},
					Stats: HistoryResponseStats{
						SuccessRate:       100,
						SuccessOperations: 2,
						TotalOperations:   2,
						ErroOperations:    0,
					},
				},
				StatusCode: http.StatusOK,
			},
			bindStruct: PresentationLayerResponse[OperationResponse]{Body: &HistoryResponse{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" - Unmarshall", func(t *testing.T) {
			serde := JSONSerde{}

			err := serde.Unmarshal([]byte(tt.inputString), &tt.bindStruct)

			require.NoError(t, err, "Unmarshall should not return an error")
			assert.Equal(t, tt.expectedStruct.StatusCode, tt.bindStruct.StatusCode, "Status codes should match")
			assert.Equal(t, tt.expectedStruct.Body, tt.bindStruct.Body, "Body should match")
			assert.Nil(t, tt.bindStruct.Err, "Error should be nil")
		})
	}
}
