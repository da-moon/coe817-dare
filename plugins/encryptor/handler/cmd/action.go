package cmd

import (
	shared "github.com/da-moon/coe817-dare/plugins/shared"
	cli "github.com/urfave/cli"
)

// Client - this is the client that talks to engine
var Client shared.EncryptorInterface

// EncryptImpl - this function is called in the implmentation of encrypt cmd operation and should be set by the developer
var EncryptImpl func(ctx *cli.Context) error

// Encrypt - cmd interface implmentation
func Encrypt(ctx *cli.Context) error {
	return EncryptImpl(ctx)
}
