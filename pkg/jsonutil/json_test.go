package jsonutil_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	jsonutil "github.com/da-moon/dare-cli/pkg/jsonutil"

	"testing"
)

func TestJSONUtil_EncodeJSON(t *testing.T) {
	input := map[string]interface{}{
		"validation": "process",
		"test":       "data",
	}

	actualBytes, err := jsonutil.EncodeJSON(input)
	if err != nil {
		t.Fatalf("failed to encode JSON: %v", err)
	}

	actual := strings.TrimSpace(string(actualBytes))
	expected := `{"validation":"process","test":"data"}`

	if actual != expected {
		t.Fatalf("bad: encoded JSON: expected:%s\nactual:%s\n", expected, string(actualBytes))
	}
}

func TestJSONUtil_DecodeJSON(t *testing.T) {
	input := `{"test":"data","validation":"process"}`

	var actual map[string]interface{}

	err := jsonutil.DecodeJSON([]byte(input), &actual)
	if err != nil {
		fmt.Printf("decoding err: %v\n", err)
	}

	expected := map[string]interface{}{
		"test":       "data",
		"validation": "process",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}

func TestJSONUtil_DecodeJSONFromReader(t *testing.T) {
	input := `{"test":"data","validation":"process"}`

	var actual map[string]interface{}

	err := jsonutil.DecodeJSONFromReader(bytes.NewReader([]byte(input)), &actual)
	if err != nil {
		fmt.Printf("decoding err: %v\n", err)
	}

	expected := map[string]interface{}{
		"test":       "data",
		"validation": "process",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}
func BenchmarkJSONUtil_EncodeJSON(b *testing.B) {
	input := map[string]interface{}{
		"test":       "data",
		"validation": "process",
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		_, err := jsonutil.EncodeJSON(input)
		if err != nil {
			b.Fatalf("failed to encode JSON: %v", err)
		}
	}

}

func BenchmarkStdLib_EncodeJSON(b *testing.B) {
	input := map[string]interface{}{
		"test":       "data",
		"validation": "process",
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		err := enc.Encode(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkJSONUtil_DecodeJSON(b *testing.B) {
	input := `{"test":"data","validation":"process"}`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		var actual map[string]interface{}

		err := jsonutil.DecodeJSON([]byte(input), &actual)
		if err != nil {
			b.Logf("decoding err: %v\n", err)
		}
	}
}
func BenchmarkStdLib_DecodeJSON(b *testing.B) {
	input := `{"test":"data","validation":"process"}`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var out map[string]interface{}
		err := json.Unmarshal([]byte(input), &out)
		if err != nil {
			b.Fatal(err)
		}
	}

}
func BenchmarkJSONUtil_DecodeJSONFromReader(b *testing.B) {
	input := `{"test":"data","validation":"process"}`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		var actual map[string]interface{}

		err := jsonutil.DecodeJSONFromReader(bytes.NewReader([]byte(input)), &actual)
		if err != nil {
			fmt.Printf("decoding err: %v\n", err)
		}
	}

}

func BenchmarkStdlib_DecodeJSONFromReader(b *testing.B) {
	input := `{"test":"data","validation":"process"}`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var out map[string]interface{}
		r := bytes.NewReader([]byte(input))
		dec := json.NewDecoder(r)
		dec.UseNumber()
		err := dec.Decode(&out)
		if err != nil {
			b.Fatal(err)
		}
	}
}
