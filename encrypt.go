package dare

import (
	encryptor "github.com/da-moon/coe817-dare/pkg/dare/encryptor"
	stacktrace "github.com/palantir/stacktrace"
	"io"
	"os"
)

// EncryptFile a convenience method that uses encrypt reader
// to encrypt data as it is reading it from a source file path
// and would write the encrypted data to a target file path
func EncryptFile(source string, destination string, key [32]byte) (int64, [24]byte, error) {
	fi, err := os.Stat(source)
	if err == nil {
		if fi.Size() == 0 {
			os.Remove(source)
			return 0, [24]byte{}, nil
		}
	}
	srcFile, err := os.Open(source)
	if srcFile != nil {
		defer srcFile.Close()
	}
	if err != nil {
		err = stacktrace.Propagate(err, "could not encrypt due to failure in opening source file at %s", source)
		return 0, [24]byte{}, err
	}
	dstFile, err := os.OpenFile(
		destination,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0600)

	if dstFile == nil {
		err = stacktrace.NewError("could not successfully get a file handle for %s", destination)
		return 0, [24]byte{}, err
	}

	if dstFile != nil {
		defer dstFile.Close()
	}
	if err != nil {
		err = stacktrace.Propagate(err, "Could not create empty file at (%s) ", destination)
		return 0, [24]byte{}, err
	}
	n, nonce, err := Encrypt(
		dstFile,
		srcFile,
		key,
	)
	if err != nil {
		err = stacktrace.Propagate(err, "Could not Encrypt file at '%s' and store it in '%s' ", source, destination)
		return 0, [24]byte{}, err

	}
	return n, nonce, nil

}

// Encrypt a convenience method that uses encrypt reader
// to encrypt data as it is reading it from io stream
// and would write the encrypted data to another sink
// through an io.writer
func Encrypt(writer io.Writer, reader io.Reader, key [32]byte) (int64, [24]byte, error) {
	encryptReader, nonce, err := EncryptReader(reader, key)
	if err != nil {
		return 0, [24]byte{}, err
	}
	n, err := io.Copy(writer, encryptReader)
	if err != nil {
		if err != io.EOF {
			return 0, [24]byte{}, err
		}
	}
	return n, nonce, nil
}

// EncryptReader returns an io.reader that
// encrypts data with a passed key as it is reading it
// from an io stream (eg socket , file).
func EncryptReader(reader io.Reader, key [32]byte) (io.Reader, [24]byte, error) {
	nonce, err := RandomNonce()
	if err != nil {
		err = stacktrace.Propagate(err, "could not create an encrypted io reader due to failure in generating random nonce")
		return nil, [24]byte{}, err
	}
	return encryptor.NewReader(reader, nonce, &key), nonce, nil
}
