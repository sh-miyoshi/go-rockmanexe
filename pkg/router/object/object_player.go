package object

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/action"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/queue"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/skill"
)

type playerAct struct {
	actType       int
	count         int
	pObject       *gameinfo.Object
	info          []byte
	getPanelInfo  func(pos common.Point) battlecommon.PanelInfo
	ownerClientID string
	endCount      int
}

type Player struct {
	objectInfo      gameinfo.Object
	gameInfo        *gameinfo.GameInfo
	actionQueueID   string
	hpMax           int
	act             playerAct
	invincibleCount int
}

func NewPlayer(info gameinfo.Object, gameInfo *gameinfo.GameInfo, actionQueueID string) *Player {
	res := &Player{
		objectInfo:    info,
		gameInfo:      gameInfo,
		actionQueueID: actionQueueID,
		hpMax:         info.HP,
		act: playerAct{
			actType:       -1,
			ownerClientID: info.OwnerClientID,
		},
		invincibleCount: 0,
	}
	res.act.pObject = &res.objectInfo
	res.act.getPanelInfo = res.gameInfo.GetPanelInfo

	return res
}

func (p *Player) GetCurrentObjectTypePointer() *int {
	return &p.objectInfo.Type
}

func (p *Player) Process() (bool, error) {
	if p.invincibleCount > 0 {
		p.invincibleCount--
	}

	// Action処理中
	if p.act.Process() {
		return false, nil
	}

	// Actionしてないときは標準ポーズにする
	p.objectInfo.Type = TypePlayerStand

	act := queue.Pop(p.actionQueueID)
	if act != nil {
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

	// インビジ中は無効、ただしRecover系は使えるようにする
	if p.invincibleCount > 0 && dm.Power >= 0 {
		logger.Debug("got damage, but the object is in invincible")
		return false
	}

	// 自分宛のダメージだがObjectが自分じゃない時は無視
	if dm.TargetType == damage.TargetPlayer && dm.OwnerClientID != p.objectInfo.OwnerClientID {
		return false
	}

	// 敵宛のダメージだがObjectが自分の時は無視
	if dm.TargetType == damage.TargetEnemy && dm.OwnerClientID == p.objectInfo.OwnerClientID {
		return false
	}

	p.objectInfo.HP -= dm.Power
	if p.objectInfo.HP < 0 {
		p.objectInfo.HP = 0
	}
	if p.objectInfo.HP > p.hpMax {
		p.objectInfo.HP = p.hpMax
	}
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
	p.MakeInvisible(battlecommon.PlayerDefaultInvincibleTime)
	logger.Debug("Player damaged: %+v", *dm)
	return true
}

func (p *Player) GetParam() objanim.Param {
	info := NetInfo{
		ActCount:      p.act.count,
		OwnerClientID: p.gameInfo.ClientID,
		IsInvincible:  p.invincibleCount > 0,
	}

	return objanim.Param{
		Param: anim.Param{
			ObjID:     p.objectInfo.ID,
			Pos:       p.objectInfo.Pos,
			DrawType:  anim.DrawTypeObject,
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

func (p *Player) useChip(chipInfo action.UseChip) {
	c := chip.Get(chipInfo.ChipID)
	logger.Debug("Use chip: %+v", c)
	target := damage.TargetEnemy
	if c.ForMe {
		target = damage.TargetPlayer
	}

	s := skill.GetByChip(chipInfo.ChipID, skill.Argument{
		AnimObjID:     chipInfo.AnimID,
		OwnerObjectID: p.objectInfo.ID,
		OwnerClientID: chipInfo.ChipUserClientID,
		Power:         c.Power,
		TargetType:    target,

		GameInfo: p.gameInfo,
	})
	routeranim.AnimNew(chipInfo.ChipUserClientID, s)

	if c.PlayerAct != -1 {
		p.act.SetAnim(c.PlayerAct, nil, s.GetEndCount())
	}
}

// Process method returns true if processing now
func (a *playerAct) Process() bool {
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
				battlecommon.MoveObject(&a.pObject.Pos, move.Direct, battlecommon.PanelTypePlayer, true, a.getPanelInfo)
			case action.MoveTypeAbs:
				target := common.Point{X: move.AbsPosX, Y: move.AbsPosY}
				battlecommon.MoveObjectDirect(&a.pObject.Pos, target, battlecommon.PanelTypePlayer, true, a.getPanelInfo)
			}

			a.actType = -1
			a.count = 0
			return false
		}
	case battlecommon.PlayerActBuster:
		if a.count == 1 {
			var buster action.Buster
			buster.Unmarshal(a.info)

			damageAdd := func(pos common.Point, power int) bool {
				if a.getPanelInfo(pos).ObjectID != "" {
					logger.Debug("Rock buster damage set %d to (%d, %d)", buster.Power, pos.X, pos.Y)
					pos.X = battlecommon.FieldNum.X - pos.X - 1
					eff := battlecommon.EffectTypeHitSmall
					if buster.IsCharged {
						eff = battlecommon.EffectTypeHitBig
					}

					routeranim.DamageManager(a.ownerClientID).New(damage.Damage{
						OwnerClientID: a.ownerClientID,
						Pos:           pos,
						Power:         power,
						TTL:           1,
						TargetType:    damage.TargetEnemy,
						HitEffectType: eff,
						DamageType:    damage.TypeNone,
					})
					return true
				}
				return false
			}

			y := a.pObject.Pos.Y
			for x := a.pObject.Pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				pos := common.Point{X: x, Y: y}
				if damageAdd(pos, buster.Power) {
					break
				}
			}

			a.actType = -1
			a.count = 0
			return false
		}
	case battlecommon.PlayerActCannon, battlecommon.PlayerActSword, battlecommon.PlayerActBomb, battlecommon.PlayerActDamage, battlecommon.PlayerActShot, battlecommon.PlayerActPick, battlecommon.PlayerActThrow:
		// TODO set pick anim end count

		// No special action
		if a.count >= a.endCount {
			a.actType = -1
			a.count = 0
			a.endCount = 0
			return false
		}
		return true
	default:
		panic(fmt.Sprintf("Invalid player anim type %d was specified.", a.actType))
	}

	return true // processing now
}

func (a *playerAct) SetAnim(actType int, actInfo []byte, endCount int) {
	a.actType = actType
	a.info = actInfo
	a.count = 0
	a.endCount = endCount

	switch actType {
	case battlecommon.PlayerActMove:
		a.pObject.Type = TypePlayerMove
	case battlecommon.PlayerActBuster:
		a.pObject.Type = TypePlayerBuster
	case battlecommon.PlayerActCannon:
		a.pObject.Type = TypePlayerCannon
	case battlecommon.PlayerActShot:
		a.pObject.Type = TypePlayerShot
	case battlecommon.PlayerActSword:
		a.pObject.Type = TypePlayerSword
	default:
		logger.Error("Invalid player act type %d was specified", actType)
	}
}
