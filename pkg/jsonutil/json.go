package jsonutil

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	stacktrace "github.com/palantir/stacktrace"
	"io"
)

// EncodeJSONWithoutErr - Encodes/Marshals the given object into JSON but does not return an err
func EncodeJSONWithoutErr(in interface{}) []byte {
	res, _ := EncodeJSON(in)
	return res
}

// EncodeJSON - Encodes/Marshals the given object into JSON
func EncodeJSON(in interface{}) ([]byte, error) {
	if in == nil {
		return nil, stacktrace.NewError("input for encoding is nil")
	}
	stream := jsoniter.ConfigFastest.BorrowStream(nil)
	defer jsoniter.ConfigFastest.ReturnStream(stream)
	stream.WriteVal(in)
	if stream.Error != nil {
		return nil, stacktrace.Propagate(stream.Error, "Failed to encode JSON")
	}
	return stream.Buffer(), nil
}

// DecodeJSON -
func DecodeJSON(data []byte, out interface{}) error {
	if len(data) == 0 {
		return stacktrace.NewError("'data' being decoded is nil")
	}
	if out == nil {
		return stacktrace.NewError("output parameter 'out' is nil")
	}
	iter := jsoniter.ConfigFastest.BorrowIterator(data)
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	iter.ReadVal(&out)
	if iter.Error != nil {
		return stacktrace.Propagate(iter.Error, "Failed to decode JSON Blob")
	}
	return nil
}

// EncodeJSONWithIndentation - Encodes/Marshals the given object into JSON
// DEPRACATED
func EncodeJSONWithIndentation(in interface{}) ([]byte, error) {
	// if in == nil {
	// 	return nil, stacktrace.NewError("input for encoding is nil")
	// }
	buf := new(bytes.Buffer)
	EncodeJSONToWriter(buf, in, "", "    ")
	return buf.Bytes(), nil
}

// EncodeJSONToWriter - encodes/marshals a given interface
// to an io writer. it can also indent the output
func EncodeJSONToWriter(w io.Writer, in interface{}, prefix, indent string) error {
	if w == nil {
		return stacktrace.NewError("io.Writer is nil")
	}
	enc := jsoniter.NewEncoder(w)
	enc.SetEscapeHTML(true)
	if len(prefix) != 0 && len(indent) != 0 {
		enc.SetIndent(prefix, indent)
	}
	return enc.Encode(in)
}

// DecodeJSONFromReader - Decodes/Unmarshals the given
// io.Reader pointing to a JSON, into a desired object
func DecodeJSONFromReader(r io.Reader, out interface{}) error {
	if r == nil {
		return stacktrace.NewError("'io.Reader' being decoded is nil")
	}
	if out == nil {
		return stacktrace.NewError("output parameter 'out' is nil")
	}
	dec := jsoniter.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(out)
}
