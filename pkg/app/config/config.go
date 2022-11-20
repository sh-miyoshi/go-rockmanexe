package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Log struct {
		FileName     string `yaml:"file"`
		DebugEnabled bool   `yaml:"debug_enabled"`
	} `yaml:"log"`
	Debug struct {
		ShowDebugData      bool `yaml:"show_debug_data"`
		SkipTitle          bool `yaml:"skip_title"`
		SkipMenu           bool `yaml:"skip_menu"`
		SkipBattleOpening  bool `yaml:"skip_battle_opening"`
		StartContinue      bool `yaml:"start_continue"`
		InitSleepSec       int  `yaml:"init_sleep_sec"`
		RunAlways          bool `yaml:"run_always"`
		EnableDevFeature   bool `yaml:"enable_dev_feature"`
		UsePrivateResource bool `yaml:"use_private_res"`
		UseDebugFolder     bool `yaml:"use_debug_folder"`
	} `yaml:"debug"`
	BGM struct {
		Disabled bool `yaml:"disabled"`
	} `yaml:"bgm"`
	Net struct {
		Insecure   bool   `yaml:"insecure"`
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
