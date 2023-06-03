package config

import (
	"fmt"
	"os"

	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	APIAddr        string `yaml:"api_addr"`
	DataStreamAddr string `yaml:"data_stream_addr"`
	Log            struct {
		DebugLog bool   `yaml:"debug_log"`
		FileName string `yaml:"file"`
	} `yaml:"log"`
	AcceptableVersion string `yaml:"acceptable_version"`
	ChipFilePath      string `yaml:"chip_file_path"`
	Debug             struct {
		InvincibleCount *int `yaml:"invincible_count"`
	}
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

	if inst.Debug.InvincibleCount != nil {
		battlecommon.PlayerDefaultInvincibleTime = *inst.Debug.InvincibleCount
	}

	return nil
}

func Get() *Config {
	return &inst
}
