package shared

import (
	model "github.com/da-moon/coe817-dare/model"
	plugin "github.com/hashicorp/go-plugin"
)

// EncryptorInterface - this is the interface that we're exposing as a plugin.
type EncryptorInterface interface {
	Encrypt(req *model.EncryptRequest) (*model.EncryptResponse, error)
}

// DecryptorInterface - this is the interface that we're exposing as a plugin.
type DecryptorInterface interface {
	Decrypt(req *model.DecryptRequest) (*model.DecryptResponse, error)
}

// HandshakeConfig - engine-interface handshake configuration
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}
