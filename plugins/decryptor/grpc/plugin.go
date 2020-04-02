package grpc

import (
	"context"
	model "github.com/da-moon/coe817-dare/model"
	shared "github.com/da-moon/coe817-dare/plugins/shared"
	plugin "github.com/hashicorp/go-plugin"
	grpcx "google.golang.org/grpc"
)

// GRPCClient is an implementation of shared that talks over gRPC.
type Plugin struct {
	// Plugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl shared.DecryptorInterface
}

// GRPCClient - Required method to implement Plugin interface
func (p *Plugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpcx.ClientConn) (interface{}, error) {
	return &Client{client: model.NewDecryptorClient(c)}, nil
}

// GRPCServer - Required method to implement Plugin interface
func (p *Plugin) GRPCServer(broker *plugin.GRPCBroker, s *grpcx.Server) error {
	model.RegisterDecryptorServer(s, &Server{Impl: p.Impl})
	return nil
}
