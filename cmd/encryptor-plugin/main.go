package main

import (
	"fmt"
	model "github.com/da-moon/coe817-dare/model"
	grpc "github.com/da-moon/coe817-dare/plugins/encryptor/grpc"
	shared "github.com/da-moon/coe817-dare/plugins/shared"
	plugin "github.com/hashicorp/go-plugin"
)

// Encrypt - this is the struct that implements engine operations
type Encrypt struct {
}

// Encrypt - Implementation of Encrypt method for go engine
func (e *Encrypt) Encrypt(req *model.EncryptRequest) (*model.EncryptResponse, error) {
	fmt.Println("ENCRYPT PLUGIN BSCKEND")
	srcHash := &model.Hash{
		Md5:    "[encrypt] src md5hash",
		Sha256: "[encrypt] src sha256hash",
	}
	dstHash := &model.Hash{
		Md5:    "[encrypt] dst md5hash",
		Sha256: "[encrypt] dst sha256hash",
	}
	result := &model.EncryptResponse{
		SourceHash:      srcHash,
		DestinationHash: dstHash,
		RandomNonce:     "adajsiodjaiofaio",
		RandomKey:       "adajsiodjaiofaio",
	}
	fmt.Println(result)
	return result, nil

}
func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"encryptor": &grpc.Plugin{Impl: &Encrypt{}},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
