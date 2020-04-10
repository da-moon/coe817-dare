package dare

import (
	decryptor "github.com/da-moon/coe817-dare/pkg/dare/decryptor"
	stacktrace "github.com/palantir/stacktrace"
	"io"
	"os"
)

// DecryptFile a convenience method that uses decrypt reader
// to decrypt data as it is reading it from a source file path
// and would write the decrypted data to a target file path
func DecryptFile(source string, destination string, key [32]byte, nonce [24]byte) (int64, error) {
	fi, err := os.Stat(source)
	if err == nil {
		if fi.Size() == 0 {
			os.Remove(source)
			return 0, nil
		}
	}
	srcFile, err := os.Open(source)
	if srcFile != nil {
		defer srcFile.Close()
	}
	if err != nil {
		err = stacktrace.Propagate(err, "could not decrypt due to failure in opening source file at %s", source)
		return 0, err
	}
	dstFile, err := os.OpenFile(
		destination,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0600)

	if dstFile == nil {
		err = stacktrace.NewError("could not successfully get a file handle for %s", destination)
		return 0, err
	}

	if dstFile != nil {
		defer dstFile.Close()
	}
	if err != nil {
		err = stacktrace.Propagate(err, "Could not create empty file at (%s) ", destination)
		return 0, err
	}
	n, err := Decrypt(
		dstFile,
		srcFile,
		key,
		nonce,
	)
	if err != nil {
		err = stacktrace.Propagate(err, "Could not Decrypt file at '%s' and store it in '%s' ", source, destination)
		return 0, err

	}
	return n, nil

}

// Decrypt ...
func Decrypt(writer io.Writer, reader io.Reader, key [32]byte, nonce [24]byte) (int64, error) {
	decryptReader := DecryptReader(reader, key, nonce)
	n, err := io.Copy(writer, decryptReader)
	if err != nil {
		if err != io.EOF {
			return 0, err
		}
	}
	return n, nil
}

// DecryptReader returns an io.reader that
// Decrypts data with a passed key as it is reading it
// from an io stream (eg socket , file).
func DecryptReader(reader io.Reader, key [32]byte, nonce [24]byte) io.Reader {
	return decryptor.NewReader(reader, nonce, &key)
}
