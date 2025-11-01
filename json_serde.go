package main

import "encoding/json"

var jsonserde = JSONSerde{}

type JSONSerde struct{}

// Marshal implements Serde.
func (j JSONSerde) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal implements Serde.
func (j JSONSerde) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
