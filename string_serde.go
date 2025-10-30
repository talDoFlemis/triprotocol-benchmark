package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type (
	StringSerde struct{}
)

// Marshal implements Serde.
func (s StringSerde) Marshal(v any) ([]byte, error) {
	typ := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	if typ.Kind() != reflect.Struct && typ.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("pointers or structs are supported, found %s", typ.Kind().String())
	}

	if typ.Kind() == reflect.Pointer && value.IsNil() {
		return nil, fmt.Errorf("nil pointer provided")
	}

	if typ.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	t := value.Type()

	r, ok := v.(PresentationLayerRequest)
	if !ok {
		return nil, fmt.Errorf("only presentation layer requests are supported, found %s", t.Kind().String())
	}

	args := []string{}
	body := r.Body

	prefix := body.CommandOrOperationName()

	if body.IsOperation() {
		prefix = "OP"
	}

	args = append(args, prefix)

	if r.Token != "" {
		args = append(args, fmt.Sprintf("token=%s", r.Token))
	}

	if body.IsOperation() {
		arg := fmt.Sprintf("operacao=%s", body.CommandOrOperationName())
		args = append(args, arg)
	}

	value = reflect.ValueOf(r.Body)
	t = value.Type()

	for i := range value.NumField() {
		field := value.Field(i)
		fieldTagValue := t.Field(i).Tag.Get("json")

		fieldValue := s.getStrFieldRepresentation(field)

		arg := fmt.Sprintf("%s=%v", fieldTagValue, fieldValue)
		args = append(args, arg)
	}

	// Add terminator
	args = append(args, "FIM")
	result := strings.Join(args, "|")

	result += "\n"

	return []byte(result), nil
}

func (s StringSerde) getStrFieldRepresentation(field reflect.Value) string {
	var fieldValue string

	inter := field.Interface()

	switch value := inter.(type) {
	case []int:
		numbers := []string{}
		for _, number := range value {
			numbers = append(numbers, strconv.Itoa(number))
		}
		fieldValue = strings.Join(numbers, ",")
	case time.Time:
		fieldValue = value.Format(time.RFC3339)
	case bool:
		fieldValue = strconv.FormatBool(value)
	case int:
		fieldValue = strconv.Itoa(value)
	default:
		fieldValue = field.String()
	}

	return fieldValue
}

func (s StringSerde) getCommandName(t reflect.Type) (string, int) {
	// Command name defaults to struct name
	commandName := t.Name()
	commandIdx := -1

	for i := range t.NumField() {
		field := t.Field(i)

		if field.Tag.Get(STRINGS_TAG) == "id" {
			commandName = strings.ToUpper(field.Name)
			commandIdx = i
			break
		}
	}

	return commandName, commandIdx
}

// Unmarshal implements Serde.
func (s StringSerde) Unmarshal(data []byte, v any) error {
	typ := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	if typ.Kind() != reflect.Pointer {
		return fmt.Errorf("pointers are supported, found %s", typ.Kind().String())
	}

	if value.IsNil() {
		return fmt.Errorf("nil pointer provided")
	}

	value = value.Elem()
	typ = reflect.TypeOf(value)

	_, ok := v.(*PresentationLayerResponse)
	if !ok {
		return fmt.Errorf("expected presentatino layer response, found %v", value.Type().Kind())
	}

	dataArgs := strings.Split(string(data), "|")

	if len(dataArgs) < 3 {
		return fmt.Errorf("invalid response from server, expected at least 3 parameters, found %d, data %s", len(dataArgs), string(data))
	}

	// Ignore FIM token
	dataArgs = dataArgs[:len(dataArgs)-1]

	status := dataArgs[0]

	// Remove status from slice
	dataArgs = dataArgs[1:]

	var statusCode int

	switch status {
	case "OK":
		statusCode = http.StatusOK
	case "INVALIDO":
		statusCode = http.StatusUnprocessableEntity
	case "ERROR":
		statusCode = http.StatusInternalServerError
	default:
		return fmt.Errorf("unexpected status code %s", status)
	}

	properties := make(map[string]string)

	for _, arg := range dataArgs {
		split := strings.Split(arg, "=")
		if len(split) != 2 {
			return fmt.Errorf("expected 2 args after spliting argument, found %d", len(split))
		}

		property := split[0]
		value := split[1]
		properties[property] = value
	}

	statusField := value.FieldByName("StatusCode")
	statusField.SetInt(int64(statusCode))

	if statusCode != http.StatusOK {
		err := PresentationLayerErrorResponse{
			Code:    http.StatusText(statusCode),
			Message: properties["msg"],
			Details: make(map[string]any),
		}

		errField := value.FieldByName("Err")
		errField.Set(reflect.ValueOf(err))

		return nil
	}

	return nil
}
