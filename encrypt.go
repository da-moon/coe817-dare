package dare

import (
	// "fmt"
	config "github.com/da-moon/coe817-dare/pkg/dare/config"
	encryptor "github.com/da-moon/coe817-dare/pkg/dare/encryptor"
	// stacktrace "github.com/palantir/stacktrace"
	"io"
)

// EncryptWithWriter encrypts data with a passed key as it is writing it
// to an io stream (eg socket , file).
func EncryptWithWriter(
	dstwriter io.Writer,
	srcReader io.Reader,
	key [32]byte,
	nonce [24]byte,
) error {
	encWriter := encryptor.NewWriter(dstwriter, nonce, &key)
	for {
		buffer := make([]byte, config.DefaultChunkSize)
		bytesRead, err := srcReader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		encWriter.Write(buffer[:bytesRead])
	}
	return nil
}
