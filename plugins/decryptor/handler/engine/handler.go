package engine

import (
	model "github.com/da-moon/coe817-dare/model"
	grpc "github.com/da-moon/coe817-dare/plugins/decryptor/grpc"
	shared "github.com/da-moon/coe817-dare/plugins/shared"
	plugin "github.com/hashicorp/go-plugin"
)

// ServeConfig - This is the plugin config thet is used in main function of engine
func ServeConfig() *plugin.ServeConfig {
	return &plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"decrypt": &grpc.Plugin{Impl: &Decrypt{}},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	}
}

// Decrypt - this is the struct that implements engine operations
type Decrypt struct{}

// DecryptImpl - this function is called in the implmentation of decrypt operation and should be set by the developer
var DecryptImpl func(source string, destination string) (*model.Hash, *model.Hash, error)

// Decrypt - Implementation of Decrypt method for go engine
func (Decrypt) Decrypt(source string, destination string) (*model.Hash, *model.Hash, error) {
	return DecryptImpl(source, destination)
}
