package daemon

import (
	dare "github.com/da-moon/coe817-dare/dare"
	stacktrace "github.com/palantir/stacktrace"
	"io"
	"log"
	"os"
	"sync"
)

// Core ...
type Core struct {
	conf         *dare.CoreConfig
	logger       *log.Logger
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// Create ...
func Create(coreConf *Config, conf *dare.CoreConfig, logOutput io.Writer) (*Core, error) {

	if logOutput == nil {
		logOutput = os.Stderr
	}
	conf.LogOutput = logOutput
	conf.DevelopmentMode = coreConf.DevelopmentMode
	conf.MasterKey = coreConf.MasterKey
	conf.Protocol = uint8(coreConf.Protocol)
	conf.Init()
	core := &Core{
		conf:       conf,
		logger:     log.New(logOutput, "", log.LstdFlags),
		shutdownCh: make(chan struct{}),
	}
	return core, nil
}

// Start ...
func (a *Core) Start() error {
	a.logger.Printf("[INFO] dare daemon core: starting...")
	if len(a.conf.MasterKey) == 0 {
		return stacktrace.NewError("master key could not be found")
	}
	return nil
}

// Shutdown ...
func (a *Core) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()
	a.logger.Println("[INFO] dare daemon core: shutdown complete")
	a.shutdown = true
	// close(a.shutdownCh)
	return nil
}

// ShutdownCh ...
func (a *Core) ShutdownCh() <-chan struct{} {
	return a.shutdownCh
}
