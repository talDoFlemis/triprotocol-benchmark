package main

import (
	"encoding/json"
	"fmt"
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
	return json.Unmarshal(data, v)
}
