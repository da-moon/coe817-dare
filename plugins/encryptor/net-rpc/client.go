package netrpc

import (
	model "github.com/da-moon/coe817-dare/model"
	stacktrace "github.com/palantir/stacktrace"
	rpc "net/rpc"
)

// Client is an implementation of shared that talks over RPC.
type Client struct{ client *rpc.Client }

// Encrypt ...
func (c *Client) Encrypt(source string, destination string) (*model.Hash, *model.Hash, error) {
	var _resp model.EncryptResponse
	err := c.client.Call("Plugin.Encrypt", &model.EncryptRequest{
		Source:      source,
		Destination: destination,
	}, &_resp)
	if err != nil {
		err = stacktrace.Propagate(err, "Encrypt call failed with request %#v", &model.EncryptRequest{
			Source:      source,
			Destination: destination,
		})
	}
	return _resp.SourceHash, _resp.DestinationHash, err
}
