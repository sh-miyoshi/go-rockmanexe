package object

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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
	if dm == nil {
		return false
	}

	// TODO: インビジブル関係

	myType := damage.TargetPlayer
	if p.isReverse() {
		myType = damage.TargetEnemy
	}

	if dm.TargetType&myType != 0 {
		p.objectInfo.HP -= dm.Power
		if p.objectInfo.HP < 0 {
			p.objectInfo.HP = 0
		}
		// TODO: HPMax
		// TODO: damage effect

		for i := 0; i < dm.PushLeft; i++ {
			if !battlecommon.MoveObject(&p.objectInfo.Pos, common.DirectLeft, p.gameInfo.GetPanelType(p.objectInfo.ID), true, p.gameInfo.GetPanelInfo) {
				break
			}
		}
		for i := 0; i < dm.PushRight; i++ {
			if !battlecommon.MoveObject(&p.objectInfo.Pos, common.DirectRight, p.gameInfo.GetPanelType(p.objectInfo.ID), true, p.gameInfo.GetPanelInfo) {
				break
			}
		}

		if dm.Power <= 0 {
			// Not damage, maybe recover or special anim
			return true
		}

		if !dm.BigDamage {
			return true
		}

		// TODO: sound

		// TODO: Stop current animation
		logger.Debug("Player damaged: %+v", *dm)
		return true
	}
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
	if p.isReverse() {
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
		battlecommon.MoveObject(&p.objectInfo.Pos, moveInfo.Direct, p.gameInfo.GetPanelType(p.objectInfo.ID), true, p.gameInfo.GetPanelInfo)
	case action.MoveTypeAbs:
		target := common.Point{X: moveInfo.AbsPosX, Y: moveInfo.AbsPosY}
		battlecommon.MoveObjectDirect(&p.objectInfo.Pos, target, p.gameInfo.GetPanelType(p.objectInfo.ID), true, p.gameInfo.GetPanelInfo)
	}
}

func (p *Player) AddBuster(busterInfo action.Buster) {
	// TODO: このタイミングで動作させず、アニメーションの追加のみにする

	damageAdd := func(pos common.Point, power int, targetType int) {
		if p.gameInfo.GetPanelInfo(pos).ObjectID != "" {
			logger.Debug("Rock buster damage set %d to (%d, %d)", busterInfo.Power, pos.X, pos.Y)
			damage.New(damage.Damage{
				Pos:           pos,
				Power:         power,
				TTL:           1,
				TargetType:    targetType,
				HitEffectType: 0, // TODO: 正しくセットする
				DamageType:    damage.TypeNone,
			})
		}
	}

	y := p.objectInfo.Pos.Y
	if p.isReverse() {
		for x := p.objectInfo.Pos.X - 1; x >= 0; x-- {
			pos := common.Point{X: x, Y: y}
			damageAdd(pos, busterInfo.Power, damage.TargetPlayer)
		}
	} else {
		for x := p.objectInfo.Pos.X + 1; x < battlecommon.FieldNum.X; x++ {
			pos := common.Point{X: x, Y: y}
			damageAdd(pos, busterInfo.Power, damage.TargetEnemy)
		}
	}
}

func (p *Player) isReverse() bool {
	return p.gameInfo.ReverseClientID == p.objectInfo.OwnerClientID
}
