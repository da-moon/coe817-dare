package dare

import (
	"io"
	"log"
	"os"
	"sync"
)

// CoreProtocolVersionMap ...
var CoreProtocolVersionMap map[uint8]uint8

func init() {
	CoreProtocolVersionMap = map[uint8]uint8{
		1: 1,
	}
}

const (
	// CoreVersionMin ...
	CoreVersionMin = 1
	// CoreVersionMax ...
	CoreVersionMax = 2
)

// CoreConfig ...
type CoreConfig struct {
	sync.Once
	Initialized     bool
	LogOutput       io.Writer
	Protocol        uint8
	Logger          *log.Logger
	MasterKey       string
	ShutdownCh      chan struct{}
	DevelopmentMode bool
}

// Init ...
func (c *CoreConfig) Init() {
	c.Do(func() {
		if len(c.MasterKey) == 0 {
			c.MasterKey = os.Getenv("DARE_MASTER_KEY")
			if len(c.MasterKey) == 0 {
				c.MasterKey = "b6c4bba7a385aef779965cb0b7d66316ab091704042606797871"
			}
		}
		logOutput := c.LogOutput
		if logOutput == nil {
			logOutput = os.Stderr
		}
		c.Logger = log.New(logOutput, "", log.LstdFlags)
		c.Initialized = true
	})
}

// DefaultCoreConfig ...
func DefaultCoreConfig() *CoreConfig {
	return &CoreConfig{
		Protocol:   1,
		LogOutput:  os.Stderr,
		ShutdownCh: make(chan struct{}),
	}
}
