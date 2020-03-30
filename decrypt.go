package dare

import (
	// config "github.com/da-moon/coe817-dare/pkg/config"
	decryptor "github.com/da-moon/coe817-dare/pkg/decryptor"
	// header "github.com/da-moon/coe817-dare/pkg/header"
	// segment "github.com/da-moon/coe817-dare/pkg/segment"
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
