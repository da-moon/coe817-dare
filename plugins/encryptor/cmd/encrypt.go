package cmd

import (
	fmt "fmt"
	grpc "github.com/da-moon/coe817-dare/plugins/encryptor/grpc"
	handler "github.com/da-moon/coe817-dare/plugins/encryptor/handler/cmd"
	rpc "github.com/da-moon/coe817-dare/plugins/encryptor/net-rpc"
	shared "github.com/da-moon/coe817-dare/plugins/shared"
	plugin "github.com/hashicorp/go-plugin"
	stacktrace "github.com/palantir/stacktrace"
	cli "github.com/urfave/cli"
	exec "os/exec"
)

// Encrypt - cli command based on service name
var Encrypt = []cli.Command{
	{
		Name:    "encrypt",
		Usage:   "Encryptor Engine ",
		Aliases: []string{"encryptor-engine"},
		Flags: cli.FlagsByName{
			cli.StringFlag{
				Name:  "encrypt-binary",
				Usage: "encrypt engine binary path",
				Value: "",
			},
		},
		Before: func(ctx *cli.Context) error {
			// We don't want to see the plugin logs.
			// log.SetOutput(ioutil.Discard)
			// We're a host. Start by launching the plugin process.
			path := ctx.String("encrypt-binary")
			if len(path) == 0 {
				err := stacktrace.NewError("encrypt plugin engine binary path is empty")
				fmt.Println(err)
				return err
			}
			client := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: shared.HandshakeConfig,
				Plugins: map[string]plugin.Plugin{
					"encryptor_grpc": &grpc.Plugin{},
					"encrypt":        &rpc.Plugin{},
				},
				Cmd: exec.Command("sh", "-c", path),
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
			})
			defer client.Kill()
			// Connect via RPC
			rpcClient, err := client.Client()
			if err != nil {
				err = stacktrace.Propagate(err, "Client failed to return the protocol client for encrypt engine connection")
				fmt.Println(err)
				return err
			}
			// Request the plugin
			raw, err := rpcClient.Dispense("encryptor_grpc")
			if err != nil {
				err = stacktrace.Propagate(err, "RPC Client could not dispense a new instance of encryptor_grpc")
				fmt.Println(err)
				return err
			}
			// We should have a encrypt store now! This feels like a normal interface
			// implementation but is in fact over an RPC connection.
			handler.Client = raw.(shared.EncryptorInterface)
			return nil
		},

		Subcommands: []cli.Command{
			EncryptSubcommand,
		},
	},
}

// EncryptSubcommand - Encrypt cli subcommand
var EncryptSubcommand = cli.Command{
	Name:    "encrypt",
	Usage:   "",
	Aliases: []string{},
	Action:  handler.Encrypt,
}
