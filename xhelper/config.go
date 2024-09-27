package xhelper

import (
	"github.com/thisisdevelopment/go-dockly/v2/xconfig"
	"github.com/thisisdevelopment/go-dockly/v2/xlogger"
)

// GetLogger returns the application default logger
func (h *Helper) GetLogger() *xlogger.Logger {
	l, err := xlogger.New(new(xlogger.Config))
	h.suite.Require().NoError(err)

	return l
}

// GetConfig returns the config.ServiceConfig
func (h *Helper) GetConfig(path string, cfg interface{}) {
	h.suite.Require().NoError(xconfig.LoadConfig(path, cfg))
}
