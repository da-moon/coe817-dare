package dare_test

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	dare "github.com/da-moon/coe817-dare"
	config "github.com/da-moon/coe817-dare/internal/config"
	log "github.com/da-moon/coe817-dare/pkg/log"
	assert "github.com/stretchr/testify/assert"
	hkdf "golang.org/x/crypto/hkdf"
	"io"
	"testing"
)

type TestCase struct {
	datasize    int
	buffersize  int
	payloadsize int
}

func TestBasicEncrypt(t *testing.T) {
	log.SetTestLogger(t)
	tests := []TestCase{
		{
			datasize: config.MaxBufferSize / 2,
		},
		{
			datasize: config.MaxBufferSize,
		},
		{
			datasize: config.MaxBufferSize + 1,
		},
		{
			datasize: 2 * config.MaxBufferSize,
		},
		{
			datasize: 2*config.MaxBufferSize + 1,
		},
		// {
		// 	datasize: 2*config.MaxBufferSize + 6,
		// 	// buffersize:  config.MaxPayloadSize + 1,
		// 	// payloadsize: config.MaxPayloadSize,
		// },

		// {
		// 	datasize: (2 * config.MaxPayloadSize),
		// 	// buffersize:  config.MaxPayloadSize + 1,
		// 	// payloadsize: config.MaxPayloadSize,
		// },
		// {
		// 	datasize: 1024*1024 + 1,
		// 	// buffersize:  config.MaxPayloadSize + 1,
		// 	// payloadsize: config.MaxPayloadSize,
		// },
		// {
		// 	datasize: 3*config.MaxPayloadSize + 1,
		// 	// buffersize:  3 * config.MaxPayloadSize,
		// 	// payloadsize: config.MaxPayloadSize,
		// },
		// {
		// 	datasize: 1024 * 1024,
		// 	// buffersize:  2 * 1024 * 1024,
		// 	// payloadsize: config.MaxPayloadSize,
		// },
		// {
		// 	datasize:    (3 * config.MaxPayloadSize) + (config.MaxPayloadSize / 2),
		// 	buffersize:  config.MaxPayloadSize + 1,
		// 	payloadsize: config.MaxPayloadSize,
		// },
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			var (
				nonce [32]byte
				key   [32]byte
			)
			keyString, err := randomHex(26)
			assert.NoError(t, err)
			masterkey, err := hex.DecodeString(keyString)
			assert.NoError(t, err, "Cannot decode hex key")
			_, err = io.ReadFull(rand.Reader, nonce[:])
			assert.NoError(t, err, "failed to generate random data for nonce")
			// driving master key ...
			kdf := hkdf.New(sha256.New, masterkey, nonce[:], nil)
			_, err = io.ReadFull(kdf, key[:])
			assert.NoError(t, err, "could not drive an encryption key. masterkey=%v nonce=%v", masterkey, nonce[:])
			data := make([]byte, test.datasize)
			_, err = io.ReadFull(rand.Reader, data)
			assert.NoError(t, err, "could not generate random data for encryption")
			output := bytes.NewBuffer(nil)
			_, err = dare.Encrypt(
				output,
				bytes.NewReader(data),
				key[:],
			)
			assert.NoError(t, err, "could not encrypt data")

			decrypted := bytes.NewBuffer(nil)
			n, err := dare.Decrypt(
				decrypted,
				output,
				key[:],
			)
			assert.NoError(t, err)
			assert.Equal(t, int64(test.datasize), n, "decrypt expected read=%v actual read=%v", int64(test.datasize), n)
			assert.True(t, bytes.Equal(data, decrypted.Bytes()))
			// if !bytes.Equal(data, decrypted.Bytes()) {
			// 	t.Errorf("Failed to encrypt and decrypt data")
			// }
		})
	}
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
