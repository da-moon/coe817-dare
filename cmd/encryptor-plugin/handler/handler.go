package handler

import (
	"crypto/rand"
	// "encoding/base64"
	"encoding/hex"
	"fmt"
	dare "github.com/da-moon/coe817-dare"
	model "github.com/da-moon/coe817-dare/model"
	// hashsink "github.com/da-moon/coe817-dare/pkg/hashsink"
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
	fmt.Printf("encrypt backend nonce %v\n", nonce)
	// result.RandomNonce = base64.StdEncoding.EncodeToString(nonce[:])
	result.RandomNonce = hex.EncodeToString(nonce[:])
	var key [32]byte
	if len(req.Key) == 0 {
		_, err := io.ReadFull(rand.Reader, key[:])
		if err != nil {
			err = stacktrace.Propagate(err, "could not encrypt data due to failure in generating random key")
			return nil, err
		}
		// result.RandomKey = base64.StdEncoding.EncodeToString(key[:])
		result.RandomKey = hex.EncodeToString(key[:])
		fmt.Printf("encrypt backend rand key %v\n", key)
	} else {
		// decoded, err := base64.StdEncoding.DecodeString(req.Key)
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
		fmt.Printf("encrypt backend given key %v\n", key)

	}
	fi, err := os.Stat(req.Source)
	if err == nil {
		if fi.Size() == 0 {
			os.Remove(req.Source)
			err = stacktrace.NewError("decryption failure due to file with empty size at '%v'", req.Source)
			fmt.Printf("%v\n", err.Error())
			return nil, err
		}
	}
	if err != nil {
		err = stacktrace.Propagate(err, "could not stat src at '$v'", req.Source)
		fmt.Printf("%v\n", err.Error())
		return nil, err
	}
	fmt.Printf("encrypt backend stated")
	srcFile, err := os.Open(req.Source)
	if err != nil {
		err = stacktrace.Propagate(err, "could not encrypt due to failure in opening source file at %s", req.Source)
		return nil, err
	}
	defer srcFile.Close()
	dstFile, err := os.OpenFile(
		req.Destination,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0600)

	if dstFile == nil {
		err = stacktrace.NewError("could not successfully get a file handle for %s", req.Destination)
		return nil, err
	}
	defer dstFile.Close()
	if err != nil {
		err = stacktrace.Propagate(err, "Could not create empty file at (%s) ", req.Destination)
		return nil, err
	}
	fmt.Printf("about to encrypt")

	// srcReader := hashsink.NewReader(srcFile, 0)
	// dstWriter := hashsink.NewWriter(dstFile)

	_, err = dare.Encrypt(
		dstFile,
		srcFile,
		key,
		nonce,
	)
	if err != nil {
		err = stacktrace.Propagate(err, "Could not Encrypt file at '%s' and store it in '%s' ", req.Source, req.Destination)
		return nil, err
	}
	fmt.Printf("about to Sync")
	err = dstFile.Sync()
	if err != nil {
		err = stacktrace.Propagate(err, "Could not Encrypt file at '%s' and store it in '%s' due to flush to disk failure", req.Source, req.Destination)
		return nil, err
	}
	fmt.Printf("about to return ")
	// result.SourceHash = &model.Hash{
	// 	// Md5:    srcReader.MD5HexString(),
	// 	// Sha256: srcReader.SHA256HexString(),
	// }
	// result.DestinationHash = &model.Hash{
	// 	// Md5:    srcReader.MD5HexString(),
	// 	// Sha256: srcReader.SHA256HexString(),
	// }
	return result, nil

}
