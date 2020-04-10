package handler

import (
	// "encoding/base64"
	"encoding/hex"
	"fmt"
	dare "github.com/da-moon/coe817-dare"
	model "github.com/da-moon/coe817-dare/model"
	// hashsink "github.com/da-moon/coe817-dare/pkg/hashsink"
	stacktrace "github.com/palantir/stacktrace"
	"os"
)

// Decrypt - this is the struct that implements engine operations
type Decrypt struct{}

// Decrypt - Implementation of Decrypt method for go engine
func (Decrypt) Decrypt(req *model.DecryptRequest) (*model.DecryptResponse, error) {
	result := &model.DecryptResponse{}
	fmt.Printf("decrypt backend src %v \n", req.Source)
	fmt.Printf("decrypt backend dst %v \n", req.Destination)
	fmt.Printf("decrypt backend key %v \n", req.Key)
	// decKey, err := base64.StdEncoding.DecodeString(req.Key)
	decKey, err := hex.DecodeString(req.Key)
	if err != nil {
		err = stacktrace.Propagate(err, "could not decrypt file '%s' due to failure in decoding encryption key from hex string '%s'", req.Source, req.Key)
		return nil, err
	}
	if len(decKey) != 32 {
		err = stacktrace.NewError("could not decrypt data since given encoded encryption key is %d bytes. We expect 32 byte keys", len(decKey))
		return nil, err

	}
	var key [32]byte
	copy(key[:], decKey[:32])
	fmt.Printf("decrypt backend decKey key %v \n", key)
	// decNonce, err := base64.StdEncoding.DecodeString(req.Nonce)
	fmt.Printf("decrypt backend req nonce  %v \n", req.Nonce)
	decNonce, err := hex.DecodeString(req.Nonce)
	if err != nil {
		err = stacktrace.Propagate(err, "could not decrypt file '%s' due to failure in decoding nonce from hex string '%s'", req.Source, req.Nonce)
		return nil, err
	}
	if len(decNonce) != 24 {
		err = stacktrace.NewError("could not decrypt data since given encoded nonce is %d bytes. We expect 24 byte keys", len(decNonce))
		return nil, err

	}

	var nonce [24]byte
	copy(nonce[:], decNonce[:24])

	fmt.Printf("decrypt backend req decNonce  %v \n", nonce)
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
		fmt.Printf("%v\n", err.Error())
		return nil, err
	}

	fmt.Printf("decrypt backend src stated\n")
	encSourceFile, err := os.Open(req.Source)
	if err != nil {
		err = stacktrace.Propagate(err, "could not decrypt due to failure in opening source file at %s", req.Source)
		return nil, err
	}
	defer encSourceFile.Close()
	decDestinationFile, err := os.Create(req.Destination)
	if err != nil {
		err = stacktrace.NewError("could not successfully create a new empty file for %s", req.Destination)
		return nil, err
	}
	defer decDestinationFile.Close()
	fmt.Printf("about to decrypt\n")
	err = dare.DecryptWithWriter(decDestinationFile, encSourceFile, key, nonce)

	// srcReader := hashsink.NewReader(srcFile, 0)
	// dstWriter := hashsink.NewWriter(dstFile)
	// _, err = dare.Decrypt(
	// 	dstFile,
	// 	srcReader,
	// 	key,
	// 	nonce,
	// )
	// fmt.Printf("decrypt backend dstFile opened")

	if err != nil {
		err = stacktrace.Propagate(err, "Could not Decrypt file at '%s' and store it in '%s' with key '%s' and nonce '%s'", req.Source, req.Destination, req.Key, req.Nonce)
		return nil, err

	}
	// result.SourceHash = &model.Hash{
	// 	Md5:    srcReader.MD5HexString(),
	// 	Sha256: srcReader.SHA256HexString(),
	// }
	// result.DestinationHash = &model.Hash{
	// 	Md5:    dstWriter.MD5HexString(),
	// 	Sha256: dstWriter.SHA256HexString(),
	// }
	return result, nil
}
