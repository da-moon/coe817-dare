package daemon

import (
	"fmt"
	view "github.com/da-moon/coe817-dare/pkg/view"
	stacktrace "github.com/palantir/stacktrace"
	"io"
	"net"
	"strconv"
	"time"
)

// AddrParts ...

func (c *Command) startAPIEngine(
	config *Config,
	core *Core,
	logWriter *view.LogWriter,
	logOutput io.Writer,
) *API {

	x, _ := TCPAddressFromString(config.APIAddr)
	l, err := net.Listen("tcp", x.String())
	if err != nil {
		l.Close()
		c.Ui.Error(fmt.Sprintf("Failed to start the daemon api listener: %v", err))
		return nil
	}
	apiListener := &tcpKeepAliveListener{l.(*net.TCPListener)}

	ipc := NewAPIEngine(
		core,
		apiListener,
		logOutput,
		logWriter)
	c.Ui.Output("Bifrost api running!")
	c.Ui.Info(fmt.Sprintf("                   API addr: '%s'", config.APIAddr))

	return ipc
}

// TCPAddress -
func TCPAddress(ip string, port int) *net.TCPAddr {
	result := &net.TCPAddr{IP: net.ParseIP(ip), Port: port}

	return result
}

// TCPAddressFromString -
func TCPAddressFromString(addr string) (*net.TCPAddr, error) {
	h, p, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, stacktrace.Propagate(err, "[Err] could not generate tcp address from string")
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		return nil, stacktrace.Propagate(err, "[Err] could not generate tcp address from string due to invalid port")

	}

	return TCPAddress(h, port), nil
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used so dead TCP connections eventually go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept ...
func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(30 * time.Second)
	return tc, nil
}
