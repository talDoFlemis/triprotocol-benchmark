package main

import (
	"encoding/binary"
	"testing"
	"time"
)

func BenchmarkMarshallProtocols(b *testing.B) {
	// Arrange
	benchmarks := []struct {
		name        string
		inputStruct PresentationLayerRequest
	}{
		{
			name: "AuthRequest",
			inputStruct: PresentationLayerRequest{
				Body: AuthRequest{
					StudentID: "538349",
					Timestamp: time.Date(2025, 10, 10, 14, 30, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "EchoRequest",
			inputStruct: PresentationLayerRequest{
				Body: EchoRequest{
					Message: "Ola mundo",
				},
				Token: "tokenecho",
			},
		},
		{
			name: "SumRequest",
			inputStruct: PresentationLayerRequest{
				Body: SumRequest{
					Numbers: []int{1, 2, 3},
				},
				Token: "tokensum",
			},
		},
		{
			name: "HistoryRequest",
			inputStruct: PresentationLayerRequest{
				Body: HistoryRequest{
					Limit: 1,
				},
				Token: "tokenhistory",
			},
		},
		{
			name: "TimestampRequest",
			inputStruct: PresentationLayerRequest{
				Body:  TimestampRequest{},
				Token: "tokentimestamp",
			},
		},
		{
			name: "StatusRequest",
			inputStruct: PresentationLayerRequest{
				Body: StatusRequest{
					Detailed: true,
				},
				Token: "tokenstatus",
			},
		},
		{
			name: "LogoutRequest",
			inputStruct: PresentationLayerRequest{
				Body:  LogoutRequest{},
				Token: "tokenlogout",
			},
		},
	}

	serdes := []struct {
		serde Serde
		name  string
	}{
		{
			name:  "StringSerde",
			serde: &StringSerde{},
		},
		{
			name:  "JSONSerde",
			serde: &JSONSerde{},
		},
		{
			name:  "ProtobufSerde",
			serde: &ProtobufSerde{},
		},
	}

	for _, bm := range benchmarks {
		for _, serde := range serdes {
			b.Run(bm.name+"/"+serde.name, func(b *testing.B) {
				for b.Loop() {
					b.ReportAllocs()

					// Act
					buf, err := serde.serde.Marshal(bm.inputStruct)
					b.ReportMetric(float64(len(buf)), "B/op")

					// Assert
					if err != nil {
						b.Fatalf("Error marshalling: %v", err)
					}
				}
			})
		}
	}
}

func BenchmarkUnmarshallProtocols(b *testing.B) {
	// Arrange
	benchmarks := []struct {
		name          string
		stringInput   string
		jsonInput     string
		protobufInput string
		bindStruct    *PresentationLayerResponse[OperationResponse]
	}{
		{
			name:        "AuthResponse",
			stringInput: "OK|token=538349:191.6.14.5:1762435811:5bde37f8efe9550fbc6bec2a06ac9f4b24d8f4cc9951baf0166b5366f82e56a0|nome=SAID CAVALCANTE RODRIGUES|matricula=538349|timestamp=2025-11-06T13:30:11.110376|FIM",
			jsonInput: `
{
  "sucesso": true,
  "mensagem": "Autenticação realizada com sucesso",
  "timestamp": "2025-11-06T13:30:34.382931",
  "token": "tokenhere",
  "dados_aluno": { "nome": "SAID CAVALCANTE RODRIGUES" },
  "sessao_id": 207
}`,
			protobufInput: "\n\xdc\x01\n\x04AUTH\x12\x18\n\x10timeout_segundos\x12\x043600\x12\x13\n\tmatricula\x12\x06538349\x12!\n\x04nome\x12\x19SAID CAVALCANTE RODRIGUES\x12f\n\x05token\x12]538349:191.6.14.5:1762435941:784a6a4951762a12fd7eac9682cb79021b9768ef80c25ef63aa159ed86692fb4\x1a\x1a2025-11-06T13:32:21.401940",
			bindStruct: &PresentationLayerResponse[OperationResponse]{
				Body: &AuthResponse{},
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name+"/"+"StringSerde", func(b *testing.B) {
			for b.Loop() {
				b.ReportAllocs()

				buf := []byte(bm.stringInput)
				b.ReportMetric(float64(len(buf)), "B/op")

				// Act
				err := StringSerde{}.Unmarshal(buf, bm.bindStruct)
				// Assert
				if err != nil {
					b.Fatalf("Error unmarshalling: %v", err)
				}
			}
		})

		b.Run(bm.name+"/"+"JSONSerde", func(b *testing.B) {
			for b.Loop() {
				b.ReportAllocs()

				buf := []byte(bm.jsonInput)
				b.ReportMetric(float64(len(buf)), "B/op")

				// Act
				err := JSONSerde{}.Unmarshal(buf, bm.bindStruct)
				// Assert
				if err != nil {
					b.Fatalf("Error unmarshalling: %v", err)
				}
			}
		})

		b.Run(bm.name+"/"+"ProtobufSerde", func(b *testing.B) {
			for b.Loop() {
				b.ReportAllocs()

				msgBytes := []byte(bm.protobufInput)
				msgLen := len(msgBytes)
				data := make([]byte, 4+msgLen)

				binary.BigEndian.PutUint32(data, uint32(msgLen))

				copy(data[4:], msgBytes)

				b.ReportMetric(float64(len(data)), "B/op")

				// Act
				err := ProtobufSerde{}.Unmarshal(data, bm.bindStruct)
				// Assert
				if err != nil {
					b.Fatalf("Error unmarshalling: %v", err)
				}
			}
		})
	}
}
