package gamehandler

import (
	"bytes"
	"encoding/gob"
	"sort"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/action"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
)

type PanelInfo struct {
	OwnerClientID string
}

type GameHandler struct {
	Panels  [config.FieldNumX][config.FieldNumY]PanelInfo
	Objects map[string]*object.Object
}

func NewHandler() *GameHandler {
	return &GameHandler{
		Objects: make(map[string]*object.Object),
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
	g.Objects[param.ID] = &object.Object{
		ID:            param.ID,
		OwnerClientID: clientID,
		HP:            param.HP,
		Pos:           common.Point{X: param.X, Y: param.Y},
	}
}

func (g *GameHandler) MoveObject(moveInfo action.Move) {
	obj, ok := g.Objects[moveInfo.ObjectID]
	if !ok {
		logger.Info("Failed to find move target object: %+v", moveInfo)
		return
	}

	switch moveInfo.Type {
	case action.MoveTypeDirect:
		battlecommon.MoveObject(&obj.Pos, moveInfo.Direct, g.getPlayerPanelType(obj.OwnerClientID), true, g.getPanelInfo)
	case action.MoveTypeAbs:
		target := common.Point{X: moveInfo.AbsPosX, Y: moveInfo.AbsPosY}
		battlecommon.MoveObjectDirect(&obj.Pos, target, g.getPlayerPanelType(obj.OwnerClientID), true, g.getPanelInfo)
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

func (g *GameHandler) getPlayerPanelType(clientID string) int {
	var ids []string
	for objID := range g.Objects {
		ids = append(ids, objID)
	}
	sort.Strings(ids)
	for i, id := range ids {
		if clientID == id {
			return i
		}
	}
	return 0
}

func (g *GameHandler) getPanelInfo(pos common.Point) battlecommon.PanelInfo {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return battlecommon.PanelInfo{}
	}

	pn := g.Panels[pos.X][pos.Y]
	objectID := ""
	for _, obj := range g.Objects {
		if obj.Pos.X == pos.X && obj.Pos.Y == pos.Y {
			objectID = obj.ID
			break
		}
	}

	return battlecommon.PanelInfo{
		Type:     g.getPlayerPanelType(pn.OwnerClientID),
		ObjectID: objectID,

		// TODO: 適切な値を入れる
		Status:    battlecommon.PanelStatusNormal,
		HoleCount: 0,
	}
}
