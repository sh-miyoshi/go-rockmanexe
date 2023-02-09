package gameinfo

import (
	"bytes"
	"encoding/gob"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
)

type PanelInfo struct {
	OwnerClientID string
	ShowHitArea   bool
}

type GameInfo struct {
	Panels  [config.FieldNumX][config.FieldNumY]PanelInfo
	Objects map[string]object.Object
}

func (g *GameInfo) Init(clientIDs [2]string) {
	for y := 0; y < config.FieldNumY; y++ {
		hx := config.FieldNumX / 2
		for x := 0; x < hx; x++ {
			g.Panels[x][y] = PanelInfo{OwnerClientID: clientIDs[0]}
			g.Panels[x+hx][y] = PanelInfo{OwnerClientID: clientIDs[1]}
		}
	}
	g.Objects = make(map[string]object.Object)
}

func (g *GameInfo) AddObject(param object.InitParam) {
	id := uuid.New().String()
	g.Objects[id] = object.Object{
		ID:  id,
		HP:  param.HP,
		Pos: common.Point{X: param.X, Y: param.Y},
	}
}

func (g *GameInfo) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(g)
	return buf.Bytes()
}

func (g *GameInfo) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(g)
}
