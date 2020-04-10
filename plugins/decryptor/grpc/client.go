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
func (c *Client) Decrypt(req *model.DecryptRequest) (*model.DecryptResponse, error) {
	_resp, err := c.client.Decrypt(context.Background(), req)
	if err != nil {
		err = stacktrace.Propagate(err, "Decrypt call failed with request %#v", req)
	}
	return _resp, err
}
