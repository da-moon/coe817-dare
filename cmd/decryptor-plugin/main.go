package main

import (
	model "github.com/da-moon/coe817-dare/model"
	grpc "github.com/da-moon/coe817-dare/plugins/decryptor/grpc"
	shared "github.com/da-moon/coe817-dare/plugins/shared"
	plugin "github.com/hashicorp/go-plugin"
)

// Decrypt - this is the struct that implements engine operations
type Decrypt struct{}

// Decrypt - Implementation of Decrypt method for go engine
func (Decrypt) Decrypt(source string, destination string) (*model.Hash, *model.Hash, error) {
	srcHash := &model.Hash{
		Md5:    "[Decrypt] src md5hash",
		Sha256: "[Decrypt] src sha256hash",
	}
	dstHash := &model.Hash{
		Md5:    "[Decrypt] dst md5hash",
		Sha256: "[Decrypt] dst sha256hash",
	}
	return srcHash, dstHash, nil
}

// ServeConfig - This is the plugin config thet is used in main function of engine
func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"decryptor": &grpc.Plugin{Impl: &Decrypt{}},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
