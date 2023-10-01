package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayPick = 3
)

type shockWave struct {
	ID       string
	Arg      Argument
	Direct   int
	ShowPick bool
	Speed    int
	InitWait int

	count    int
	pos      common.Point
	showWave bool
	drawer   skilldraw.DrawShockWave
}

func newShockWave(objID string, isPlayer bool, arg Argument) *shockWave {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	res := &shockWave{
		ID:     objID,
		Arg:    arg,
		Direct: common.DirectLeft,
		Speed:  5,
		pos:    pos,
	}

	res.drawer.Init() // TODO: error

	if isPlayer {
		res.Direct = common.DirectRight
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
		n := (p.count / delayPick)
		if n < len(imgPick) {
			pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
			view := battlecommon.ViewPos(pos)
			dxlib.DrawRotaGraph(view.X, view.Y-15, 1, 0, imgPick[n], true)
		}
	}
}

func (p *shockWave) Process() (bool, error) {
	if p.count < p.InitWait {
		p.count++
		return false, nil
	}

	n := resources.SkillShockWaveImageNum * p.Speed
	if p.count%(n) == 0 {
		p.showWave = true
		if p.Direct == common.DirectLeft {
			p.pos.X--
		} else if p.Direct == common.DirectRight {
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
