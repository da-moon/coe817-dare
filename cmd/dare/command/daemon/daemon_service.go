package daemon

import (
	"log"
	"net/http"
)

// Service ...
type Service struct {
	logger *log.Logger
}

// EncryptRequest ...
type EncryptRequest struct {
	Source      string
	Destination string
}

// EncryptResponse ...
type EncryptResponse struct {
	MD5hash string
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
//     --header "Content-type: application/json" \
//     --data @- \
//     http://127.0.0.1:8081/rpc
func (s *Service) Encrypt(r *http.Request, req *EncryptRequest, res *EncryptResponse) error {
	s.logger.Printf("[INFO] daemon-service: Encrypt Called")
	return nil
}

// DecryptRequest ...
type DecryptRequest struct {
	Source      string
	Destination string
}

// DecryptResponse ...
type DecryptResponse struct {
	MD5hash string
}

// Decrypt json rpc2.0 call that parses variables and passes it to
// Decrypt Plugin. Sample query
// jq -n \
//   --arg source "/tmp/encrypted" \
//   --arg destination "/tmp/decrypted" \
//   --arg id "2" \
//   --arg method "Service.Decrypt" \
//  '{"jsonrpc": "2.0", "method":$method,"params":{"source": $source, "destination":$destination},"id": $id}' | curl \
//     -X POST  \
//     --header "Content-type: application/json" \
//     --data @- \
//     http://127.0.0.1:8081/rpc
func (s *Service) Decrypt(r *http.Request, req *DecryptRequest, res *DecryptResponse) error {
	s.logger.Printf("[INFO] daemon-service: Decrypt Called")
	return nil
}
