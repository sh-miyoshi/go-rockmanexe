package player

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	supporterStatusWait int = iota
	supporterStatusMove
	supporterStatusUseChip
	supporterStatusShot
)

type SupporterParam struct {
	HP      uint
	InitPos common.Point
}

type Supporter struct {
	ID              string
	Pos             common.Point
	HP              uint
	HPMax           uint
	ShotPower       uint
	act             act
	invincibleCount int
	status          int
	waitCount       int
	nextStatus      int
}

func NewSupporter(param SupporterParam) (*Supporter, error) {
	logger.Info("Initialize battle supporter")

	res := &Supporter{
		ID:        uuid.New().String(),
		Pos:       param.InitPos,
		HP:        param.HP,
		HPMax:     param.HP,
		ShotPower: 5,
	}
	res.act.typ = -1
	res.act.pPos = &res.Pos

	res.setAction(120, supporterStatusMove)

	return res, nil
}

func (s *Supporter) Draw() {
	view := battlecommon.ViewPos(s.Pos)
	img := s.act.GetImage()
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, img, true)
}

func (s *Supporter) Process() (bool, error) {
	if s.HP <= 0 {
		return true, nil
	}

	if s.invincibleCount > 0 {
		s.invincibleCount--
	}

	if s.act.Process() {
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
			s.act.SetAnim(c.PlayerAct, c.KeepCount)
		}
		target := damage.TargetEnemy
		if c.ForMe {
			target = damage.TargetPlayer
		}

		sid := skill.GetSkillID(c.ID)
		s.act.skillInst = skill.Get(sid, skill.Argument{
			OwnerID:    s.ID,
			Power:      c.Power,
			TargetType: target,
		})
		s.act.skillID = anim.New(s.act.skillInst)
		s.setAction(60, supporterStatusMove)
	case supporterStatusShot:
		s.act.ShotPower = s.ShotPower
		s.act.SetAnim(battlecommon.PlayerActBuster, 0)
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

	if dm.TargetType&damage.TargetPlayer != 0 {
		hp := int(s.HP) - dm.Power
		if hp < 0 {
			s.HP = 0
		} else if hp > int(s.HPMax) {
			s.HP = s.HPMax
		} else {
			s.HP = uint(hp)
		}
		anim.New(effect.Get(dm.HitEffectType, s.Pos, 5))

		for i := 0; i < dm.PushLeft; i++ {
			if !battlecommon.MoveObject(&s.Pos, common.DirectLeft, battlecommon.PanelTypePlayer, true, field.GetPanelInfo) {
				break
			}
		}
		for i := 0; i < dm.PushRight; i++ {
			if !battlecommon.MoveObject(&s.Pos, common.DirectRight, battlecommon.PanelTypePlayer, true, field.GetPanelInfo) {
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

		sound.On(sound.SEDamaged)

		// Stop current animation
		if anim.IsProcessing(s.act.skillID) {
			s.act.skillInst.StopByOwner()
		}
		s.act.skillID = ""

		s.act.SetAnim(battlecommon.PlayerActDamage, 0)
		s.MakeInvisible(battlecommon.PlayerDefaultInvincibleTime)
		logger.Debug("Supporter damaged: %+v", *dm)
		return true
	}
	return false
}

func (s *Supporter) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    s.ID,
			Pos:      s.Pos,
			AnimType: anim.AnimTypeObject,
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

func (s *Supporter) setAction(interval int, next int) {
	s.status = supporterStatusWait
	s.waitCount = interval
	s.nextStatus = next
}

func (s *Supporter) moveRandom() {
	candidates := []int{
		common.DirectUp,
		common.DirectLeft,
		common.DirectDown,
		common.DirectRight,
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
			s.act.SetAnim(battlecommon.PlayerActMove, 0)
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
