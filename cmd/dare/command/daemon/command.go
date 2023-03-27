package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	daemon "github.com/da-moon/dare-cli/daemon"
	view "github.com/da-moon/dare-cli/pkg/view"
	logutils "github.com/hashicorp/logutils"
	cli "github.com/mitchellh/cli"
)

const (
	gracefulTimeout = 3 * time.Second
)

// Command ...
type Command struct {
	Ui         cli.Ui
	ShutdownCh <-chan struct{}
	args       []string
	logFilter  *logutils.LevelFilter
	logger     *log.Logger
}

var _ cli.Command = &Command{}

// Run ...
func (c *Command) Run(args []string) int {
	c.Ui = &cli.PrefixedUi{
		OutputPrefix: "==> ",
		InfoPrefix:   "    ",
		ErrorPrefix:  "==> ",
		WarnPrefix:   "==> ",
		Ui:           c.Ui,
	}

	c.args = args
	config := c.readConfig()
	if config == nil {
		return 1
	}
	logGate, logWriter, logOutput := c.setupLoggers(config)
	if logWriter == nil {
		return 1
	}
	core := c.setupCore(config, logOutput)
	if core == nil {
		return 1
	}
	defer core.Shutdown()
	err := core.Start()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("[ERROR] Failed to start the dare daemon core: %v", err))
		return 1
	}
	apiEngine := c.startAPIEngine(config, core, logWriter, logOutput)
	if apiEngine == nil {
		return 1
	}

	c.Ui.Output("Log data will now stream in as it occurs:\n")
	logGate.Flush()
	return c.handleSignals(
		config,
		core,
	)
}

func (c *Command) handleSignals(config *Config, core *Core) int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

WAIT:
	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	case <-c.ShutdownCh:
		sig = os.Interrupt
	case <-core.ShutdownCh():
		return 0

	}
	c.Ui.Output(fmt.Sprintf("Caught signal: %v", sig))
	if sig == syscall.SIGHUP {
		config = c.handleReload(config, core)
		goto WAIT
	}

	graceful := false
	if sig == os.Interrupt {
		graceful = true
	} else if sig == syscall.SIGTERM {
		graceful = true
	}

	if !graceful {
		return 1
	}

	gracefulCh := make(chan struct{})
	c.Ui.Output("Gracefully shutting down core...")
	go func() {
		if err := core.Shutdown(); err != nil {
			c.Ui.Error(fmt.Sprintf("[ERROR]: %s", err.Error()))
			return
		}
		close(gracefulCh)
	}()

	select {
	case <-signalCh:
		return 1
	case <-time.After(gracefulTimeout):
		return 1
	case <-gracefulCh:
		return 0
	}
}

func (c *Command) handleReload(config *Config, core *Core) *Config {
	c.Ui.Output("Reloading configuration...")
	newConf := c.readConfig()
	if newConf == nil {
		c.Ui.Error(fmt.Sprintf("[ERROR] Failed to reload configs"))
		return config
	}

	minLevel := logutils.LogLevel(strings.ToUpper(newConf.LogLevel))
	if view.ValidateLevelFilter(minLevel, c.logFilter) {
		c.logFilter.SetMinLevel(minLevel)
	} else {
		c.Ui.Error(fmt.Sprintf(
			"[ERROR] Invalid log level: %s. Valid log levels are: %v",
			minLevel, c.logFilter.Levels))

		newConf.LogLevel = config.LogLevel
	}

	return newConf
}

// Synopsis ...
func (c *Command) Synopsis() string {
	return "data at rest encryption daemon"
}

// Help ...
func (c *Command) Help() string {
	helpText := `
Usage: dare daemon [options]

  Starts data at rest encryption daemon. it is a long running process
  that exposes an API endpoint at /rpc which intercepts user JSON messages,
  relays them to encrypt/decrypt plugins.

Options:

  -api-addr=127.0.0.1:8080 Address to bind the daemon json API listener.
  -api-password=secret     Daemon API password, used as Authorization 
                           header when sending JSON requests .
  -dev                     starts dare agent in development mode
  -config-file=foo         Path to a JSON file to read configuration from.
                           This can be specified multiple times.
  -config-dir=foo          Path to a directory to read configuration files
                           from. This will read every file ending in ".json"
                           as configuration in this directory in alphabetical
                           order.
  -log-level=info          Log level used in the encryption process.
  -encryptor-path=foo      Path encryptor plugin is located at.
  -decryptor-path=foo      Path decryptor plugin is located at.
`
	return strings.TrimSpace(helpText)
}

func (c *Command) setupCore(config *Config, logOutput io.Writer) *Core {
	coreConfig := daemon.DefaultCoreConfig()
	coreConfig.Protocol = uint8(config.Protocol)
	c.Ui.Output("Creating dare daemon core...")
	core, err := Create(config, coreConfig, logOutput)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("[ERROR] Failed to create the dare daemon core: %v", err))
		return nil
	}
	return core
}
func (c *Command) setupLoggers(config *Config) (*view.GatedWriter, *view.LogWriter, io.Writer) {
	logGate := &view.GatedWriter{
		Writer: &cli.UiWriter{Ui: c.Ui},
	}
	c.logFilter = view.LevelFilter()
	c.logFilter.MinLevel = logutils.LogLevel(strings.ToUpper(config.LogLevel))
	c.logFilter.Writer = logGate
	if !view.ValidateLevelFilter(c.logFilter.MinLevel, c.logFilter) {
		c.Ui.Error(fmt.Sprintf(
			"[ERROR] Invalid log level: %s. Valid log levels are: %v",
			c.logFilter.MinLevel, c.logFilter.Levels))
		return nil, nil, nil
	}
	LogWriter := view.NewLogWriter(512)
	var logOutput io.Writer
	logOutput = io.MultiWriter(c.logFilter, LogWriter)
	c.logger = log.New(logOutput, "", log.LstdFlags)
	return logGate, LogWriter, logOutput
}
