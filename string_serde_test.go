package main

import (
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
		bindStruct     PresentationLayerResponse
		expectedStruct PresentationLayerResponse
	}{
		{
			name:        "AUTH Response",
			inputString: "OK|token=tokenauth|nome=SAID CAVALCANTE RODRIGUES|matricula=538349|timestamp=2025-10-30T18:16:04.585339|FIM",
			expectedStruct: PresentationLayerResponse{
				Body: AuthResponse{
					Token:     "tokenauth",
					Name:      "SAID CAVALCANTE RODRIGUES",
					Timestamp: time.Date(2025, 10, 30, 18, 16, 4, 585339, time.Local),
				},
			},
			bindStruct: PresentationLayerResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" - Unmarshall", func(t *testing.T) {
			serde := StringSerde{}

			err := serde.Unmarshal([]byte(tt.inputString), &tt.bindStruct)

			require.NoError(t, err, "Unmarshall should not return an error")
			assert.Equal(t, tt.expectedStruct, tt.bindStruct)
		})
	}
}
