package daemon

import (
	model "github.com/da-moon/dare-cli/model"
	decryptorGrpc "github.com/da-moon/dare-cli/plugins/decryptor/grpc"
	decryptorRpc "github.com/da-moon/dare-cli/plugins/decryptor/net-rpc"
	encryptorGrpc "github.com/da-moon/dare-cli/plugins/encryptor/grpc"
	encryptorRpc "github.com/da-moon/dare-cli/plugins/encryptor/net-rpc"
	shared "github.com/da-moon/dare-cli/plugins/shared"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	stacktrace "github.com/palantir/stacktrace"
	"log"
	"net/http"
	"os/exec"
)

// Service ...
type Service struct {
	logger       *log.Logger
	pluginLogger hclog.Logger
	encryptor    string
	decryptor    string
	dev          bool
}

// Encrypt json rpc2.0 call that pparses variables and passes it to
// Encrypt Plugin. Sample query
// jq -n \
//   --arg source "/tmp/plain" \
//   --arg destination "/tmp/encrypted" \
//   --arg key "63b76723eb3f9d4f4862b73ff7e39b93c4de129feb4885f1f3feb74dd456e3a5" \
//   --arg id "1" \
//   --arg method "Service.Encrypt" \
//  '{"jsonrpc": "2.0", "method":$method,"params":{"source": $source, "destination":$destination,"key":$key},"id": $id}' | curl \
//     -X POST  \
//  	--silent \
//     --header "Authorization: 12445" \
// 	--header "Content-type: application/json" \
//     --data @- \
//     http://127.0.0.1:8081/rpc  | jq -r
func (s *Service) Encrypt(r *http.Request, req *model.EncryptRequest, res *model.EncryptResponse) error {
	s.logger.Printf("[INFO] daemon-service: Encrypt Called")
	// We don't want to see the plugin logs.
	log.SetOutput(s.logger.Writer())

	// We're a host. Start by launching the plugin process.
	path := s.encryptor
	if len(path) == 0 {
		err := stacktrace.NewError("encryptor plugin engine binary path is empty")
		return err
	}
	s.logger.Printf("[DEBUG] encryptor path is %s", path)
	// todo update handshake config dynamically
	client := plugin.NewClient(&plugin.ClientConfig{
		Logger:          s.pluginLogger,
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"encryptor_grpc": &encryptorGrpc.Plugin{},
			"encryptor":      &encryptorRpc.Plugin{},
		},
		Cmd: exec.Command("sh", "-c", path),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})

	defer client.Kill()
	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		err = stacktrace.Propagate(err, "failed to return the protocol client for encrypt engine connection")
		return err
	}
	s.logger.Printf("[DEBUG] encryptor plugin client is ready")

	// Request the plugin
	raw, err := rpcClient.Dispense("encryptor_grpc")
	if err != nil {
		err = stacktrace.Propagate(err, "RPC Client could not dispense a new instance of encryptor grpc")
		return err
	}
	s.logger.Printf("[DEBUG] a new instance of encryptor grpc was dispenced")

	// We should have a encrypt store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	encryptor, ok := raw.(shared.EncryptorInterface)
	if !ok {
		err = stacktrace.NewError("could not convert the dispensed interface to implementation")
		s.logger.Printf("[DEBUG] error : %#v", err)
		return err

	}
	s.logger.Printf("[DEBUG] encryptor interface converted to impl")
	s.logger.Printf("[DEBUG] encryptor recieved req %#v", req)
	resTmp, err := encryptor.Encrypt(req)
	if err != nil {
		err = stacktrace.Propagate(err, "encryptor failed to encrypt given input")
		s.logger.Printf("[DEBUG] error : %#v", err)
		return err
	}
	res.OutputHash = resTmp.OutputHash
	res.RandomNonce = resTmp.RandomNonce
	res.RandomKey = resTmp.RandomKey
	s.logger.Printf("[DEBUG] encryptor plugin response is %#v ", res)
	return nil
}

// Decrypt json rpc2.0 call that parses variables and passes it to
// Decrypt Plugin. Sample query
// jq -n \
//   --arg source "/tmp/encrypted" \
//   --arg destination "/tmp/decrypted" \
// 	 --arg nonce "e12ffdfa6cb6e56238935e32604cfa5538d3ad51a3542daa" \
//   --arg key "63b76723eb3f9d4f4862b73ff7e39b93c4de129feb4885f1f3feb74dd456e3a5" \
//   --arg id "2" \
//   --arg method "Service.Decrypt" \
//  '{"jsonrpc": "2.0", "method":$method,"params":{"source": $source, "destination":$destination, "nonce":$nonce, "key":$key},"id": $id}' | curl \
// 	-X POST  \
// 	--silent \
//     --header "Authorization: 12445" \
//     --header "Content-type: application/json" \
//     --data @- \
//     http://127.0.0.1:8081/rpc | jq -r
func (s *Service) Decrypt(r *http.Request, req *model.DecryptRequest, res *model.DecryptResponse) error {
	s.logger.Printf("[INFO] daemon-service: Decrypt Called")
	log.SetOutput(s.logger.Writer())
	path := s.decryptor
	if len(path) == 0 {
		err := stacktrace.NewError("decryptor plugin engine binary path is empty")
		return err
	}
	s.logger.Printf("[DEBUG] decryptor path is %s", path)
	// todo update handshake config dynamically
	client := plugin.NewClient(&plugin.ClientConfig{
		Logger:          s.pluginLogger,
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"decryptor_grpc": &decryptorGrpc.Plugin{},
			"decryptor":      &decryptorRpc.Plugin{},
		},
		Cmd: exec.Command("sh", "-c", path),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	defer client.Kill()
	_, err := client.Start()
	if err != nil {
		err = stacktrace.Propagate(err, "failed to start the protocol client for decrypt engine connection")
		return err
	}
	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		err = stacktrace.Propagate(err, "failed to return the protocol client for decrypt engine connection")
		return err
	}
	// Request the plugin
	raw, err := rpcClient.Dispense("decryptor_grpc")
	if err != nil {
		err = stacktrace.Propagate(err, "RPC Client could not dispense a new instance of decryptor_grpc")
		return err
	}
	// We should have a encrypt store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	decryptor := raw.(shared.DecryptorInterface)
	resTmp, err := decryptor.Decrypt(req)
	if err != nil {
		err = stacktrace.Propagate(err, "decryptor failed to decrypt given input")
		return err
	}
	res.OutputHash = resTmp.OutputHash
	return nil
}
