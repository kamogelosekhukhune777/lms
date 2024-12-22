package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	return r.PathValue(key)
}

type validator interface {
	Validate() error
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
// If the provided value is a struct then it is checked for validation tags.
// If the value implements a validate function, it is executed.
func Decode(r *http.Request, val any) error {
	if err := unmarshalWithStrictCheck(r.Body, val); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	if v, ok := val.(validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

//=================================================================================================

func unmarshalWithStrictCheck(body io.Reader, val interface{}) error {
	// Decode into a map to check for unknown fields.
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(body).Decode(&raw); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	// Decode into the target struct.
	decoder := json.NewDecoder(bytes.NewReader(encodeMap(raw)))
	decoder.DisallowUnknownFields() // Ensure no unknown fields exist.
	if err := decoder.Decode(val); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	return nil
}

func encodeMap(raw map[string]json.RawMessage) []byte {
	data, _ := json.Marshal(raw) // This should not fail since `raw` came from JSON.
	return data
}
