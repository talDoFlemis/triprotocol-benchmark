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
