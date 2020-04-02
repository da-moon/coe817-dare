package grpc

import (
	"context"
	model "github.com/da-moon/coe817-dare/model"
	stacktrace "github.com/palantir/stacktrace"
)

// Client is an implementation of shared.Decrypt that talks over gRPC.
type Client struct {
	client model.DecryptorClient
}

// Decrypt ...
func (c *Client) Decrypt(source string, destination string) (*model.Hash, *model.Hash, error) {
	_resp, err := c.client.Decrypt(context.Background(), &model.DecryptRequest{
		Source:      source,
		Destination: destination,
	})
	if err != nil {
		err = stacktrace.Propagate(err, "Decrypt call failed with request %#v", &model.DecryptRequest{
			Source:      source,
			Destination: destination,
		})
	}
	return _resp.SourceHash, _resp.DestinationHash, err
}
