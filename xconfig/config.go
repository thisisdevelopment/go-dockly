package xconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// LoadConfig reads in a toml file and inits the ServiceConfig
func LoadConfig(filePath string, cfg any) error {
	bytes, err := os.ReadFile(filePath)

	if err != nil {
		return fmt.Errorf("unable to read file %s: %w", filePath, err)
	}

	switch path.Ext(filePath) {
	case ".toml":
		err = toml.Unmarshal(bytes, cfg)
	case ".yaml":
		fallthrough
	case ".yml":
		err = yaml.Unmarshal(bytes, cfg)
	case ".json":
		err = json.Unmarshal(bytes, cfg)
	}

	if err != nil {
		return fmt.Errorf("error while parsing config file %s: %w", string(bytes), err)
	}

	return nil
}

// MustConfig load config and panic if fails
func MustConfig(filePath string, cfg any) {
	err := LoadConfig(filePath, cfg)
	if err != nil {
		panic(err)
	}
}
