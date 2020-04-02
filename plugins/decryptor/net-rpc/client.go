package netrpc

import (
	model "github.com/da-moon/coe817-dare/model"
	stacktrace "github.com/palantir/stacktrace"
	rpc "net/rpc"
)

// Client is an implementation of shared that talks over RPC.
type Client struct{ client *rpc.Client }

// Decrypt ...
func (c *Client) Decrypt(source string, destination string) (*model.Hash, *model.Hash, error) {
	var _resp model.DecryptResponse
	err := c.client.Call("Plugin.Decrypt", &model.DecryptRequest{
		Source:      source,
		Destination: destination,
	}, &_resp)
	if err != nil {
		err = stacktrace.Propagate(err, "Decrypt call failed with request %#v", &model.DecryptRequest{
			Source:      source,
			Destination: destination,
		})
	}
	return _resp.SourceHash, _resp.DestinationHash, err
}
