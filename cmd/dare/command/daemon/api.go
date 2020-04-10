package daemon

import (
	"context"
	response "github.com/da-moon/coe817-dare/pkg/http/response"
	router "github.com/da-moon/coe817-dare/pkg/http/router"
	view "github.com/da-moon/coe817-dare/pkg/view"
	gorillaHandlers "github.com/gorilla/handlers"
	hclog "github.com/hashicorp/go-hclog"
	"strings"

	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// API ...
type API struct {
	sync.Mutex
	sync.Once
	core      *Core
	config    *Config
	listener  net.Listener
	logger    *log.Logger
	logWriter *view.LogWriter
	stop      bool
	stopCh    chan struct{}
	// ....
	api *http.Server
}

// NewAPIEngine ...
func NewAPIEngine(
	config *Config,
	core *Core,
	listener net.Listener,
	logOutput io.Writer,
	logWriter *view.LogWriter,
) *API {
	if logOutput == nil {
		logOutput = os.Stderr
	}

	backend := &API{
		config:    config,
		core:      core,
		listener:  listener,
		logger:    log.New(logOutput, "", log.LstdFlags),
		logWriter: logWriter,
		stopCh:    make(chan struct{}),
	}
	backend.Do(func() {
		go backend.startAPI()
	})
	return backend
}

// Shutdown ...
func (a *API) Shutdown() {
	a.Lock()
	defer a.Unlock()

	if a.stop {
		return
	}
	a.stop = true
	close(a.stopCh)
	// @TODO fix this
	a.logger.Printf("[INFO] api: gracefully shutting down api")
	a.api.Shutdown(context.Background())
	a.logger.Printf("[INFO] api: gracefully closing down api listener")
	a.listener.Close()
}

// add tls config
func (a *API) startAPI() {
	service := new(Service)
	service.logger = a.logger
	service.pluginLogger = hclog.New(&hclog.LoggerOptions{
		Level:  strToHCLogLevel(a.config.LogLevel),
		Output: a.logger.Writer(),
	})
	service.encryptor = a.core.conf.EncryptorPath
	service.decryptor = a.core.conf.DecryptorPath
	service.dev = a.core.conf.DevelopmentMode
	baseRouter := router.GenerateRPC2Routes([]router.JSON2{
		{
			Namespace: "",
			Endpoint:  "/rpc",
			Handler:   service,
		},
	})

	baseRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.WriteErrorJSON(&w, r, http.StatusMethodNotAllowed, "The specified method is not allowed against this resource")
	})
	recoveredRouter := gorillaHandlers.RecoveryHandler()(baseRouter)
	apiRouter := gorillaHandlers.LoggingHandler(a.logger.Writer(), recoveredRouter)
	a.api = &http.Server{
		Addr:        a.listener.Addr().String(),
		Handler:     apiRouter,
		IdleTimeout: 90 * time.Second,
	}
	err := a.api.Serve(a.listener)
	if err != nil && err != http.ErrServerClosed {
		a.logger.Printf("[ERR] agent.api: start failed: %v", err)
	}

}
func strToHCLogLevel(input string) hclog.Level {
	input = strings.ToUpper(input)
	var result hclog.Level
	switch input {
	case "TRACE":
		{
			result = hclog.Trace
			break
		}
	case "DEBUG":
		{
			result = hclog.Debug
			break

		}
	case "INFO":
		{
			result = hclog.Info
			break

		}
	case "WARN":
		{
			result = hclog.Warn
			break

		}
	case "ERROR":
		{
			result = hclog.Error
			break

		}
	default:
		{
			result = hclog.NoLevel
			break
		}
	}
	return result
}
