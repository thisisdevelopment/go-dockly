package xhelper

import (
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/thisisdevelopment/go-dockly/v2/xlogger"
	"gopkg.in/yaml.v2"
)

// GetLogger returns the application default logger
func (h *Helper) GetLogger() *xlogger.Logger {
	l, err := xlogger.New(new(xlogger.Config))
	h.suite.Require().NoError(err)

	return l
}

// GetConfig returns the config.ServiceConfig
func (h *Helper) GetConfig(path string, cfg interface{}) {
	var err error
	b := h.BytesFromFile(path)

	switch true {
	case strings.Contains(path, "toml"):
		err = toml.Unmarshal(b, cfg)
	case strings.Contains(path, "yaml") || strings.Contains(path, "yml"):
		err = yaml.Unmarshal(b, cfg)
	}
	h.suite.Require().NoError(err)
}
