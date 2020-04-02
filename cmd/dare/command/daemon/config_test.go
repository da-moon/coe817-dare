package daemon

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigEncryptBytes(t *testing.T) {
	// Test with some input
	src := []byte("abc")
	c := &Config{
		MasterKey: base64.StdEncoding.EncodeToString(src),
	}

	result, err := c.EncryptBytes()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !bytes.Equal(src, result) {
		t.Fatalf("bad: %#v", result)
	}

	// Test with no input
	c = &Config{}
	result, err = c.EncryptBytes()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(result) > 0 {
		t.Fatalf("bad: %#v", result)
	}
}
func TestDecodeConfig(t *testing.T) {
	// Without a protocol
	input := `{"master_key": "foo"}`
	config, err := DecodeConfig(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if config.MasterKey != "foo" {
		t.Fatalf("bad: %#v", config)
	}

	if config.Protocol != 0 {
		t.Fatalf("bad: %#v", config)
	}
	// With a protocol
	input = `{"master_key": "foo", "protocol": 7}`
	config, err = DecodeConfig(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if config.MasterKey != "foo" {
		t.Fatalf("bad: %#v", config)
	}

	if config.Protocol != 7 {
		t.Fatalf("bad: %#v", config)
	}

}

func TestDecodeConfig_unknownDirective(t *testing.T) {
	input := `{"unknown_directive": "titi"}`
	_, err := DecodeConfig(bytes.NewReader([]byte(input)))
	if err == nil {
		t.Fatal("should have err")
	}
}

func TestMergeConfig(t *testing.T) {
	a := &Config{
		MasterKey: "foo",
		Protocol:  7,
	}

	b := &Config{
		Protocol:        -1,
		MasterKey:       "foo",
		DevelopmentMode: true,
	}

	c := MergeConfig(a, b)

	if c.Protocol != 7 {
		t.Fatalf("bad: %#v", c)
	}

	if c.MasterKey != "foo" {
		t.Fatalf("bad: %#v", c.MasterKey)
	}

	if !c.DevelopmentMode {
		t.Fatalf("bad: %#v", c)
	}
}

func TestReadConfigPaths_badPath(t *testing.T) {
	_, err := ReadConfigPaths([]string{"/i/shouldnt/exist/ever/rainbows"})
	if err == nil {
		t.Fatal("should have err")
	}
}

func TestReadConfigPaths_file(t *testing.T) {
	tf, err := ioutil.TempFile("", "dare")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	tf.Write([]byte(`{"master_key":"bar"}`))
	tf.Close()
	defer os.Remove(tf.Name())

	config, err := ReadConfigPaths([]string{tf.Name()})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if config.MasterKey != "bar" {
		t.Fatalf("bad: %#v", config)
	}
}

func TestReadConfigPaths_dir(t *testing.T) {
	td, err := ioutil.TempDir("", "dare")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer os.RemoveAll(td)

	err = ioutil.WriteFile(filepath.Join(td, "a.json"),
		[]byte(`{"master_key": "bar"}`), 0644)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = ioutil.WriteFile(filepath.Join(td, "b.json"),
		[]byte(`{"master_key": "baz"}`), 0644)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// A non-json file, shouldn't be read
	err = ioutil.WriteFile(filepath.Join(td, "c"),
		[]byte(`{"master_key": "bad"}`), 0644)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	config, err := ReadConfigPaths([]string{td})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if config.MasterKey != "baz" {
		t.Fatalf("bad: %#v", config)
	}
}
