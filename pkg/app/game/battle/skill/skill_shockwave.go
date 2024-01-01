package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type shockWave struct {
	ID       string
	Arg      skillcore.Argument
	Direct   int
	ShowPick bool
	Speed    int
	InitWait int

	count      int
	pos        point.Point
	showWave   bool
	drawer     skilldraw.DrawShockWave
	pickDrawer skilldraw.DrawPick
}

func newShockWave(objID string, isPlayer bool, arg skillcore.Argument) *shockWave {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	res := &shockWave{
		ID:     objID,
		Arg:    arg,
		Direct: config.DirectLeft,
		Speed:  5,
		pos:    pos,
	}

	if isPlayer {
		res.Direct = config.DirectRight
		res.Speed = resources.SkillShockWavePlayerSpeed
		res.ShowPick = true
		res.InitWait = resources.SkillShockWaveInitWait
	}

	return res
}

func (p *shockWave) Draw() {
	if p.showWave {
		view := battlecommon.ViewPos(p.pos)
		p.drawer.Draw(view, p.count, p.Speed, p.Direct)
	}

	if p.ShowPick {
		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		view := battlecommon.ViewPos(pos)
		p.pickDrawer.Draw(view, p.count)
	}
}

func (p *shockWave) Process() (bool, error) {
	if p.count < p.InitWait {
		p.count++
		return false, nil
	}

	n := resources.SkillShockWaveImageNum * p.Speed
	if p.count%n == 0 {
		p.showWave = true
		if p.Direct == config.DirectLeft {
			p.pos.X--
		} else if p.Direct == config.DirectRight {
			p.pos.X++
		}

		pn := field.GetPanelInfo(p.pos)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		sound.On(resources.SEShockWave)
		localanim.DamageManager().New(damage.Damage{
			DamageType:    damage.TypePosition,
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           n - 2,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			ShowHitArea:   true,
			BigDamage:     true,
			Element:       damage.ElementNone,
		})
	}
	p.count++

	if p.pos.X < 0 || p.pos.X > battlecommon.FieldNum.X {
		return true, nil
	}
	return false, nil
}

func (p *shockWave) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *shockWave) StopByOwner() {
	if p.count <= p.InitWait {
		localanim.AnimDelete(p.ID)
	}
}
