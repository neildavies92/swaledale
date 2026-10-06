package financeconfig

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRejectInvalidJSON(t *testing.T) {
	valid := encoded(t, fixture())
	for name, input := range map[string]string{
		"empty": "", "truncated": `{"version":1`, "root null": "null", "root array": "[]",
		"second document": valid + valid, "trailing garbage": valid + "oops",
		"unknown field":               strings.Replace(valid, `"version":1`, `"version":1,"password":"synthetic"`, 1),
		"unknown nested field":        strings.Replace(valid, `"name":"Alex"`, `"name":"Alex","databaseId":123`, 1),
		"external identity forbidden": strings.Replace(valid, `"connector":"manual"`, `"connector":"manual","externalAccountId":"synthetic"`, 1),
		"duplicate field":             strings.Replace(valid, `"version":1`, `"version":2,"version":1`, 1),
		"case duplicate":              strings.Replace(valid, `"version":1`, `"version":2,"Version":1`, 1),
		"escaped duplicate":           strings.Replace(valid, `"version":1`, `"version":2,"\u0076ersion":1`, 1),
		"nested duplicate":            strings.Replace(valid, `"name":"Alex"`, `"name":"Other","name":"Alex"`, 1),
		"null scalar":                 strings.Replace(valid, `"shareBasisPoints":7000`, `"shareBasisPoints":null`, 1),
		"null member":                 strings.Replace(valid, `{"key":"member_a","name":"Alex"}`, `null`, 1),
		"fractional share":            strings.Replace(valid, `"shareBasisPoints":7000`, `"shareBasisPoints":7000.1`, 1),
		"string share":                strings.Replace(valid, `"shareBasisPoints":7000`, `"shareBasisPoints":"7000"`, 1),
		"wrong type":                  strings.Replace(valid, `"version":1`, `"version":"1"`, 1),
		"date":                        strings.Replace(valid, `2026-01-01T00:00:00Z`, `2026-01-01`, 1),
		"invalid UTF8":                string([]byte{0xff}),
		"too large":                   strings.Repeat(" ", MaxBytes+1),
		"deeply nested":               strings.Repeat("[", 34) + "1" + strings.Repeat("]", 34),
	} {
		t.Run(name, func(t *testing.T) {
			c, err := Load(strings.NewReader(input))
			if err == nil {
				t.Fatal("expected error")
			}
			if c.Version != 0 {
				t.Fatal("returned invalid partial config")
			}
		})
	}
}
func TestExampleLoads(t *testing.T) {
	c, err := LoadFile("../../config/finance.example.json")
	requireValid(t, err)
	if len(c.Accounts) != 11 || len(c.Members) != 3 || len(c.Bindings) != 11 {
		t.Fatal("example no longer covers household inventory")
	}
	for _, a := range c.Accounts {
		requireValid(t, a.definition().ValidateDefinition())
	}
}
func TestLoadFileAndReaderErrors(t *testing.T) {
	_, err := LoadFile(filepath.Join(t.TempDir(), "missing.json"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("got %v", err)
	}
	sentinel := errors.New("reader failed")
	_, err = Load(errorReader{sentinel})
	if !errors.Is(err, sentinel) {
		t.Fatalf("got %v", err)
	}
}

type errorReader struct{ err error }

func (r errorReader) Read([]byte) (int, error) { return 0, r.err }

var _ io.Reader = errorReader{}
