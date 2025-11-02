package main

import (
	"fmt"
	"reflect"
)

type Serde interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

var strserde = StringSerde{}

var (
	_ Serde = (*StringSerde)(nil)
	_ Serde = (*JSONSerde)(nil)
	_ Serde = (*ProtobufSerde)(nil)
)

type (
	SerdeMarshall  = func(v any) ([]byte, error)
	SerdeUnmarshal = func(data []byte, v any) error
)

func copyStruct(src, dst any) error {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	// Ensure destination is a settable pointer to a struct
	if dstVal.Kind() != reflect.Pointer || dstVal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a pointer to a struct")
	}

	// Ensure source is a struct or a pointer to a struct
	if srcVal.Kind() == reflect.Pointer {
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
