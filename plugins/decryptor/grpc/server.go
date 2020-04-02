package grpc

import (
	"context"
	model "github.com/da-moon/coe817-dare/model"
	shared "github.com/da-moon/coe817-dare/plugins/shared"
	stacktrace "github.com/palantir/stacktrace"
)

// Server - Here is the gRPC server that Client talks to.
type Server struct {
	Impl shared.DecryptorInterface
}

// Decrypt ...
func (s *Server) Decrypt(ctx context.Context, _req *model.DecryptRequest) (*model.DecryptResponse, error) {
	sourceHash, destinationHash, err := s.Impl.Decrypt(
		_req.Source,
		_req.Destination,
	)
	if err != nil {
		err = stacktrace.Propagate(err, "Decrypt call failed with request %#v", &model.DecryptRequest{
			Source:      _req.Source,
			Destination: _req.Destination,
		})
	}
	return &model.DecryptResponse{
		SourceHash:      sourceHash,
		DestinationHash: destinationHash,
	}, nil
}
