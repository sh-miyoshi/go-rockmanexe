package gamehandler

import (
	"sort"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/action"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
)

type PanelInfo struct {
	OwnerClientID string
}

type GameHandler struct {
	info gameinfo.GameInfo
}

func NewHandler() *GameHandler {
	return &GameHandler{
		info: gameinfo.GameInfo{
			Objects: make(map[string]object.Object),
		},
	}
}

func (g *GameHandler) Init(clientIDs [2]string) error {
	for y := 0; y < config.FieldNumY; y++ {
		hx := config.FieldNumX / 2
		for x := 0; x < hx; x++ {
			g.info.Panels[x][y] = gameinfo.PanelInfo{OwnerClientID: clientIDs[0]}
			g.info.Panels[x+hx][y] = gameinfo.PanelInfo{OwnerClientID: clientIDs[1]}
		}
	}
	g.info.ReverseClientID = clientIDs[1]
	return nil
}

func (g *GameHandler) AddObject(clientID string, param object.InitParam) {
	x := param.X
	if g.info.ReverseClientID == clientID {
		x = battlecommon.FieldNum.X - x - 1
	}

	g.info.Objects[param.ID] = object.Object{
		ID:            param.ID,
		OwnerClientID: clientID,
		HP:            param.HP,
		Pos:           common.Point{X: x, Y: param.Y},
	}
}

func (g *GameHandler) MoveObject(moveInfo action.Move) {
	obj, ok := g.info.Objects[moveInfo.ObjectID]
	if !ok {
		logger.Info("Failed to find move target object: %+v", moveInfo)
		return
	}

	// TODO: このタイミングで移動せず、アニメーションの追加のみにする
	switch moveInfo.Type {
	case action.MoveTypeDirect:
		battlecommon.MoveObject(&obj.Pos, moveInfo.Direct, g.getPlayerPanelType(obj.OwnerClientID), true, g.getPanelInfo)
	case action.MoveTypeAbs:
		target := common.Point{X: moveInfo.AbsPosX, Y: moveInfo.AbsPosY}
		battlecommon.MoveObjectDirect(&obj.Pos, target, g.getPlayerPanelType(obj.OwnerClientID), true, g.getPanelInfo)
	}

	g.info.Objects[moveInfo.ObjectID] = obj
}

func (g *GameHandler) GetInfo() []byte {
	return g.info.Marshal()
}

func (g *GameHandler) getPlayerPanelType(clientID string) int {
	var ids []string
	for objID := range g.info.Objects {
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

	pn := g.info.Panels[pos.X][pos.Y]
	objectID := ""
	for _, obj := range g.info.Objects {
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
