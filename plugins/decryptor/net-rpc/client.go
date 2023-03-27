package netrpc

import (
	model "github.com/da-moon/dare-cli/model"
	stacktrace "github.com/palantir/stacktrace"
	rpc "net/rpc"
)

// Client is an implementation of shared that talks over RPC.
type Client struct{ client *rpc.Client }

// Decrypt ...
func (c *Client) Decrypt(req *model.DecryptRequest) (*model.DecryptResponse, error) {
	var _resp model.DecryptResponse
	err := c.client.Call("Plugin.Decrypt", req, &_resp)
	if err != nil {
		err = stacktrace.Propagate(err, "Decrypt call failed with request %#v", req)
	}
	return &_resp, err
}
