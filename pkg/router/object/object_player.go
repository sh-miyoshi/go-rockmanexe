package object

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/action"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
)

type Player struct {
	objectInfo gameinfo.Object
	gameInfo   *gameinfo.GameInfo
}

func NewPlayer(info gameinfo.Object, gameInfo *gameinfo.GameInfo) *Player {
	return &Player{
		objectInfo: info,
		gameInfo:   gameInfo,
	}
}

func (p *Player) Process() (bool, error) {
	// TODO
	return false, nil
}

func (p *Player) Draw() {
	// nothing to do at router
}

func (p *Player) DamageProc(dm *damage.Damage) bool {
	// TODO
	return false
}

func (p *Player) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    p.objectInfo.ID,
			Pos:      p.objectInfo.Pos,
			AnimType: anim.AnimTypeObject,
		},
		HP: p.objectInfo.HP,
	}
}

func (p *Player) GetObjectType() int {
	if p.gameInfo.ReverseClientID == p.objectInfo.OwnerClientID {
		return objanim.ObjTypeEnemy
	}
	return objanim.ObjTypePlayer
}

func (p *Player) MakeInvisible(count int) {
	// TODO
}

func (p *Player) AddMove(moveInfo action.Move) {
	// TODO: このタイミングで移動せず、アニメーションの追加のみにする
	switch moveInfo.Type {
	case action.MoveTypeDirect:
		battlecommon.MoveObject(&p.objectInfo.Pos, moveInfo.Direct, p.getPanelType(p.objectInfo.ID), true, p.getPanelInfo)
	case action.MoveTypeAbs:
		target := common.Point{X: moveInfo.AbsPosX, Y: moveInfo.AbsPosY}
		battlecommon.MoveObjectDirect(&p.objectInfo.Pos, target, p.getPanelType(p.objectInfo.ID), true, p.getPanelInfo)
	}
}

func (p *Player) getPanelInfo(pos common.Point) battlecommon.PanelInfo {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return battlecommon.PanelInfo{}
	}

	pn := p.gameInfo.Panels[pos.X][pos.Y]
	return battlecommon.PanelInfo{
		Type:     p.getPanelType(pn.OwnerClientID),
		ObjectID: pn.ObjectID,

		// TODO: 適切な値を入れる
		Status:    battlecommon.PanelStatusNormal,
		HoleCount: 0,
	}
}

func (p *Player) getPanelType(clientID string) int {
	if p.gameInfo.ReverseClientID == clientID {
		return 1
	}
	return 0
}
