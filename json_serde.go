package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

var jsonserde = JSONSerde{}

type JSONSerde struct{}

type jsonRequestWrapper struct {
	Kind      string `json:"tipo"`
	Operation string `json:"operacao,omitempty"`
	Token     string `json:"token,omitempty"`
	Params    any    `json:"parametros,omitempty"`
	StudentID string `json:"aluno_id,omitempty"`
}

type studenData struct {
	Name string `json:"nome"`
}

type jsonResponseWrapper struct {
	Message     string         `json:"mensagem,omitempty"`
	Token       string         `json:"token,omitempty"`
	Success     bool           `json:"sucesso,omitempty"`
	Result      any            `json:"resultado,omitempty"`
	StudentData *studenData    `json:"dados_aluno,omitempty"`
	Timestamp   NonISO8601Time `json:"timestamp"`
}

// Marshal implements Serde.
func (j JSONSerde) Marshal(v any) ([]byte, error) {
	typ := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	if typ.Kind() != reflect.Struct && typ.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("pointers or structs are supported, found %s", typ.Kind().String())
	}

	if typ.Kind() == reflect.Pointer && value.IsNil() {
		return nil, fmt.Errorf("nil pointer provided")
	}

	r, ok := v.(PresentationLayerRequest)
	if !ok {
		return nil, fmt.Errorf("only presentation layer requests are supported, found %s", typ.Kind().String())
	}

	request := jsonRequestWrapper{}

	switch r.Body.CommandOrOperationName() {
	case "LOGOUT":
		request.Kind = "logout"
		request.Token = r.Token
	case "AUTH":
		request.Kind = "autenticar"
		request.StudentID = r.Body.(AuthRequest).StudentID
	default:
		request.Kind = "operacao"
		request.Operation = r.Body.CommandOrOperationName()
		request.Params = r.Body
		request.Token = r.Token
	}

	return json.Marshal(request)
}

// Unmarshal implements Serde.
func (j JSONSerde) Unmarshal(data []byte, v any) error {
	typ := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	if typ.Kind() != reflect.Pointer {
		return fmt.Errorf("pointers are supported, found %s", typ.Kind().String())
	}

	if value.IsNil() {
		return fmt.Errorf("nil pointer provided")
	}

	value = value.Elem()

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

	// Get the actual body type to unmarshal the result into
	bodyElem := bodyField.Elem()

	responseWrapper := jsonResponseWrapper{
		Result: bodyElem.Addr().Interface(),
	}

	err := json.Unmarshal(data, &responseWrapper)
	if err != nil {
		return err
	}

	if !responseWrapper.Success {
		err := PresentationLayerErrorResponse{
			Code:    http.StatusText(http.StatusInternalServerError),
			Message: responseWrapper.Message,
			Details: make(map[string]any),
		}
		statusField := value.FieldByName("StatusCode")
		statusField.SetInt(int64(http.StatusInternalServerError))

		errField := value.FieldByName("Err")
		errField.Set(reflect.ValueOf(&err))
	}

	statusField := value.FieldByName("StatusCode")
	statusField.SetInt(int64(http.StatusOK))

	value.FieldByName("Err").Set(reflect.Zero(value.FieldByName("Err").Type()))

	timestampField := bodyElem.FieldByName("Timestamp")
	if timestampField.IsValid() && timestampField.CanSet() {
		timestampField.Set(reflect.ValueOf(responseWrapper.Timestamp))
	}

	bodyFieldType := bodyElem.Type()
	switch bodyFieldType.Name() {
	case "AuthResponse":
		bodyElem.FieldByName("Token").SetString(responseWrapper.Token)
		bodyElem.FieldByName("Name").Set(reflect.ValueOf(responseWrapper.StudentData.Name))
	case "LogoutResponse":
		bodyElem.FieldByName("Message").SetString(responseWrapper.Message)
	}

	return nil
}
