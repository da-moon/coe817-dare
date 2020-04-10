package dare

import (
	"io"

	encryptor "github.com/da-moon/coe817-dare/pkg/dare/encryptor"
	stacktrace "github.com/palantir/stacktrace"
)

// Encrypt a convenience method that uses encrypt reader
// to encrypt data as it is reading it from io stream
// and would write the encrypted data to another sink
// through an io.writer
func Encrypt(writer io.Writer, reader io.Reader, key [32]byte, nonce [24]byte) (int64, error) {
	fmt.Printf("[TRACE] dare.Encrypt called")
	encryptReader := EncryptReader(reader, key, nonce)
	fmt.Printf("[TRACE] dare.Encrypt about to io.Copy")
	n, err := io.Copy(writer, encryptReader)
	if err != nil {
		if err != io.EOF {
			err = stacktrace.Propagate(err, "encryption failed due to an issue with io.copy")
			fmt.Printf("[TRACE] io.copy err", err.Error())
			return 0, err
		}
	}
	return n, nil
}

// EncryptReader returns an io.reader that
// encrypts data with a passed key as it is reading it
// from an io stream (eg socket , file).
func EncryptReader(reader io.Reader, key [32]byte, nonce [24]byte) io.Reader {
	return encryptor.NewReader(reader, nonce, &key)
}
