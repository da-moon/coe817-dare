package config_test

import (
	"fmt"
	"testing"

	config "github.com/da-moon/coe817-dare/internal/config"
	header "github.com/da-moon/coe817-dare/internal/header"
	segment "github.com/da-moon/coe817-dare/internal/segment"
	stacktrace "github.com/palantir/stacktrace"
	assert "github.com/stretchr/testify/assert"
)

type TestCase struct {
	input       uint64
	result      uint64
	description string
	success     bool
}

func TestEncryptedSize(t *testing.T) {
	var tests = []TestCase{
		{
			input:       config.MaxDecryptedSize + 1,
			result:      0,
			description: fmt.Sprintf("MoreThanMaxSize"),
			success:     false,
		},
		{
			input:       config.MaxDecryptedSize,
			result:      config.MaxEncryptedSize,
			description: fmt.Sprintf("MaxSupportedSize"),
			success:     true,
		},
		{
			input:       config.MaxPayloadSize,
			result:      config.MaxPayloadSize + config.MetadataSize,
			description: fmt.Sprintf("RegularOne"),
			success:     true,
		},
		{
			input:       2 * config.MaxPayloadSize,
			result:      2 * (config.MaxPayloadSize + config.MetadataSize),
			description: fmt.Sprintf("RegularTwo"),
			success:     true,
		},
		{
			input:       3 * config.MaxPayloadSize,
			result:      3 * (config.MaxPayloadSize + config.MetadataSize),
			description: fmt.Sprintf("RegularThree"),
			success:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			input, err := getEncryptedSize(test.input)
			if test.success {
				assert.NoError(t, err)
				assert.Equal(t, input, test.result, "expected=%v got=%v", input, test.result)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDecryptedSize(t *testing.T) {
	var tests = []TestCase{

		{
			input:       config.MaxBufferSize + 1,
			result:      0,
			description: fmt.Sprintf("MoreThanMaxBufferSize"),
			success:     false,
		},
		{
			input:       config.MaxEncryptedSize,
			result:      config.MaxDecryptedSize,
			description: fmt.Sprintf("MaxSupportedSize"),
			success:     true,
		},
		{
			input:       config.MaxPayloadSize + config.MetadataSize,
			result:      config.MaxPayloadSize,
			description: fmt.Sprintf("RegularOne"),
			success:     true,
		},
		{
			input:       2 * (config.MaxPayloadSize + config.MetadataSize),
			result:      2 * config.MaxPayloadSize,
			description: fmt.Sprintf("RegularTwo"),
			success:     true,
		},
		{
			input:       3 * (config.MaxPayloadSize + config.MetadataSize),
			result:      3 * config.MaxPayloadSize,
			description: fmt.Sprintf("RegularThree"),
			success:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			input, err := getDecryptedSize(test.input)
			if test.success {
				assert.NoError(t, err)
				assert.Equal(t, input, test.result, "expected=%v got=%v", input, test.result)
			} else {
				assert.Error(t, err)
			}
		})
	}

}

// getEncryptedSize ...
func getEncryptedSize(size uint64) (uint64, error) {
	if size > config.MaxDecryptedSize {
		return 0, stacktrace.NewError("given input size (%v) is larger than maximum supported size (%v) for decrypted data", size, config.MaxDecryptedSize)
	}
	result := (size / config.MaxPayloadSize) * config.MaxBufferSize
	remainder := size % config.MaxPayloadSize
	if remainder > 0 {
		result += remainder + (header.HeaderSize + segment.TagSize)
	}
	return result, nil
}

// getDecryptedSize ...
func getDecryptedSize(size uint64) (uint64, error) {
	if size > config.MaxEncryptedSize {
		return 0, stacktrace.NewError("given input size (%v) is larger than maximum supported size (%v) for encrypted data", size, config.MaxEncryptedSize)
	}
	result := (size / config.MaxBufferSize) * config.MaxPayloadSize
	remainder := size % config.MaxBufferSize
	if remainder > 0 {
		if remainder <= header.HeaderSize+segment.TagSize {
			return 0, stacktrace.NewError("invalid decrypted size")
		}
		result += remainder - (header.HeaderSize + segment.TagSize)
	}
	return result, nil
}
