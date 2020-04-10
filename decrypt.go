package dare

import (
	config "github.com/da-moon/coe817-dare/pkg/dare/config"
	decryptor "github.com/da-moon/coe817-dare/pkg/dare/decryptor"
	// stacktrace "github.com/palantir/stacktrace"
	"io"
	// "os"
)

// // DecryptReader returns an io.reader that
// // Decrypts data with a passed key as it is reading it
// // from an io stream (eg socket , file).
// func DecryptReader(reader io.Reader, key [32]byte, nonce [24]byte) io.Reader {
// 	return decryptor.NewReader(reader, nonce, &key)
// }
func DecryptWithWriter(
	dstwriter io.Writer,
	srcReader io.Reader,
	key [32]byte,
	nonce [24]byte,
) error {
	decWriter := decryptor.NewWriter(dstwriter, nonce, &key)
	for {
		buffer := make([]byte, config.DefaultChunkSize+config.DefaultOverhead)
		bytesRead, err := srcReader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		decWriter.Write(buffer[:bytesRead])
	}
	return nil
}
