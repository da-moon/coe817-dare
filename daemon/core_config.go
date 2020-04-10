package daemon

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"crypto/rand"
	"encoding/base64"
	"sync"
)

const jwtSecretLen = 32

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
	EncryptorPath   string
	DecryptorPath   string
	APIAddr         string
	APIPassword     string
	ShutdownCh      chan struct{}
	DevelopmentMode bool
}

// Init ...
func (c *CoreConfig) Init() {
	c.Do(func() {
		logOutput := c.LogOutput
		if logOutput == nil {
			logOutput = os.Stderr
		}
		c.Logger = log.New(logOutput, "", log.LstdFlags)
		if len(c.EncryptorPath) == 0 {
			c.EncryptorPath = os.Getenv("DARE_ENCRYPTOR_PLUGIN")
			if len(c.EncryptorPath) == 0 {
				path, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				c.EncryptorPath = filepath.Join(path, "bin/encryptor-plugin")
			}
		}
		if len(c.DecryptorPath) == 0 {
			c.DecryptorPath = os.Getenv("DARE_DECRYPTOR_PLUGIN")
			if len(c.DecryptorPath) == 0 {
				path, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				c.DecryptorPath = filepath.Join(path, "bin/decryptor-plugin")
			}
		}
		if len(c.APIAddr) == 0 {
			c.APIAddr = os.Getenv("DARE_API_ADDRESS")
			if len(c.APIAddr) == 0 {
				c.APIAddr = "127.0.0.1:8080"
			}
		}
		if len(c.APIPassword) == 0 {
			c.APIPassword = os.Getenv("DARE_API_PASSWORD")
			if len(c.APIPassword) == 0 {
				bytes := make([]byte, jwtSecretLen)
				if _, err := rand.Read(bytes); err != nil {
					panic(err)
				}
				c.APIPassword = base64.StdEncoding.EncodeToString(bytes)
			}
		}
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
