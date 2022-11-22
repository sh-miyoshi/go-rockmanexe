package player

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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
	ChipFolder      []player.ChipInfo
	SelectedChips   []player.ChipInfo
	act             act
	invincibleCount int
}

func NewSupporter(param SupporterParam) (*Supporter, error) {
	logger.Info("Initialize battle supporter")

	res := &Supporter{
		ID:        uuid.New().String(),
		Pos:       param.InitPos,
		HP:        param.HP,
		HPMax:     param.HP,
		ShotPower: 1,
		// TODO: ChipFolder
	}
	res.act.typ = -1
	res.act.pPos = &res.Pos

	return res, nil
}

func (s *Supporter) Draw() {
	view := battlecommon.ViewPos(s.Pos)
	img := s.act.GetImage()
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, img, true)
}

func (s *Supporter) Process() (bool, error) {
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
			if !battlecommon.MoveObject(&s.Pos, common.DirectLeft, field.PanelTypePlayer, true, field.GetPanelInfo) {
				break
			}
		}
		for i := 0; i < dm.PushRight; i++ {
			if !battlecommon.MoveObject(&s.Pos, common.DirectRight, field.PanelTypePlayer, true, field.GetPanelInfo) {
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
		// if anim.IsProcessing(s.act.skillID) {
		// 	s.act.skillInst.StopByOwner()
		// }
		// s.act.skillID = ""

		// s.act.SetAnim(battlecommon.PlayerActDamage, 0)
		// s.DamageNum++
		s.MakeInvisible(battlecommon.PlayerDefaultInvincibleTime)
		logger.Debug("Supporter damaged: %+v", *dm)
		return true
	}
	return false
}

func (s *Supporter) GetParam() anim.Param {
	return anim.Param{
		ObjID:    s.ID,
		Pos:      s.Pos,
		AnimType: anim.AnimTypeObject,
	}
}
func (s *Supporter) GetObjectType() int {
	return objanim.ObjTypePlayer
}
func (s *Supporter) MakeInvisible(count int) {
	s.invincibleCount = count
}
