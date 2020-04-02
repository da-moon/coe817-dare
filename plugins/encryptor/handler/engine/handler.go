package engine

import (
	model "github.com/da-moon/coe817-dare/model"
	grpc "github.com/da-moon/coe817-dare/plugins/encryptor/grpc"
	shared "github.com/da-moon/coe817-dare/plugins/shared"
	plugin "github.com/hashicorp/go-plugin"
)

// ServeConfig - This is the plugin config thet is used in main function of engine
func ServeConfig() *plugin.ServeConfig {
	return &plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"encrypt": &grpc.Plugin{Impl: &Encrypt{}},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	}
}

// Encrypt - this is the struct that implements engine operations
type Encrypt struct{}

// EncryptImpl - this function is called in the implmentation of encrypt operation and should be set by the developer
var EncryptImpl func(source string, destination string) (*model.Hash, *model.Hash, error)

// Encrypt - Implementation of Encrypt method for go engine
func (Encrypt) Encrypt(source string, destination string) (*model.Hash, *model.Hash, error) {
	return EncryptImpl(source, destination)
}
