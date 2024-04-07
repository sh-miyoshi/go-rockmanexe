package config

import (
	"fmt"
	"os"

	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Enabled bool `yaml:"enabled"`
		Session struct {
			ID         string `yaml:"id"`
			ClientID1  string `yaml:"client_1_id"`
			ClientKey1 string `yaml:"client_1_key"`
			ClientID2  string `yaml:"client_2_id"`
			ClientKey2 string `yaml:"client_2_key"`
		} `yaml:"session"`
	} `yaml:"server"`
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

	setByEnv()

	if inst.Debug.InvincibleCount != nil {
		battlecommon.PlayerDefaultInvincibleTime = *inst.Debug.InvincibleCount
	}

	return nil
}

func Get() *Config {
	return &inst
}

func APIAddr() string {
	if inst.Server.Enabled {
		return "http://localhost:3000"
	}
	return inst.APIAddr
}

func setByEnv() {
	if id := os.Getenv("CLIENT_1_ID"); id != "" {
		inst.Server.Session.ClientID1 = id
	}
	if key := os.Getenv("CLIENT_1_KEY"); key != "" {
		inst.Server.Session.ClientKey1 = key
	}
	if id := os.Getenv("CLIENT_2_ID"); id != "" {
		inst.Server.Session.ClientID2 = id
	}
	if key := os.Getenv("CLIENT_2_KEY"); key != "" {
		inst.Server.Session.ClientKey2 = key
	}
	if addr := os.Getenv("DATA_ADDR"); addr != "" {
		inst.DataStreamAddr = addr
	}
}
