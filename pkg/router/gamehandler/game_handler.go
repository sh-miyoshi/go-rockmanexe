package gamehandler

import (
	"sort"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/action"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/skill"
)

type GameHandler struct {
	info gameinfo.GameInfo
}

func NewHandler() *GameHandler {
	return &GameHandler{
		info: gameinfo.GameInfo{
			Objects: make(map[string]gameinfo.Object),
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

	g.info.Objects[param.ID] = gameinfo.Object{
		ID:            param.ID,
		OwnerClientID: clientID,
		HP:            param.HP,
		Pos:           common.Point{X: x, Y: param.Y},
	}
	g.info.Panels[x][param.Y].ObjectID = param.ID
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
	g.updatePanelObject()
}

func (g *GameHandler) AddBuster(clientID string, busterInfo action.Buster) {
	obj, ok := g.info.Objects[busterInfo.ObjectID]
	if !ok {
		logger.Info("Failed to find buster send object: %+v", busterInfo)
		return
	}

	// TODO: このタイミングで実行せず、アニメーションの追加のみにする
	// TODO: Busterアニメーション
	damageToObj := func(objID string, power int) {
		o := g.info.Objects[objID]
		o.HP -= power
		if o.HP < 0 {
			o.HP = 0
		}
		g.info.Objects[objID] = o
	}

	if g.info.ReverseClientID == clientID {
		for x := obj.Pos.X - 1; x >= 0; x-- {
			if objID := g.info.Panels[x][obj.Pos.Y].ObjectID; objID != "" {
				damageToObj(objID, busterInfo.Power)
			}
		}
	} else {
		for x := obj.Pos.X + 1; x < battlecommon.FieldNum.X; x++ {
			if objID := g.info.Panels[x][obj.Pos.Y].ObjectID; objID != "" {
				damageToObj(objID, busterInfo.Power)
			}
		}
	}
}

func (g *GameHandler) UseChip(clientID string, chipInfo action.UseChip) {
	// TODO
	s := skill.GetByChip(chipInfo.ChipID, skill.Argument{
		OwnerID:    clientID,
		Power:      40, // debug
		TargetType: 0,  // debug
	})
	anim.New(s)
}

func (g *GameHandler) GetInfo() []byte {
	return g.info.Marshal()
}

func (g *GameHandler) UpdateGameStatus() {
	if err := anim.MgrProcess(); err != nil {
		logger.Error("Failed to manage animation: %+v", err)
		// TODO: 処理を終了する
	}
	// TODO
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

	g.updatePanelObject()
	pn := g.info.Panels[pos.X][pos.Y]

	return battlecommon.PanelInfo{
		Type:     g.getPlayerPanelType(pn.OwnerClientID),
		ObjectID: pn.ObjectID,

		// TODO: 適切な値を入れる
		Status:    battlecommon.PanelStatusNormal,
		HoleCount: 0,
	}
}

func (g *GameHandler) updatePanelObject() {
	// Cleanup at first
	for y := 0; y < battlecommon.FieldNum.Y; y++ {
		for x := 0; x < battlecommon.FieldNum.X; x++ {
			g.info.Panels[x][y].ObjectID = ""
		}
	}
	for _, obj := range g.info.Objects {
		g.info.Panels[obj.Pos.X][obj.Pos.Y].ObjectID = obj.ID
	}
}
