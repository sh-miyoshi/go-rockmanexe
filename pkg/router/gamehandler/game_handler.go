package gamehandler

import (
	"bytes"
	"encoding/gob"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/action"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
)

type PanelInfo struct {
	OwnerClientID string
	ShowHitArea   bool
}

type GameHandler struct {
	Panels  [config.FieldNumX][config.FieldNumY]PanelInfo
	Objects map[string]object.Object
}

func NewHandler() *GameHandler {
	return &GameHandler{
		Objects: make(map[string]object.Object),
	}
}

func (g *GameHandler) Init(clientIDs [2]string) error {
	for y := 0; y < config.FieldNumY; y++ {
		hx := config.FieldNumX / 2
		for x := 0; x < hx; x++ {
			g.Panels[x][y] = PanelInfo{OwnerClientID: clientIDs[0]}
			g.Panels[x+hx][y] = PanelInfo{OwnerClientID: clientIDs[1]}
		}
	}
	return nil
}

func (g *GameHandler) AddObject(clientID string, param object.InitParam) {
	g.Objects[param.ID] = object.Object{
		ID:            param.ID,
		OwnerClientID: clientID,
		HP:            param.HP,
		Pos:           common.Point{X: param.X, Y: param.Y},
	}
}

func (g *GameHandler) MoveObject(moveInfo action.Move) {
	switch moveInfo.Type {
	case action.MoveTypeDirect:
	case action.MoveTypeAbs:

	}
}

func (g *GameHandler) GetInfo() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(g)
	return buf.Bytes()
}

func (g *GameHandler) ParseInfo(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(g)
}
