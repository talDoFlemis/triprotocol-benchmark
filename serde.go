package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"google.golang.org/protobuf/proto"
)

const STRINGS_TAG = "strings"

type Serde interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

var (
	strserde   = StringSerde{}
	jsonserde  = JSONSerde{}
	protoserde = ProtobufSerde{}
)

var (
	_ Serde = (*StringSerde)(nil)
	_ Serde = (*JSONSerde)(nil)
	_ Serde = (*ProtobufSerde)(nil)
)

type (
	SerdeMarshall  = func(v any) ([]byte, error)
	SerdeUnmarshal = func(data []byte, v any) error
)

type (
	JSONSerde     struct{}
	ProtobufSerde struct{}
)

// Marshal implements Serde.
func (j JSONSerde) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal implements Serde.
func (j JSONSerde) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// Marshal implements Serde.
func (p ProtobufSerde) Marshal(v any) ([]byte, error) {
	msg, err := fromDomainToProto(v)
	if err != nil {
		slog.Error("Error converting domain to proto", slog.String("error", err.Error()))
		return nil, err
	}

	return proto.Marshal(msg)
}

// Unmarshal implements Serde.
func (p ProtobufSerde) Unmarshal(data []byte, v any) error {
	var msg proto.Message

	switch v.(type) {
	case *AuthRequest:
		msg = nil
	default:
		slog.Error("Error matching type into proto", slog.String("error", "invalid type"))
		return errors.New("invalid type")
	}

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

func fromDomainToProto(v any) (proto.Message, error) {
	var msg proto.Message

	switch v.(type) {
	default:
		return nil, errors.New("invalid type from domain to proto")
	}

	return msg, nil
}

func fromProtoToDomain(msg proto.Message) (any, error) {
	var domain any

	switch msg.(type) {
	default:
		return nil, errors.New("invalid type from proto to domain")
	}

	return domain, nil
}

func copyStruct(src, dst any) error {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	// Ensure destination is a settable pointer to a struct
	if dstVal.Kind() != reflect.Ptr || dstVal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a pointer to a struct")
	}

	// Ensure source is a struct or a pointer to a struct
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}
	if srcVal.Kind() != reflect.Struct {
		return fmt.Errorf("source must be a struct or a pointer to a struct")
	}

	// Ensure types are identical
	if srcVal.Type() != dstVal.Elem().Type() {
		return fmt.Errorf("source and destination structs must be of the same type")
	}

	// Copy fields
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.Elem().Field(i)

		// Only copy exported fields (uppercase) and settable fields
		if dstField.CanSet() {
			dstField.Set(srcField)
		}
	}
	return nil
}
