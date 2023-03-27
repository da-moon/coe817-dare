package encryptor

import (
	"golang.org/x/crypto/nacl/box"
	"io"
	"sync"

	config "github.com/da-moon/dare-cli/pkg/dare/config"
	"github.com/palantir/stacktrace"
)

// Writer ...
type Writer struct {
	stateLock sync.Mutex
	writer    io.Writer
	nonce     *[24]byte
	sharedKey *[32]byte
	chunkSize int
	buf       []byte
}

// NewWriter ...
func NewWriter(
	writer io.Writer,
	nonce [24]byte,
	sharedKey *[32]byte,
) *Writer {
	return &Writer{
		writer:    writer,
		nonce:     &nonce,
		sharedKey: sharedKey,
		chunkSize: config.DefaultChunkSize,
	}
}

// encrypt ..
func (e *Writer) Write(p []byte) (n int, err error) {
	e.encrypt(p)
	n, err = e.writer.Write(e.buf)
	if err != nil {
		err = stacktrace.Propagate(err, "encryptor could not finish writing due to failure of underlying io.reader")

		return 0, err
	}
	return n, nil
}

// encrypt ...
func (e *Writer) encrypt(p []byte) {
	e.stateLock.Lock()
	defer e.stateLock.Unlock()
	e.buf = box.SealAfterPrecomputation(nil, p, e.nonce, e.sharedKey)
	// copying first 24 bytes of output as current nonce for nonce chaining
	copy(e.nonce[:], e.buf[:24])
}
