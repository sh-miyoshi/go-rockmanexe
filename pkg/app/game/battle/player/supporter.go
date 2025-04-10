package player

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/player/drawer"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	supporterStatusWait int = iota
	supporterStatusMove
	supporterStatusUseChip
	supporterStatusShot
)

type SupporterParam struct {
	HP      uint
	InitPos point.Point
}

type Supporter struct {
	ID              string
	Pos             point.Point
	HP              uint
	HPMax           uint
	ShotPower       uint
	act             BattlePlayerAct
	invincibleCount int
	status          int
	waitCount       int
	nextStatus      int
	animMgr         *manager.Manager
	playerDrawer    drawer.PlayerDrawer
}

func NewSupporter(param SupporterParam, animMgr *manager.Manager) (*Supporter, error) {
	logger.Info("Initialize battle supporter")

	res := &Supporter{
		ID:        uuid.New().String(),
		Pos:       param.InitPos,
		HP:        param.HP,
		HPMax:     param.HP,
		ShotPower: 5,
		animMgr:   animMgr,
	}
	res.act.Init(&res.Pos, animMgr)

	if err := res.playerDrawer.Init(); err != nil {
		return nil, err
	}

	res.setAction(120, supporterStatusMove)

	return res, nil
}

func (s *Supporter) Draw() {
	view := battlecommon.ViewPos(s.Pos)
	cnt, typ := s.act.GetParams()
	s.playerDrawer.Draw(cnt, view, typ, s.act.IsParalyzed())
}

func (s *Supporter) Update() (bool, error) {
	if s.HP <= 0 {
		return true, nil
	}

	if s.invincibleCount > 0 {
		s.invincibleCount--
	}

	if s.act.Update() {
		return false, nil
	}

	switch s.status {
	case supporterStatusWait:
		s.waitCount--
		if s.waitCount <= 0 {
			s.status = s.nextStatus
		}
	case supporterStatusMove:
		s.moveRandom()
		s.decideNextAction()
	case supporterStatusUseChip:
		targetChip := chip.IDSpreadGun
		c := chip.Get(targetChip)
		if c.PlayerAct != -1 {
			s.act.SetAnim(resources.SoulUnisonNone, c.PlayerAct, c.KeepCount)
		}
		target := damage.TargetEnemy
		if c.ForMe {
			target = damage.TargetPlayer
		}

		sid := skillcore.GetIDByChipID(c.ID)
		s.act.SetSkill(sid, skillcore.Argument{
			OwnerID:    s.ID,
			Power:      c.Power,
			TargetType: target,
		})
		s.setAction(60, supporterStatusMove)
	case supporterStatusShot:
		s.act.ShotPower = s.ShotPower
		s.act.SetAnim(resources.SoulUnisonNone, battlecommon.PlayerActBuster, 0)
		s.setAction(60, supporterStatusMove)
	}

	return false, nil
}

func (s *Supporter) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}

	// Recover系は使えるようにする
	if s.invincibleCount > 0 && dm.Power >= 0 {
		return false
	}

	if dm.TargetObjType&damage.TargetPlayer != 0 {
		hp := int(s.HP) - dm.Power
		if hp < 0 {
			s.HP = 0
		} else if hp > int(s.HPMax) {
			s.HP = s.HPMax
		} else {
			s.HP = uint(hp)
		}
		s.animMgr.EffectAnimNew(effect.Get(dm.HitEffectType, s.Pos, 5))

		for i := 0; i < dm.PushLeft; i++ {
			if !battlecommon.MoveObject(&s.Pos, config.DirectLeft, battlecommon.PanelTypePlayer, true, field.GetPanelInfo) {
				break
			}
		}
		for i := 0; i < dm.PushRight; i++ {
			if !battlecommon.MoveObject(&s.Pos, config.DirectRight, battlecommon.PanelTypePlayer, true, field.GetPanelInfo) {
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

		sound.On(resources.SEDamaged)

		// Stop current animation
		if s.animMgr.IsAnimProcessing(s.act.skillObjID) {
			s.act.skillInst.StopByOwner()
		}
		s.act.skillObjID = ""

		if dm.IsParalyzed {
			s.act.SetAnim(resources.SoulUnisonNone, battlecommon.PlayerActParalyzed, battlecommon.DefaultParalyzedTime)
		} else {
			s.act.SetAnim(resources.SoulUnisonNone, battlecommon.PlayerActDamage, 0)
			if dm.StrengthType == damage.StrengthHigh {
				s.MakeInvisible(battlecommon.PlayerDefaultInvincibleTime)
			}
		}
		logger.Debug("Supporter damaged: %+v", *dm)
		return true
	}
	return false
}

func (s *Supporter) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID: s.ID,
			Pos:   s.Pos,
		},
		HP: int(s.HP),
	}
}

func (s *Supporter) GetObjectType() int {
	return objanim.ObjTypePlayer
}

func (s *Supporter) MakeInvisible(count int) {
	s.invincibleCount = count
}

func (s *Supporter) AddBarrier(hp int) {}

func (s *Supporter) SetCustomGaugeMax() {}

func (s *Supporter) setAction(interval int, next int) {
	s.status = supporterStatusWait
	s.waitCount = interval
	s.nextStatus = next
}

func (s *Supporter) moveRandom() {
	candidates := []int{
		config.DirectUp,
		config.DirectLeft,
		config.DirectDown,
		config.DirectRight,
	}
	// shuffule candidates
	for i := 0; i < 10; i++ {
		for j := 0; j < len(candidates); j++ {
			n := rand.Intn(len(candidates))
			candidates[j], candidates[n] = candidates[n], candidates[j]
		}
	}

	for _, direct := range candidates {
		if battlecommon.MoveObject(&s.Pos, direct, battlecommon.PanelTypePlayer, false, field.GetPanelInfo) {
			s.act.MoveDirect = direct
			s.act.SetAnim(resources.SoulUnisonNone, battlecommon.PlayerActMove, 0)
			return
		}
	}
}

func (s *Supporter) decideNextAction() {
	n := rand.Intn(100)
	wait := 30 + rand.Intn(60)
	if n < 40 {
		s.setAction(wait, supporterStatusMove)
	} else if n < 80 {
		s.setAction(wait, supporterStatusShot)
	} else {
		s.setAction(wait, supporterStatusUseChip)
	}
}
