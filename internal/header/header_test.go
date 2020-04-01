package header_test

import (
	"crypto/rand"
	"io"

	header "github.com/da-moon/coe817-dare/internal/header"
	assert "github.com/stretchr/testify/assert"

	"testing"
)

type TestCase struct {
	length int
	final  bool
}

var tests = []TestCase{
	{
		length: 16,
		final:  false,
	},
	{
		length: 32,
		final:  true,
	},
	{
		length: 64,
		final:  false,
	},
	{
		length: 1024,
		final:  true,
	},
}

func TestHeaderLength(t *testing.T) {
	for _, test := range tests {
		h := header.Header(make([]byte, header.HeaderSize))
		h.SetLength(test.length)
		assert.Equal(t, test.length, h.GetLength())

	}
}
func TestHeaderNonce(t *testing.T) {
	for _, test := range tests {
		h := header.Header(make([]byte, header.HeaderSize))
		h.SetLength(test.length)
		var randVal [header.NonceFieldSize]byte
		_, err := io.ReadFull(rand.Reader, randVal[:])
		assert.Nil(t, err)
		h.SetNonce(randVal[:], test.final)
		if test.final {
			assert.True(t, h.IsFinal())
			assert.Equal(t, randVal[0]&0x7F, h[4])
		} else {
			assert.Equal(t, randVal[0]|0x80, h[4])
			assert.False(t, h.IsFinal())
		}
	}
}
