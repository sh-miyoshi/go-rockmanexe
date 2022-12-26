package mapinfo

import (
	"fmt"
	"io/ioutil"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"gopkg.in/yaml.v2"
)

type Wall struct {
	X1 int `yaml:"x1"`
	Y1 int `yaml:"y1"`
	X2 int `yaml:"x2"`
	Y2 int `yaml:"y2"`
}

type Event struct {
	No int `yaml:"no"`
	X  int `yaml:"x"`
	Y  int `yaml:"y"`
	R  int `yaml:"r"`
}

type MapInfo struct {
	ID             int     `yaml:"id"`
	Name           string  `yaml:"name"`
	CollisionWalls []Wall  `yaml:"walls"`
	Events         []Event `yaml:"events"`

	Image int
	Size  common.Point
}

const (
	// 順番をmapInfo.yamlと合わせる

	ID_犬小屋 int = iota

	idMax
)

var (
	mapInfo []MapInfo
)

func Init(fname string) error {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(buf, &mapInfo); err != nil {
		return err
	}

	return nil
}

func Load(id int) (*MapInfo, error) {
	if id < 0 || id >= idMax {
		return nil, fmt.Errorf("no such as map %d", id)
	}

	m := mapInfo[id]
	fname := fmt.Sprintf("%smap/field/%d_%s.png", common.ImagePath, m.ID, m.Name)
	res := &MapInfo{
		Image:          dxlib.LoadGraph(fname),
		CollisionWalls: m.CollisionWalls,
		Events:         m.Events,
	}
	if res.Image == -1 {
		return nil, fmt.Errorf("failed to load image: %s", fname)
	}
	dxlib.GetGraphSize(res.Image, &res.Size.X, &res.Size.Y)
	if err := background.Set(background.Type秋原町); err != nil {
		return nil, fmt.Errorf("failed to load background: %w", err)
	}

	return res, nil
}
