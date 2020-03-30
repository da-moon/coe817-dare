package dare

import (
	"context"
	"log"
	"os"
	"sync"

	stacktrace "github.com/palantir/stacktrace"
)

// EncryptorEngine ...
type EncryptorEngine struct {
	stateLock  sync.Mutex
	state      EncryptorEngineState
	logger     *log.Logger
	config     *CoreConfig
	shutdownCh chan struct{}
}

// CreateEncryptorEngine ...
func CreateEncryptorEngine(config *CoreConfig) (*EncryptorEngine, error) {
	config.Init()
	if config.Protocol < EncryptorEngineVersionMin {
		return nil, stacktrace.NewError("Encryptor Engine version '%d' too low. Must be in range: [%d, %d]",
			config.Protocol, EncryptorEngineVersionMin, EncryptorEngineVersionMin)
	} else if config.Protocol > EncryptorEngineVersionMax {
		return nil, stacktrace.NewError("Encryptor Engine version '%d' too high. Must be in range: [%d, %d]",
			config.Protocol, EncryptorEngineVersionMin, EncryptorEngineVersionMax)
	}

	logger := config.Logger
	if logger == nil {
		logOutput := config.LogOutput
		if logOutput == nil {
			logOutput = os.Stderr
		}
		logger = log.New(logOutput, "", log.LstdFlags)
	}

	result := &EncryptorEngine{
		config:     config,
		logger:     logger,
		shutdownCh: make(chan struct{}),
		state:      EncryptorEngineRunning,
	}
	// Start background tasks for this engine here ....
	// go result.handleReap()
	return result, nil
}

// Version ...
func (s *EncryptorEngine) Version() uint8 {
	return s.config.Protocol
}

// Shutdown ...
func (s *EncryptorEngine) Shutdown(ctx context.Context) error {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()

	if s.state == EncryptorEngineShutdown {
		return nil
	}
	return nil
}

// ShutdownCh ...
func (s *EncryptorEngine) ShutdownCh() <-chan struct{} {
	return s.shutdownCh
}

// Run ...
func (s *EncryptorEngine) Run() error {
	return nil
}
