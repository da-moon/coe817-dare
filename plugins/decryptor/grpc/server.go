package grpc

import (
	"context"
	model "github.com/da-moon/dare-cli/model"
	shared "github.com/da-moon/dare-cli/plugins/shared"
	stacktrace "github.com/palantir/stacktrace"
)

// Server - Here is the gRPC server that Client talks to.
type Server struct {
	Impl shared.DecryptorInterface
}

// Decrypt ...
func (s *Server) Decrypt(ctx context.Context, _req *model.DecryptRequest) (*model.DecryptResponse, error) {
	resp, err := s.Impl.Decrypt(_req)
	if err != nil {
		err = stacktrace.Propagate(err, "Decrypt call failed with request %#v", &model.DecryptRequest{
			Source:      _req.Source,
			Destination: _req.Destination,
		})
	}
	return resp, nil
}
