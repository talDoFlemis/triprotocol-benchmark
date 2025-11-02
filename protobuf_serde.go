package main

import (
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

	return proto.Marshal(msg)
}

// Unmarshal implements Serde.
func (p ProtobufSerde) Unmarshal(data []byte, v any) error {
	resp, ok := v.(*PresentationLayerResponse[OperationResponse])
	if !ok {
		return fmt.Errorf("expected PresentationLayerResponse, found %v", reflect.TypeOf(v).Name())
	}

	bodyField := reflect.ValueOf(resp.Body)

	if bodyField.Kind() != reflect.Pointer {
		return fmt.Errorf("expected pointer to change inner value, found %v", reflect.TypeOf(resp.Body).Name())
	}

	msg := &protogenerated.Resposta{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		slog.Error("Error unmarshaling proto", slog.String("error", err.Error()))
		return err
	}

	domain, err := fromProtoToDomain(msg)
	if err != nil {
		slog.Error("Error converting", slog.String("error", err.Error()))
		return err
	}

	return copyStruct(domain, v)
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

		operationName := body.CommandOrOperationName()
		params := make(map[string]string)

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
				Operacao:   operationName,
				Parametros: params,
			},
		}
	default:
		return nil, fmt.Errorf("invalid type from domain to proto: %v", reflect.TypeOf(body).Name())
	}

	return msg, nil
}

func fromProtoToDomain(msg *protogenerated.Resposta) (*PresentationLayerResponse[OperationResponse], error) {
	resp := &PresentationLayerResponse[OperationResponse]{}
	respBody := reflect.ValueOf(resp)

	errorMsg := msg.GetErro()
	if errorMsg != nil {
		resp.StatusCode = http.StatusInternalServerError
		details := make(map[string]any)

		for key, value := range errorMsg.Detalhes {
			details[key] = value
		}

		resp.Err = &PresentationLayerErrorResponse{
			Code:    http.StatusText(http.StatusInternalServerError),
			Message: errorMsg.Mensagem,
			Details: details,
		}

		return resp, nil
	}

	okMsg := msg.GetOk()
	if okMsg == nil {
		return nil, fmt.Errorf("expected to have a pointer to RespostaOK")
	}

	timestampField := respBody.FieldByName("Timestamp")
	if timestampField.IsValid() && timestampField.CanSet() {
		timestampField.Set(reflect.ValueOf(okMsg.Timestamp))
	}

	for key, value := range okMsg.Dados {
		slog.Info("here", "key", key, "value", value)
	}

	return resp, nil
}
