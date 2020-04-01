package dare

import (
	// config "github.com/da-moon/coe817-dare/internal/config"
	decryptor "github.com/da-moon/coe817-dare/internal/decryptor"
	// header "github.com/da-moon/coe817-dare/internal/header"
	// segment "github.com/da-moon/coe817-dare/internal/segment"
	"io"
)

// Decrypt ...
func Decrypt(dst io.Writer, reader io.Reader, key []byte) (n int64, err error) {
	decReader, err := DecryptReader(reader, key)
	if err != nil {
		return 0, err
	}
	return io.Copy(dst, decReader)
}

// DecryptReader ...
func DecryptReader(reader io.Reader, key []byte) (io.Reader, error) {
	return decryptor.New(reader, nil, key)
}
