package cmd

import (
	shared "github.com/da-moon/coe817-dare/plugins/shared"
	cli "github.com/urfave/cli"
)

// Client - this is the client that talks to engine
var Client shared.DecryptorInterface

// DecryptImpl - this function is called in the implmentation of decrypt cmd operation and should be set by the developer
var DecryptImpl func(ctx *cli.Context) error

// Decrypt - cmd interface implmentation
func Decrypt(ctx *cli.Context) error {
	return DecryptImpl(ctx)
}
