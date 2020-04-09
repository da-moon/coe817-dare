package daemon

import (
	daemon "github.com/da-moon/coe817-dare/daemon"
	stacktrace "github.com/palantir/stacktrace"
	"io"
	"log"
	"os"
	"sync"
)

// Core ...
type Core struct {
	conf         *daemon.CoreConfig
	logger       *log.Logger
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// Create ...
func Create(coreConf *Config, conf *daemon.CoreConfig, logOutput io.Writer) (*Core, error) {

	if logOutput == nil {
		logOutput = os.Stderr
	}
	conf.LogOutput = logOutput
	conf.DevelopmentMode = coreConf.DevelopmentMode
	conf.Protocol = uint8(coreConf.Protocol)
	// todo remove this ... it may be very useless
	conf.EncryptorPath = coreConf.EncryptorPath
	conf.DecryptorPath = coreConf.DecryptorPath
	conf.APIAddr = coreConf.APIAddr
	conf.APIPassword = coreConf.APIPassword
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
