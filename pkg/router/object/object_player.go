package object

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/action"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/queue"
)

type Player struct {
	objectInfo    gameinfo.Object
	gameInfo      *gameinfo.GameInfo
	actionQueueID string
}

func NewPlayer(info gameinfo.Object, gameInfo *gameinfo.GameInfo, actionQueueID string) *Player {
	return &Player{
		objectInfo:    info,
		gameInfo:      gameInfo,
		actionQueueID: actionQueueID,
	}
}

func (p *Player) GetCurrentObjectTypePointer() *int {
	return &p.objectInfo.Type
}

func (p *Player) Process() (bool, error) {
	act := queue.Pop(p.actionQueueID)
	if act != nil {
		switch act.GetType() {
		case pb.Request_MOVE:
			var move action.Move
			move.Unmarshal(act.GetRawData())

			p.addMove(move)
		case pb.Request_BUSTER:
			var buster action.Buster
			buster.Unmarshal(act.GetRawData())
			p.addBuster(buster)
		case pb.Request_CHIPUSE:
			var chipInfo action.UseChip
			chipInfo.Unmarshal(act.GetRawData())
			p.useChip(chipInfo)
		default:
			return false, fmt.Errorf("invalid action type %d is specified", act.GetType())
		}
	}

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

	if dm.TargetType&damage.TargetPlayer != 0 {
		p.objectInfo.HP -= dm.Power
		if p.objectInfo.HP < 0 {
			p.objectInfo.HP = 0
		}
		// TODO: HPMax
		// TODO: damage effect

		for i := 0; i < dm.PushLeft; i++ {
			if !battlecommon.MoveObject(&p.objectInfo.Pos, common.DirectLeft, battlecommon.PanelTypePlayer, true, p.gameInfo.GetPanelInfo) {
				break
			}
		}
		for i := 0; i < dm.PushRight; i++ {
			if !battlecommon.MoveObject(&p.objectInfo.Pos, common.DirectRight, battlecommon.PanelTypePlayer, true, p.gameInfo.GetPanelInfo) {
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
	return objanim.ObjTypePlayer
}

func (p *Player) MakeInvisible(count int) {
	// TODO
}

func (p *Player) addMove(moveInfo action.Move) {
	// TODO: このタイミングで移動せず、アニメーションの追加のみにする
	switch moveInfo.Type {
	case action.MoveTypeDirect:
		battlecommon.MoveObject(&p.objectInfo.Pos, moveInfo.Direct, battlecommon.PanelTypePlayer, true, p.gameInfo.GetPanelInfo)
	case action.MoveTypeAbs:
		target := common.Point{X: moveInfo.AbsPosX, Y: moveInfo.AbsPosY}
		battlecommon.MoveObjectDirect(&p.objectInfo.Pos, target, battlecommon.PanelTypePlayer, true, p.gameInfo.GetPanelInfo)
	}
}

func (p *Player) addBuster(busterInfo action.Buster) {
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
	for x := p.objectInfo.Pos.X + 1; x < battlecommon.FieldNum.X; x++ {
		pos := common.Point{X: x, Y: y}
		damageAdd(pos, busterInfo.Power, damage.TargetEnemy)
	}
}

func (p *Player) useChip(chipInfo action.UseChip) {
	/*
		TODO

		c := chip.Get(chipInfo.ChipID)
		logger.Debug("Use chip: %+v", c)

		var targetType int
		if g.info.ReverseClientID == clientID {
			if c.ForMe {
				targetType = damage.TargetEnemy
			} else {
				targetType = damage.TargetPlayer
			}
		} else {
			if c.ForMe {
				targetType = damage.TargetPlayer
			} else {
				targetType = damage.TargetEnemy
			}
		}

		s := skill.GetByChip(chipInfo.ChipID, skill.Argument{
			AnimObjID:  chipInfo.AnimID,
			OwnerID:    chipInfo.ChipUserClientID,
			Power:      c.Power,
			TargetType: targetType,

			GameInfo: &g.info,
		})
		anim.New(s)

		// TODO: player_act
	*/
}
