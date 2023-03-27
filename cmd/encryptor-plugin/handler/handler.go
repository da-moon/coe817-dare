package handler

import (
	"crypto/rand"
	"encoding/hex"
	dare "github.com/da-moon/dare-cli"
	model "github.com/da-moon/dare-cli/model"
	hashsink "github.com/da-moon/dare-cli/pkg/hashsink"
	stacktrace "github.com/palantir/stacktrace"
	"io"
	"os"
)

// Encrypt - this is the struct that implements engine operations
type Encrypt struct {
}

// Encrypt - Implementation of Encrypt method for go engine
func (e *Encrypt) Encrypt(req *model.EncryptRequest) (*model.EncryptResponse, error) {
	result := &model.EncryptResponse{}
	nonce, err := dare.RandomNonce()
	if err != nil {
		err = stacktrace.Propagate(err, "could not encrypt data due to failure in generating random nonce")
		return nil, err
	}
	result.RandomNonce = hex.EncodeToString(nonce[:])
	var key [32]byte
	if len(req.Key) == 0 {
		_, err := io.ReadFull(rand.Reader, key[:])
		if err != nil {
			err = stacktrace.Propagate(err, "could not encrypt data due to failure in generating random key")
			return nil, err
		}
		result.RandomKey = hex.EncodeToString(key[:])
	} else {
		decoded, err := hex.DecodeString(req.Key)
		if err != nil {
			err = stacktrace.Propagate(err, "could not encrypt data due to failure in decoding encryption key")
			return nil, err
		}
		if len(decoded) != 32 {
			err = stacktrace.NewError("could not encrypt data since given encoded encryption key is %d bytes. We expect 32 byte keys", len(decoded))
			return nil, err
		}
		copy(key[:], decoded[:32])

	}
	fi, err := os.Stat(req.Source)
	if err == nil {
		if fi.Size() == 0 {
			os.Remove(req.Source)
			err = stacktrace.NewError("decryption failure due to file with empty size at '%v'", req.Source)
			return nil, err
		}
	}
	if err != nil {
		err = stacktrace.Propagate(err, "could not stat src at '$v'", req.Source)
		return nil, err
	}
	srcFile, err := os.Open(req.Source)
	if err != nil {
		err = stacktrace.Propagate(err, "could not encrypt due to failure in opening source file at %s", req.Source)
		return nil, err
	}
	defer srcFile.Close()
	os.Remove(req.Destination)
	destinationFile, err := os.Create(req.Destination)
	if err != nil {
		err = stacktrace.NewError("could not successfully create a new empty file for %s", req.Destination)
		return nil, err
	}
	defer destinationFile.Close()

	dstWriter := hashsink.NewWriter(destinationFile)

	err = dare.EncryptWithWriter(dstWriter, srcFile, key, nonce)
	if err != nil {
		err = stacktrace.Propagate(err, "Could not Encrypt file at '%s' and store it in '%s' ", req.Source, req.Destination)
		return nil, err
	}
	result.OutputHash = &model.Hash{
		Md5:    dstWriter.MD5HexString(),
		Sha256: dstWriter.SHA256HexString(),
	}
	return result, nil

}
