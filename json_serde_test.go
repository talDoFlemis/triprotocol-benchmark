package main

import (
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
