package grpc

import (
	"context"
	model "github.com/da-moon/dare-cli/model"
	stacktrace "github.com/palantir/stacktrace"
)

// Client is an implementation of shared.Encrypt that talks over gRPC.
type Client struct {
	client model.EncryptorClient
}

// Encrypt ...
func (c *Client) Encrypt(req *model.EncryptRequest) (*model.EncryptResponse, error) {
	_resp, err := c.client.Encrypt(context.Background(), req)
	if err != nil {
		err = stacktrace.Propagate(err, "Encrypt call failed with request %#v", req)
	}
	return _resp, err
}
