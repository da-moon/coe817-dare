package daemon

import (
	model "github.com/da-moon/coe817-dare/model"
	encryptorGrpc "github.com/da-moon/coe817-dare/plugins/encryptor/grpc"
	encryptorRpc "github.com/da-moon/coe817-dare/plugins/encryptor/net-rpc"

	decryptorGrpc "github.com/da-moon/coe817-dare/plugins/decryptor/grpc"
	decryptorRpc "github.com/da-moon/coe817-dare/plugins/decryptor/net-rpc"
	shared "github.com/da-moon/coe817-dare/plugins/shared"
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
}

// Encrypt json rpc2.0 call that pparses variables and passes it to
// Encrypt Plugin. Sample query
// jq -n \
//   --arg source "/tmp/plain" \
//   --arg destination "/tmp/encrypted" \
//   --arg id "1" \
//   --arg method "Service.Encrypt" \
//  '{"jsonrpc": "2.0", "method":$method,"params":{"source": $source, "destination":$destination},"id": $id}' | curl \
//     -X POST  \
//  	--silent \
//     --header "Authorization: 12445" \
// 	--header "Content-type: application/json" \
//     --data @- \
//     http://127.0.0.1:8081/rpc  | jq -r
func (s *Service) Encrypt(r *http.Request, req *model.EncryptRequest, res *model.EncryptResponse) error {
	s.logger.Printf("[INFO] daemon-service: Encrypt Called")
	// We don't want to see the plugin logs.
	// log.SetOutput(ioutil.Discard)
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
	// Request the plugin
	raw, err := rpcClient.Dispense("encryptor_grpc")
	if err != nil {
		err = stacktrace.Propagate(err, "RPC Client could not dispense a new instance of encryptor_grpc")
		return err
	}
	// We should have a encrypt store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	encryptor := raw.(shared.EncryptorInterface)
	srcHash, dstHash, err := encryptor.Encrypt(req.Source, req.Destination)
	if err != nil {
		err = stacktrace.Propagate(err, "encryptor failed to encrypt given input")
		return err
	}
	res.DestinationHash = dstHash
	res.SourceHash = srcHash
	return nil
}

// Decrypt json rpc2.0 call that parses variables and passes it to
// Decrypt Plugin. Sample query
// jq -n \
//   --arg source "/tmp/encrypted" \
//   --arg destination "/tmp/decrypted" \
//   --arg id "2" \
//   --arg method "Service.Decrypt" \
//  '{"jsonrpc": "2.0", "method":$method,"params":{"source": $source, "destination":$destination},"id": $id}' | curl \
// 	-X POST  \
// 	--silent \
//     --header "Authorization: 12445" \
//     --header "Content-type: application/json" \
//     --data @- \
//     http://127.0.0.1:8081/rpc | jq -r
func (s *Service) Decrypt(r *http.Request, req *model.DecryptRequest, res *model.DecryptResponse) error {
	s.logger.Printf("[INFO] daemon-service: Decrypt Called")
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
	srcHash, dstHash, err := decryptor.Decrypt(req.Source, req.Destination)
	if err != nil {
		err = stacktrace.Propagate(err, "decryptor failed to decrypt given input")
		return err
	}
	res.DestinationHash = dstHash
	res.SourceHash = srcHash
	return nil
}
