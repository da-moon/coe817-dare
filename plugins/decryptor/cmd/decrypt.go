package cmd

import (
	fmt "fmt"
	grpc "github.com/da-moon/coe817-dare/plugins/decryptor/grpc"
	handler "github.com/da-moon/coe817-dare/plugins/decryptor/handler/cmd"
	rpc "github.com/da-moon/coe817-dare/plugins/decryptor/net-rpc"
	shared "github.com/da-moon/coe817-dare/plugins/shared"
	plugin "github.com/hashicorp/go-plugin"
	stacktrace "github.com/palantir/stacktrace"
	cli "github.com/urfave/cli"
	exec "os/exec"
)

// Decrypt - cli command based on service name
var Decrypt = []cli.Command{
	{
		Name:    "decrypt",
		Usage:   "Decryptor Engine",
		Aliases: []string{"decryptor-engine"},
		Flags: cli.FlagsByName{
			cli.StringFlag{
				Name:  "decrypt-binary",
				Usage: "decrypt engine binary path",
				Value: "",
			},
		},
		Before: func(ctx *cli.Context) error {
			// We don't want to see the plugin logs.
			// log.SetOutput(ioutil.Discard)
			// We're a host. Start by launching the plugin process.
			path := ctx.String("decrypt-binary")
			if len(path) == 0 {
				err := stacktrace.NewError("decrypt plugin engine binary path is empty")
				fmt.Println(err)
				return err
			}
			client := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: shared.HandshakeConfig,
				Plugins: map[string]plugin.Plugin{
					"decryptor_grpc": &grpc.Plugin{},
					"decrypt":        &rpc.Plugin{},
				},
				Cmd: exec.Command("sh", "-c", path),
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
			})
			defer client.Kill()
			// Connect via RPC
			rpcClient, err := client.Client()
			if err != nil {
				err = stacktrace.Propagate(err, "Client failed to return the protocol client for decrypt engine connection")
				fmt.Println(err)
				return err
			}
			// Request the plugin
			raw, err := rpcClient.Dispense("decryptor_grpc")
			if err != nil {
				err = stacktrace.Propagate(err, "RPC Client could not dispense a new instance of decryptor_grpc")
				fmt.Println(err)
				return err
			}
			// We should have a decrypt store now! This feels like a normal interface
			// implementation but is in fact over an RPC connection.
			handler.Client = raw.(shared.DecryptorInterface)
			return nil
		},

		Subcommands: []cli.Command{
			DecryptSubcommand,
		},
	},
}

// Decrypt - Decrypt cli subcommand
var DecryptSubcommand = cli.Command{
	Name:    "decrypt",
	Usage:   "",
	Aliases: []string{},
	Action:  handler.Decrypt,
}
