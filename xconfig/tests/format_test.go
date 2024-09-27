package tests

import (
	"github.com/thisisdevelopment/go-dockly/v2/xconfig"
	"testing"
)

type JsonTest struct {
	Test string `json:"test"`
}

type YamlTest struct {
	Test string `yaml:"test"`
}

type TomlTest struct {
	Test string `toml:"test"`
}

func TestLoadJson(t *testing.T) {
	var cfg JsonTest

	err := xconfig.LoadConfig("cfg.json", &cfg)
	if err != nil {
		t.Error(err)
	}

	t.Log(cfg.Test)
}

func TestLoadYaml(t *testing.T) {
	var cfg YamlTest

	err := xconfig.LoadConfig("cfg.yaml", &cfg)
	if err != nil {
		t.Error(err)
	}

	err = xconfig.LoadConfig("cfg.yml", &cfg)
	if err != nil {
		t.Error(err)
	}

	t.Log(cfg.Test)
}

func TestLoadToml(t *testing.T) {
	var cfg TomlTest

	err := xconfig.LoadConfig("cfg.toml", &cfg)
	if err != nil {
		t.Error(err)
	}

	t.Log(cfg.Test)
}
