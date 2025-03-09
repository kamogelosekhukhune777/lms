package testapp

import "encoding/json"

// TestInfo represents information about the service.
type Test struct {
	Status string
}

// Encode implements the encoder interface.
func (app Test) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}
