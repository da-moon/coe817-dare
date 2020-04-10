package main

import (
	handler "github.com/da-moon/coe817-dare/cmd/encryptor-plugin/handler"
	grpc "github.com/da-moon/coe817-dare/plugins/encryptor/grpc"
	shared "github.com/da-moon/coe817-dare/plugins/shared"
	plugin "github.com/hashicorp/go-plugin"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"encryptor": &grpc.Plugin{Impl: &handler.Encrypt{}},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
