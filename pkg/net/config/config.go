package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	APIAddr        string `yaml:"api_addr"`
	DataStreamAddr string `yaml:"data_stream_addr"`
	DB             struct {
		Type       string `yaml:"type"`
		ConnString string `yaml:"conn_string"`
	} `yaml:"db"`
	Log struct {
		DebugLog bool   `yaml:"debug_log"`
		FileName string `yaml:"file"`
	} `yaml:"log"`
	Debug struct {
		Enabled    bool   `yaml:"enabled"`
		Client1ID  string `yaml:"client1_id"`
		Client1Key string `yaml:"client1_key"`
		Client2ID  string `yaml:"client2_id"`
		Client2Key string `yaml:"client2_key"`
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
