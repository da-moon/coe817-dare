package dare

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	dareCore "github.com/da-moon/coe817-dare"
	hkdf "golang.org/x/crypto/hkdf"
	"io"
	"log"
	"os"
	"sync"

	stacktrace "github.com/palantir/stacktrace"
)

// DareEngine ...
type DareEngine struct {
	stateLock sync.Mutex
	state     DareEngineState
	logger    *log.Logger
	config    *CoreConfig
	key       [32]byte

	shutdownCh chan struct{}
}

// CreateDareEngine ...
func CreateDareEngine(config *CoreConfig) (*DareEngine, error) {
	config.Init()
	if config.Protocol < DareEngineVersionMin {
		return nil, stacktrace.NewError("Data At Rest Encryption version '%d' too low. Must be in range: [%d, %d]",
			config.Protocol, DareEngineVersionMin, DareEngineVersionMin)
	} else if config.Protocol > DareEngineVersionMax {
		return nil, stacktrace.NewError("Data At Rest Encryption version '%d' too high. Must be in range: [%d, %d]",
			config.Protocol, DareEngineVersionMin, DareEngineVersionMax)
	}

	logger := config.Logger
	if logger == nil {
		logOutput := config.LogOutput
		if logOutput == nil {
			logOutput = os.Stderr
		}
		logger = log.New(logOutput, "", log.LstdFlags)
	}
	var (
		nonce [32]byte
		key   [32]byte
	)

	// setting up enc key
	masterkey, err := hex.DecodeString(config.MasterKey)
	if err != nil {
		err = stacktrace.Propagate(err, "could not initialize encryptor engine due to failure in decoding masterkey")
		return nil, err
	}
	_, err = io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		err = stacktrace.Propagate(err, "could not initialize encryptor engine due to failure in generating random data for nonce")
		return nil, err
	}
	// driving master key ...
	kdf := hkdf.New(sha256.New, masterkey, nonce[:], nil)
	_, err = io.ReadFull(kdf, key[:])
	if err != nil {
		err = stacktrace.Propagate(err, "initialize encryptor engine due to failure in driving an encryption key. masterkey=%v nonce=%v", masterkey, nonce[:])
		return nil, err
	}
	result := &DareEngine{
		config:     config,
		logger:     logger,
		shutdownCh: make(chan struct{}),
		key:        key,
		state:      DareEngineRunning,
	}

	// Start background tasks for this engine here ....
	// go result.handleReap()
	return result, nil
}

// Version ...
func (s *DareEngine) Version() uint8 {
	return s.config.Protocol
}

// Shutdown ...
func (s *DareEngine) Shutdown(ctx context.Context) error {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()

	if s.state == DareEngineShutdown {
		return nil
	}
	return nil
}

// ShutdownCh ...
func (s *DareEngine) ShutdownCh() <-chan struct{} {
	return s.shutdownCh
}

// Run ...
func (s *DareEngine) Encrypt(source string, destination string) error {

	// reading file at source
	fi, err := os.Stat(source)
	if err == nil {
		if fi.Size() == 0 {
			s.logger.Printf("[WARN] encrypt : Target source file (%s) exists but is size zero, it may be left from some previous FS error like out-of-space ", source)
			os.Remove(source)
			return nil
		}
	}
	srcFile, err := os.Open(source)
	if srcFile != nil {
		defer srcFile.Close()
	}
	if err != nil {
		err = stacktrace.Propagate(err, "could not encrypt due to failure in opening source file at %s", source)
		return err
	}
	dstFile, err := os.OpenFile(
		destination,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0600)

	if dstFile == nil {
		err = stacktrace.NewError("could not successfully get a file handle for %s", destination)
		return err
	}

	if dstFile != nil {
		defer dstFile.Close()
	}
	if err != nil {
		err = stacktrace.Propagate(err, "Could not create empty file at (%s) ", destination)
		return err
	}
	_, err = dareCore.Encrypt(
		dstFile,
		srcFile,
		s.key[:],
	)
	return nil
}

// Run ...
func (s *DareEngine) Decrypt(source string, destination string) error {

	// reading file at source
	fi, err := os.Stat(source)
	if err == nil {
		if fi.Size() == 0 {
			s.logger.Printf("[WARN] decrypt : Target source file (%s) exists but is size zero, it may be left from some previous FS error like out-of-space ", source)
			os.Remove(source)
			return nil
		}
	}
	srcFile, err := os.Open(source)
	if srcFile != nil {
		defer srcFile.Close()
	}
	if err != nil {
		err = stacktrace.Propagate(err, "could not decrypt due to failure in opening source file at %s", source)
		return err
	}
	dstFile, err := os.OpenFile(
		destination,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0600)

	if dstFile == nil {
		err = stacktrace.NewError("could not successfully get a file handle for %s", destination)
		return err
	}

	if dstFile != nil {
		defer dstFile.Close()
	}
	if err != nil {
		err = stacktrace.Propagate(err, "Could not create empty file at (%s) ", destination)
		return err
	}
	_, err = dareCore.Decrypt(
		dstFile,
		srcFile,
		s.key[:],
	)
	return nil
}
