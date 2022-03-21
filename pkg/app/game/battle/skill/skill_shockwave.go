package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
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
}

func newShockWave(objID string, isPlayer bool, arg Argument) *shockWave {
	pos := objanim.GetObjPos(arg.OwnerID)
	res := &shockWave{
		ID:     objID,
		Arg:    arg,
		Direct: common.DirectLeft,
		Speed:  5,
		pos:    pos,
	}

	if isPlayer {
		res.Direct = common.DirectRight
		res.Speed = 3
		res.ShowPick = true
		res.InitWait = 9
	}

	return res
}

func (p *shockWave) Draw() {
	n := (p.count / p.Speed) % len(imgShockWave)
	if p.showWave && n >= 0 {
		view := battlecommon.ViewPos(p.pos)
		if p.Direct == common.DirectLeft {
			flag := int32(dxlib.TRUE)
			dxopts := dxlib.DrawRotaGraphOption{
				ReverseXFlag: &flag,
			}
			dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, imgShockWave[n], true, dxopts)
		} else if p.Direct == common.DirectRight {
			dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, imgShockWave[n], true)
		}
	}

	if p.ShowPick {
		n = (p.count / delayPick)
		if n < len(imgPick) {
			pos := objanim.GetObjPos(p.Arg.OwnerID)
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

	n := len(imgShockWave) * p.Speed
	if p.count%(n) == 0 {
		p.showWave = true
		if p.Direct == common.DirectLeft {
			p.pos.X--
		} else if p.Direct == common.DirectRight {
			p.pos.X++
		}

		pn := field.GetPanelInfo(p.pos)
		if pn.Status == field.PanelStatusHole {
			return true, nil
		}

		sound.On(sound.SEShockWave)
		damage.New(damage.Damage{
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           n - 2,
			TargetType:    p.Arg.TargetType,
			HitEffectType: effect.TypeNone,
			ShowHitArea:   true,
			BigDamage:     true,
		})
	}
	p.count++

	if p.pos.X < 0 || p.pos.X > field.FieldNum.X {
		return true, nil
	}
	return false, nil
}

func (p *shockWave) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *shockWave) AtDelete() {
	if p.Arg.RemoveObject != nil {
		p.Arg.RemoveObject(p.ID)
	}
}
