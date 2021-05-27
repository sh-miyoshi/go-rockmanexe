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
		Enabled           bool `yaml:"enabled"`
		SkipTitle         bool `yaml:"skip_title"`
		SkipMenu          bool `yaml:"skip_menu"`
		SkipBattleOpening bool `yaml:"skip_battle_opening"`
		StartContinue     bool `yaml:"start_continue"`
	} `yaml:"debug"`
	BGM struct {
		Disabled bool `yaml:"disabled"`
	} `yaml:"bgm"`
	Net struct {
		ClientID   string `yaml:"client_id"`
		ClientKey  string `yaml:"client_key"`
		StreamAddr string `yaml:"addr"`
	} `yaml:"net"`
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
