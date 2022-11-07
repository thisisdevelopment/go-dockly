package xconfig

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// LoadConfig reads in a toml file and inits the ServiceConfig
func LoadConfig(cfg interface{}, path string) error {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return errors.Wrapf(err, "unable to read file %s", path)
	}

	switch true {
	case strings.Contains(path, "toml"):
		err = toml.Unmarshal(bytes, cfg)
	case strings.Contains(path, "yaml") || strings.Contains(path, "yml"):
		err = yaml.Unmarshal(bytes, cfg)
	}

	if err != nil {
		return errors.Wrapf(err, "error while parsing config file %s", string(bytes))
	}

	return nil
}
