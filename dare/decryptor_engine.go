package dare

import (
	"context"
	"log"
	"os"
	"sync"

	stacktrace "github.com/palantir/stacktrace"
)

// DecryptorEngine ...
type DecryptorEngine struct {
	stateLock  sync.Mutex
	state      DecryptorEngineState
	logger     *log.Logger
	config     *CoreConfig
	shutdownCh chan struct{}
}

// CreateDecryptorEngine ...
func CreateDecryptorEngine(config *CoreConfig) (*DecryptorEngine, error) {
	config.Init()
	if config.Protocol < DecryptorEngineVersionMin {
		return nil, stacktrace.NewError("Decryptor Engine version '%d' too low. Must be in range: [%d, %d]",
			config.Protocol, DecryptorEngineVersionMin, DecryptorEngineVersionMin)
	} else if config.Protocol > DecryptorEngineVersionMax {
		return nil, stacktrace.NewError("Decryptor Engine version '%d' too high. Must be in range: [%d, %d]",
			config.Protocol, DecryptorEngineVersionMin, DecryptorEngineVersionMax)
	}

	logger := config.Logger
	if logger == nil {
		logOutput := config.LogOutput
		if logOutput == nil {
			logOutput = os.Stderr
		}
		logger = log.New(logOutput, "", log.LstdFlags)
	}

	result := &DecryptorEngine{
		config:     config,
		logger:     logger,
		shutdownCh: make(chan struct{}),
		state:      DecryptorEngineRunning,
	}
	// Start background tasks for this engine here ....
	// go result.handleReap()
	return result, nil
}

// Version ...
func (s *DecryptorEngine) Version() uint8 {
	return s.config.Protocol
}

// Shutdown ...
func (s *DecryptorEngine) Shutdown(ctx context.Context) error {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()

	if s.state == DecryptorEngineShutdown {
		return nil
	}
	return nil
}

// ShutdownCh ...
func (s *DecryptorEngine) ShutdownCh() <-chan struct{} {
	return s.shutdownCh
}
