package object

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/action"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/queue"
)

type playerAct struct {
	actType       int
	count         int
	pObject       *gameinfo.Object
	info          []byte
	ownerClientID string
	endCount      int
	mgr           *manager.Manager
	fieldFuncs    gameinfo.FieldFuncs
}

type Player struct {
	objectInfo      gameinfo.Object
	hpMax           int
	act             playerAct
	invincibleCount int
	mgr             *manager.Manager
	fieldFuncs      gameinfo.FieldFuncs
	skillID         string
	skillInst       skill.SkillAnim
	actQueueID      string
}

func NewPlayer(info gameinfo.Object, mgr *manager.Manager, funcs gameinfo.FieldFuncs) *Player {
	res := &Player{
		objectInfo: info,
		hpMax:      info.HP,
		act: playerAct{
			actType:       -1,
			ownerClientID: info.OwnerClientID,
			mgr:           mgr,
			fieldFuncs:    funcs,
		},
		invincibleCount: 0,
		mgr:             mgr,
		fieldFuncs:      funcs,
		actQueueID:      uuid.NewString(),
	}
	res.act.pObject = &res.objectInfo
	return res
}

func (p *Player) End() {
	queue.Delete(p.actQueueID)
}

func (p *Player) Update() (bool, error) {
	if p.invincibleCount > 0 {
		p.invincibleCount--
	}

	// Action処理中
	if p.act.Update() {
		return false, nil
	}

	// Actionしてないときは標準ポーズにする
	p.objectInfo.Type = TypePlayerStand

	tact := queue.Pop(p.actQueueID)
	if tact != nil {
		act := tact.(*pb.Request_Action)
		switch act.GetType() {
		case pb.Request_MOVE:
			p.act.SetAnim(battlecommon.PlayerActMove, act.GetRawData(), 0)
		case pb.Request_BUSTER:
			p.act.SetAnim(battlecommon.PlayerActBuster, act.GetRawData(), 0)
		case pb.Request_CHIPUSE:
			var chipInfo action.UseChip
			chipInfo.Unmarshal(act.GetRawData())
			p.useChip(chipInfo)
		default:
			return false, errors.Newf("invalid action type %d is specified", act.GetType())
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

	// インビジ中は無効、ただしRecover系は使えるようにする
	if p.invincibleCount > 0 && dm.Power >= 0 {
		logger.Debug("got damage, but the object is in invincible")
		return false
	}

	// 自分宛のダメージだがObjectが自分じゃない時は無視
	if dm.TargetObjType == damage.TargetPlayer && dm.OwnerClientID != p.objectInfo.OwnerClientID {
		return false
	}

	// 敵宛のダメージだがObjectが自分の時は無視
	if dm.TargetObjType == damage.TargetEnemy && dm.OwnerClientID == p.objectInfo.OwnerClientID {
		return false
	}

	p.objectInfo.HP -= dm.Power
	if p.objectInfo.HP < 0 {
		p.objectInfo.HP = 0
	}
	if p.objectInfo.HP > p.hpMax {
		p.objectInfo.HP = p.hpMax
	}
	if dm.HitEffectType != resources.EffectTypeNone {
		logger.Debug("Add effect %v", dm.HitEffectType)
		p.mgr.QueuePush(gameinfo.QueueTypeEffect, &gameinfo.Effect{
			ID:            uuid.New().String(),
			OwnerClientID: p.act.ownerClientID,
			Pos:           p.objectInfo.Pos,
			Type:          dm.HitEffectType,
			RandRange:     5,
		})
	}

	for i := 0; i < dm.PushLeft; i++ {
		// 敵側から見ると方向は逆になる
		if !battlecommon.MoveObject(&p.objectInfo.Pos, config.DirectRight, battlecommon.PanelTypePlayer, true, p.fieldFuncs.GetPanelInfo) {
			break
		}
	}
	for i := 0; i < dm.PushRight; i++ {
		// 敵側から見ると方向は逆になる
		if !battlecommon.MoveObject(&p.objectInfo.Pos, config.DirectLeft, battlecommon.PanelTypePlayer, true, p.fieldFuncs.GetPanelInfo) {
			break
		}
	}

	if dm.Power <= 0 {
		// Not damage, maybe recover or special anim
		return true
	}

	if dm.StrengthType == damage.StrengthNone {
		return true
	}

	p.mgr.QueuePush(gameinfo.QueueTypeSound, &gameinfo.Sound{
		ID:   uuid.New().String(),
		Type: int(resources.SEDamaged),
	})

	// Stop current animation
	if p.mgr.AnimIsProcessing(p.skillID) {
		p.skillInst.StopByOwner()
	}
	p.skillID = ""

	if dm.IsParalyzed {
		// 麻痺状態(p.act.SetAnim(battlecommon.PlayerActParalyzed, 120))
		system.SetError("TODO: not implemented yet")
	} else {
		p.act.SetAnim(battlecommon.PlayerActDamage, nil, 0)
		if dm.StrengthType == damage.StrengthHigh {
			p.MakeInvisible(battlecommon.PlayerDefaultInvincibleTime)
		}
	}

	logger.Debug("Player damaged: %+v", *dm)
	return true
}

func (p *Player) GetParam() objanim.Param {
	info := NetInfo{
		ActCount:      p.act.count,
		OwnerClientID: p.act.ownerClientID,
		IsInvincible:  p.invincibleCount > 0,
	}

	return objanim.Param{
		Param: anim.Param{
			ObjID:     p.objectInfo.ID,
			Pos:       p.objectInfo.Pos,
			ExtraInfo: info.Marshal(),
		},
		HP: p.objectInfo.HP,
	}
}

func (p *Player) GetObjectType() int {
	return objanim.ObjTypePlayer
}

func (p *Player) MakeInvisible(count int) {
	p.invincibleCount = count
}

func (p *Player) AddBarrier(hp int) {}

func (p *Player) SetCustomGaugeMax() {}

func (p *Player) HandleAction(act *pb.Request_Action) {
	queue.Push(p.actQueueID, act)
}

func (p *Player) GetAnimObjectType() int {
	return p.objectInfo.Type
}

func (p *Player) useChip(chipInfo action.UseChip) {
	c := chip.Get(chipInfo.ChipID)
	logger.Debug("Use chip: %+v", c)
	target := damage.TargetEnemy
	if c.ForMe {
		target = damage.TargetPlayer
	}

	id := skillcore.GetIDByChipID(chipInfo.ChipID)
	s := skill.Get(id, skill.Argument{
		AnimObjID:     chipInfo.AnimID,
		OwnerObjectID: p.objectInfo.ID,
		OwnerClientID: chipInfo.ChipUserClientID,
		Power:         c.Power,
		TargetType:    target,
		Manager:       p.mgr,
		FieldFuncs:    p.fieldFuncs,
	})
	p.skillID = p.mgr.SkillAnimNew(s)
	p.skillInst = s

	if c.PlayerAct != -1 {
		p.act.SetAnim(c.PlayerAct, nil, c.KeepCount)
	}
}

// Process method returns true if processing now
func (a *playerAct) Update() bool {
	a.count++

	switch a.actType {
	case -1: // No animation
		return false
	case battlecommon.PlayerActMove:
		if a.count == 2 {
			var move action.Move
			move.Unmarshal(a.info)

			switch move.Type {
			case action.MoveTypeDirect:
				battlecommon.MoveObject(&a.pObject.Pos, move.Direct, battlecommon.PanelTypePlayer, true, a.fieldFuncs.GetPanelInfo)
			case action.MoveTypeAbs:
				target := point.Point{X: move.AbsPosX, Y: move.AbsPosY}
				battlecommon.MoveObjectDirect(&a.pObject.Pos, target, battlecommon.PanelTypePlayer, true, a.fieldFuncs.GetPanelInfo)
			}

			a.actType = -1
			a.count = 0
			return false
		}
	case battlecommon.PlayerActBuster:
		if a.count == 1 {
			var buster action.Buster
			buster.Unmarshal(a.info)

			damageAdd := func(pos point.Point, power int) bool {
				if objID := a.fieldFuncs.GetPanelInfo(pos).ObjectID; objID != "" {
					logger.Debug("Rock buster damage set %d to (%d, %d)", buster.Power, pos.X, pos.Y)
					eff := resources.EffectTypeHitSmall
					if buster.IsCharged {
						eff = resources.EffectTypeHitBig
					}

					a.mgr.DamageMgr().New(damage.Damage{
						DamageType:    damage.TypeObject,
						OwnerClientID: a.ownerClientID,
						Power:         power,
						TargetObjType: damage.TargetEnemy,
						HitEffectType: eff,
						Element:       damage.ElementNone,
						TargetObjID:   objID,
					})
					return true
				}
				return false
			}

			y := a.pObject.Pos.Y
			for x := a.pObject.Pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				pos := point.Point{X: x, Y: y}
				if damageAdd(pos, buster.Power) {
					break
				}
			}

			a.actType = -1
			a.count = 0
			return false
		}
	case battlecommon.PlayerActCannon, battlecommon.PlayerActSword, battlecommon.PlayerActBomb, battlecommon.PlayerActDamage, battlecommon.PlayerActShot, battlecommon.PlayerActPick, battlecommon.PlayerActThrow:
		// No special action
		if a.count >= a.endCount {
			logger.Info("finished player act %d", a.actType)
			a.actType = -1
			a.count = 0
			a.endCount = 0
			return false
		}
		return true
	default:
		system.SetError(fmt.Sprintf("Invalid player anim type %d was specified.", a.actType))
	}

	return true // processing now
}

func (a *playerAct) SetAnim(actType int, actInfo []byte, keepCount int) {
	a.actType = actType
	a.info = actInfo
	a.count = 0
	// WIP: ソウルユニゾン
	a.endCount = battlecommon.GetPlayerActCount(resources.SoulUnisonNone, actType, keepCount)

	switch actType {
	case battlecommon.PlayerActMove:
		a.pObject.Type = TypePlayerMove
	case battlecommon.PlayerActBuster, battlecommon.PlayerActBShot:
		a.pObject.Type = TypePlayerBuster
	case battlecommon.PlayerActCannon:
		a.pObject.Type = TypePlayerCannon
	case battlecommon.PlayerActShot:
		a.pObject.Type = TypePlayerShot
	case battlecommon.PlayerActSword:
		a.pObject.Type = TypePlayerSword
	case battlecommon.PlayerActDamage:
		a.pObject.Type = TypePlayerDamaged
	case battlecommon.PlayerActPick:
		a.pObject.Type = TypePlayerPick
	case battlecommon.PlayerActBomb:
		a.pObject.Type = TypePlayerBomb
	case battlecommon.PlayerActThrow:
		a.pObject.Type = TypePlayerThrow
	default:
		logger.Error("Invalid player act type %d was specified", actType)
	}
}
