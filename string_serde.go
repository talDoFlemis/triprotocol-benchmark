package main

import (
	"fmt"
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
	el := reflect.ValueOf(v).Elem()
	t := el.Type()

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("one structs are supported, found %s", t.Kind().String())
	}

	args := strings.Split(string(data), "|")

	if len(args) < 3 {
		return fmt.Errorf("unknown size for args %d", len(args))
	}

	commandName, commandIdx := s.getCommandName(t)
	commandFound := args[0]

	if commandName != commandFound {
		return fmt.Errorf("command not compatible with struct found: %s, expected %s", commandFound, commandName)
	}

	args = args[1 : len(args)-1]

	argsMap := make(map[string]string, 1)

	for _, arg := range args {
		splittedArgs := strings.Split(arg, "=")
		if len(splittedArgs) != 2 {
			return fmt.Errorf("unknown size for splitted args found %d, expected 2, arg %s", len(splittedArgs), arg)
		}

		fieldTag := splittedArgs[0]
		fieldValue := splittedArgs[1]
		argsMap[fieldTag] = fieldValue
	}

	for i := range t.NumField() {
		field := el.Field(i)
		fieldTagValue := t.Field(i).Tag.Get(STRINGS_TAG)

		if fieldTagValue == "-" || i == commandIdx {
			continue
		}

		field.SetString(argsMap[fieldTagValue])
	}

	return nil
}
