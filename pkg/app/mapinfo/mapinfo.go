package mapinfo

import (
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
	"gopkg.in/yaml.v2"
)

type Wall struct {
	X1 int `yaml:"x1"`
	Y1 int `yaml:"y1"`
	X2 int `yaml:"x2"`
	Y2 int `yaml:"y2"`
}

type Event struct {
	No            int  `yaml:"no"`
	RequireAction bool `yaml:"require_action"`
	X             int  `yaml:"x"`
	Y             int  `yaml:"y"`
	R             int  `yaml:"r"`
}

type MapInfo struct {
	ID               int     `yaml:"id"`
	Name             string  `yaml:"name"`
	CollisionWalls   []Wall  `yaml:"walls"`
	Events           []Event `yaml:"events"`
	IsEnemyEncounter bool    `yaml:"is_enemy_encounter"`
	IsCyberWorld     bool    `yaml:"is_cyber_world"`

	Image int
	Size  point.Point
}

const (
	// 順番をmapInfo.yamlと合わせる

	ID_犬小屋 int = iota
	ID_秋原町

	idMax
)

var (
	mapInfo []MapInfo
)

func Init(fname string) error {
	buf, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(buf, &mapInfo); err != nil {
		return err
	}
	for i := 0; i < len(mapInfo); i++ {
		mapInfo[i].Image = -1
	}

	return nil
}

func Load(id int) (MapInfo, error) {
	if id < 0 || id >= idMax {
		return MapInfo{}, errors.Newf("no such as map %d", id)
	}

	if mapInfo[id].Image == -1 {
		fname := fmt.Sprintf("%smap/field/%d_%s.png", config.ImagePath, id, mapInfo[id].Name)
		mapInfo[id].Image = dxlib.LoadGraph(fname)
		if mapInfo[id].Image == -1 {
			return MapInfo{}, errors.Newf("failed to load image: %s", fname)
		}
		dxlib.GetGraphSize(mapInfo[id].Image, &mapInfo[id].Size.X, &mapInfo[id].Size.Y)
	}

	// TODO: 適切な背景をセットする
	if err := background.Set(background.Type秋原町); err != nil {
		return MapInfo{}, errors.Wrap(err, "failed to load background")
	}

	return mapInfo[id], nil
}
