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

	err := xconfig.LoadConfig(&cfg, "cfg.json")
	if err != nil {
		t.Error(err)
	}

	t.Log(cfg.Test)
}

func TestLoadYaml(t *testing.T) {
	var cfg YamlTest

	err := xconfig.LoadConfig(&cfg, "cfg.yaml")
	if err != nil {
		t.Error(err)
	}

	err = xconfig.LoadConfig(&cfg, "cfg.yml")
	if err != nil {
		t.Error(err)
	}

	t.Log(cfg.Test)
}

func TestLoadToml(t *testing.T) {
	var cfg TomlTest

	err := xconfig.LoadConfig(&cfg, "cfg.toml")
	if err != nil {
		t.Error(err)
	}

	t.Log(cfg.Test)
}
