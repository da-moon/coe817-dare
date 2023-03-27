package grpc

import (
	"context"
	model "github.com/da-moon/dare-cli/model"
	shared "github.com/da-moon/dare-cli/plugins/shared"
	stacktrace "github.com/palantir/stacktrace"
)

// Server - Here is the gRPC server that Client talks to.
type Server struct {
	Impl shared.EncryptorInterface
}

// Encrypt ...
func (s *Server) Encrypt(ctx context.Context, _req *model.EncryptRequest) (*model.EncryptResponse, error) {
	resp, err := s.Impl.Encrypt(_req)
	if err != nil {
		err = stacktrace.Propagate(err, "Encrypt call failed with request %#v", &model.EncryptRequest{
			Source:      _req.Source,
			Destination: _req.Destination,
		})
	}
	return resp, nil
}
