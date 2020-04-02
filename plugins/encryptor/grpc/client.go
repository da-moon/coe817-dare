package grpc

import (
	"context"
	model "github.com/da-moon/coe817-dare/model"
	stacktrace "github.com/palantir/stacktrace"
)

// Client is an implementation of shared.Encrypt that talks over gRPC.
type Client struct {
	client model.EncryptorClient
}

// Encrypt ...
func (c *Client) Encrypt(source string, destination string) (*model.Hash, *model.Hash, error) {
	_resp, err := c.client.Encrypt(context.Background(), &model.EncryptRequest{
		Source:      source,
		Destination: destination,
	})
	if err != nil {
		err = stacktrace.Propagate(err, "Encrypt call failed with request %#v", &model.EncryptRequest{
			Source:      source,
			Destination: destination,
		})
	}
	return _resp.SourceHash, _resp.DestinationHash, err
}
