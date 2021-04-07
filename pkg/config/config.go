package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Log struct {
		FileName string `yaml:"file"`
	} `yaml:"log"`
	Debug struct {
		Enabled   bool `yaml:"enabled"`
		SkipTitle bool `yaml:"skip_title"`
	} `yaml:"debug"`
}

var (
	inst Config
)

func Init(fname string) error {
	fp, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer fp.Close()

	if err := yaml.NewDecoder(fp).Decode(&inst); err != nil {
		return fmt.Errorf("failed to decode yaml: %v", err)
	}

	return nil
}

func Get() *Config {
	return &inst
}
