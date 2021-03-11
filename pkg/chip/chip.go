package chip

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Chip ...
type Chip struct {
	ID        int    `yaml:"id"`
	Name      string `yaml:"name"`
	Power     uint   `yaml:"power"`
	Type      int    `yaml:"type"`
	Code      string `yaml:"code"`
	PlayerAct int    `yaml:"player_act"`

	// TODO Image, ImgIcon
}

const (
	// IDCannon ...
	IDCannon = iota

	idMax
)

var (
	chipData []Chip
)

// Init ...
func Init(fname string) error {
	// Load chip data
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(buf, &chipData); err != nil {
		return err
	}

	// TODO: set image

	return nil
}

// Get ...
func Get(id int) Chip {
	return chipData[id]
}
