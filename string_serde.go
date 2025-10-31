package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
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

	_, ok := v.(*PresentationLayerResponse[OperationResponse])
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
		errField.Set(reflect.ValueOf(&err))

		return nil
	}

	value.FieldByName("Err").Set(reflect.Zero(value.FieldByName("Err").Type()))

	bodyField := value.FieldByName("Body").Elem()
	if bodyField.Kind() != reflect.Pointer || bodyField.IsNil() {
		return fmt.Errorf("body field is not a pointer or is nil: %v", bodyField.Kind())
	}

	bodyField = bodyField.Elem()

	return bindStructFields(bodyField, properties)
}

func bindStructFields(v reflect.Value, properties map[string]string) error {
	typ := v.Type()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := typ.Field(i)
		fieldTagValue := getFieldTagValue(fieldType)
		tagValues := strings.Split(fieldTagValue, ",")

		if len(tagValues) == 0 {
			return fmt.Errorf("field %s does not have a json tag", fieldType.Name)
		}

		omitEmpty := strings.Contains(fieldTagValue, "omitempty")
		propertyName := tagValues[0]

		fieldValueStr, ok := properties[propertyName]
		if !ok && !omitEmpty {
			return fmt.Errorf("property %s not found", propertyName)
		}

		if !ok && omitEmpty {
			continue
		}

		err := setFieldValueFromString(field, fieldValueStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func getFieldTagValue(field reflect.StructField) string {
	fieldTagValue := field.Tag.Get(STRINGS_TAG)

	if fieldTagValue == "" {
		fieldTagValue = field.Tag.Get("json")
	}

	if fieldTagValue == "" {
		fieldTagValue = field.Name
	}

	return fieldTagValue
}

func setFieldValueFromString(field reflect.Value, valueStr string) error {
	switch field.Kind() {
	case reflect.Pointer:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}

		field = field.Elem()
		return setFieldValueFromString(field, valueStr)
	case reflect.String:
		field.SetString(valueStr)
	case reflect.Int:
		fieldValue, err := strconv.Atoi(valueStr)
		if err != nil {
			return err
		}
		field.SetInt(int64(fieldValue))
	case reflect.Float64:
		fieldValue, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return err
		}
		field.SetFloat(fieldValue)
	case reflect.Bool:
		fieldValue, err := strconv.ParseBool(valueStr)
		if err != nil {
			return err
		}
		field.SetBool(fieldValue)
	case reflect.Map:
		mapType := field.Type()

		keyType := mapType.Key()
		if keyType.Kind() != reflect.String {
			return fmt.Errorf("existing map key is not string, got %s", keyType.Kind())
		}

		valueType := mapType.Elem()

		mapProperties, err := generateMapStrinStringFromValueStr(valueStr)
		if err != nil {
			return err
		}

		newMapType := reflect.MapOf(reflect.TypeOf(""), valueType)
		newMap := reflect.MakeMap(newMapType)

		for key, strValue := range mapProperties {
			newValue := reflect.New(valueType)
			err = setFieldValueFromString(newValue, strValue)
			if err != nil {
				return err
			}

			newKey := reflect.ValueOf(key)
			newMap.SetMapIndex(newKey, newValue.Elem())
		}

		field.Set(newMap)

	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			fieldValue, err := time.Parse("2006-01-02T15:04:05.000000", valueStr)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(fieldValue))
			return nil
		}

		// Fallback to python dict parsing, first level only
		// Second level we use recursion, as second map will be a string
		mapProperties, err := generateMapStrinStringFromValueStr(valueStr)
		if err != nil {
			return err
		}

		err = bindStructFields(field, mapProperties)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported field type %v", field.Type())
	}

	return nil
}

func generateMapStrinStringFromValueStr(stringValue string) (map[string]string, error) {
	properties := make(map[string]any)
	jsonStr := convertPythonDictTOJSONDict(stringValue)

	d := json.NewDecoder(strings.NewReader(jsonStr))
	d.UseNumber()
	err := d.Decode(&properties)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling struct field from string: %w", err)
	}

	mapProperties := make(map[string]string)
	for key, value := range properties {
		switch v := value.(type) {
		case map[string]any:
			buf, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("error marshaling map field to string: %w", err)
			}
			mapProperties[key] = string(buf)
		default:
			mapProperties[key] = fmt.Sprintf("%v", v)
		}
	}

	return mapProperties, nil
}

func convertPythonDictTOJSONDict(dictStr string) string {
	// Replace single quotes with double quotes
	str := regexp.MustCompile(`'`).ReplaceAllString(dictStr, `"`)
	// Replace Python boolean literals with JSON boolean literals
	str = regexp.MustCompile(`False`).ReplaceAllString(str, `false`)
	str = regexp.MustCompile(`True`).ReplaceAllString(str, `true`)

	return str
}
