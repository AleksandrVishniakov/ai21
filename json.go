package ai21

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func encodeJSONToReader[T any](obj T) (io.Reader, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(%+v): %w", obj, err)
	}

	return bytes.NewBuffer(buf), nil
}

func decodeJSON[T any](reader io.Reader) (T, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return *new(T), fmt.Errorf("io.ReadAll(reader): %w", err)
	}

	var obj T
	err = json.Unmarshal(buf, &obj)
	if err != nil {
		return *new(T), fmt.Errorf("json.Unmarshal(): %w", err)
	}

	return obj, nil
}
