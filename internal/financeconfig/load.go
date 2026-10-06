package financeconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// MaxBytes bounds human-managed configuration, not connector payloads.
const MaxBytes = 1 << 20

func LoadFile(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("open finance configuration: %w", err)
	}
	defer f.Close()
	return Load(f)
}

// Load accepts exactly one JSON document. Unknown fields, duplicate object
// fields (including case variants), nulls, and invalid references are errors.
// No environment interpolation, includes, network access or secret loading occurs.
func Load(r io.Reader) (Config, error) {
	data, err := io.ReadAll(io.LimitReader(r, MaxBytes+1))
	if err != nil {
		return Config{}, fmt.Errorf("read finance configuration: %w", err)
	}
	if len(data) > MaxBytes {
		return Config{}, fmt.Errorf("finance configuration exceeds %d bytes", MaxBytes)
	}
	if !utf8.Valid(data) {
		return Config{}, fmt.Errorf("finance configuration must be UTF-8")
	}
	scan := json.NewDecoder(bytes.NewReader(data))
	scan.UseNumber()
	if err := checkValue(scan, 0); err != nil {
		return Config{}, fmt.Errorf("finance JSON: %w", err)
	}
	if _, err := scan.Token(); err != io.EOF {
		return Config{}, fmt.Errorf("finance JSON: expected a single document")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var c Config
	if err := dec.Decode(&c); err != nil {
		return Config{}, fmt.Errorf("decode finance configuration: %w", err)
	}
	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}

// encoding/json otherwise silently accepts last-wins duplicate fields and null
// scalar values. This small token pass prevents ambiguous human-managed input.
func checkValue(dec *json.Decoder, depth int) error {
	if depth > 32 {
		return fmt.Errorf("nesting exceeds 32 levels")
	}
	token, err := dec.Token()
	if err != nil {
		return err
	}
	if token == nil {
		return fmt.Errorf("null is not allowed; omit optional fields")
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delimiter {
	case '{':
		fields := map[string]bool{}
		for dec.More() {
			field, err := dec.Token()
			if err != nil {
				return err
			}
			name, ok := field.(string)
			if !ok {
				return fmt.Errorf("expected object field")
			}
			normalized := strings.ToLower(name)
			if fields[normalized] {
				return fmt.Errorf("duplicate object field at byte %d", dec.InputOffset())
			}
			fields[normalized] = true
			if err := checkValue(dec, depth+1); err != nil {
				return err
			}
		}
	case '[':
		for dec.More() {
			if err := checkValue(dec, depth+1); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unexpected delimiter")
	}
	_, err = dec.Token()
	return err
}
