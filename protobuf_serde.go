package main

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/taldoflemis/triprotocol-benchmark/protogenerated"
	"google.golang.org/protobuf/proto"
)

type (
	ProtobufSerde struct{}
)

var protoserde = ProtobufSerde{}

// Marshal implements Serde.
func (p ProtobufSerde) Marshal(v any) ([]byte, error) {
	req, ok := v.(PresentationLayerRequest)
	if !ok {
		return nil, fmt.Errorf("invalid type for protobuf marshal, expected PresentationLayerRequest, found %v", reflect.TypeOf(v).Name())
	}

	msg, err := fromDomainToProto(req)
	if err != nil {
		slog.Error("Error converting domain to proto", slog.String("error", err.Error()))
		return nil, err
	}

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	data := make([]byte, 4)
	// Calculate header size
	binary.BigEndian.PutUint32(data, uint32(len(msgBytes)))

	data = append(data, msgBytes...)

	return data, nil
}

// Unmarshal implements Serde.
func (p ProtobufSerde) Unmarshal(data []byte, v any) error {
	typ := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	if typ.Kind() != reflect.Pointer {
		return fmt.Errorf("pointers are supported, found %s", typ.Kind().String())
	}

	if value.IsNil() {
		return fmt.Errorf("nil pointer provided")
	}

	value = value.Elem()

	if len(data) < 4 {
		return fmt.Errorf("data too small, expected at least 4 bytes for header, got %d", len(data))
	}

	headerSize := binary.BigEndian.Uint32(data[:4])

	// Check for overflow and ensure we have enough data
	if headerSize > uint32(len(data)-4) {
		return fmt.Errorf("data size is smaller than header size, probably corrupted data. Expected at least %d bytes, got %d bytes", headerSize+4, len(data)-4)
	}

	data = data[4 : 4+headerSize]

	msg := &protogenerated.Resposta{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		slog.Error("Error unmarshaling proto", slog.String("error", err.Error()))
		return err
	}

	err = fromProtoToDomain(msg, value)
	if err != nil {
		slog.Error("Error converting", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func fromDomainToProto(v PresentationLayerRequest) (proto.Message, error) {
	msg := &protogenerated.Requisicao{}

	bodyValue := reflect.ValueOf(v.Body)
	bodyType := reflect.TypeOf(v)
	if bodyValue.Kind() == reflect.Pointer {
		bodyValue = bodyValue.Elem()
	}

	switch body := bodyValue.Interface().(type) {
	case AuthRequest:
		msg.Tipo = &protogenerated.Requisicao_Auth{
			Auth: &protogenerated.ComandoAuth{
				AlunoId: body.StudentID,
			},
		}
	case LogoutRequest:
		msg.Tipo = &protogenerated.Requisicao_Logout{
			Logout: &protogenerated.ComandoLogout{
				Token: v.Token,
			},
		}
	case OperationRequest:
		if !body.IsOperation() {
			return nil, fmt.Errorf("expected a operation, found %v", reflect.TypeOf(body).Name())
		}

		params := make(map[string]string)

		bodyValue = reflect.ValueOf(body)
		bodyType = reflect.TypeOf(body)

		for i := range bodyValue.NumField() {
			field := bodyValue.Field(i)
			fieldType := bodyType.Field(i)

			fieldTagValue := getFieldTagValue(fieldType)
			tagValues := strings.Split(fieldTagValue, ",")

			if len(tagValues) == 0 {
				return nil, fmt.Errorf("field %s does not have a json tag", fieldType.Name)
			}

			fieldName := tagValues[0]
			fieldValue := getStrFieldRepresentation(field)

			params[fieldName] = fieldValue
		}

		msg.Tipo = &protogenerated.Requisicao_Operacao{
			Operacao: &protogenerated.ComandoOperacao{
				Token:      v.Token,
				Operacao:   body.CommandOrOperationName(),
				Parametros: params,
			},
		}
	default:
		return nil, fmt.Errorf("invalid type from domain to proto: %v", reflect.TypeOf(body).Name())
	}

	return msg, nil
}

func fromProtoToDomain(msg *protogenerated.Resposta, value reflect.Value) error {
	errorMsg := msg.GetErro()
	if errorMsg != nil {
		details := make(map[string]any)

		for key, value := range errorMsg.Detalhes {
			details[key] = value
		}

		value.FieldByName("StatusCode").SetInt(int64(http.StatusInternalServerError))

		err := PresentationLayerErrorResponse{
			Code:    http.StatusText(http.StatusInternalServerError),
			Message: errorMsg.Mensagem,
			Details: details,
		}
		value.FieldByName("Err").Set(reflect.ValueOf(&err))

		return nil
	}

	okMsg := msg.GetOk()
	if okMsg == nil {
		return fmt.Errorf("expected to have a pointer to RespostaOK")
	}

	value.FieldByName("Err").Set(reflect.Zero(value.FieldByName("Err").Type()))
	bodyField := value.FieldByName("Body")

	if bodyField.Kind() != reflect.Pointer && bodyField.Kind() != reflect.Interface {
		return fmt.Errorf("body field is not a pointer or interface: %v", bodyField.Kind())
	}

	if bodyField.Kind() == reflect.Interface {
		bodyField = bodyField.Elem()
	}

	if bodyField.IsNil() {
		bodyField.Set(reflect.New(bodyField.Type().Elem()))
	}

	bodyField = bodyField.Elem()

	okMsg.Dados["timestamp"] = okMsg.Timestamp

	return bindStructFields(bodyField, okMsg.Dados)
}
