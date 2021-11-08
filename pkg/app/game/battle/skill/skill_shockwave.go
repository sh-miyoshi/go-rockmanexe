package skill

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type shockWave struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int
	Direct     int
	ShowPick   bool
	Speed      int
	InitWait   int

	count    int
	x, y     int
	showWave bool
}

func (p *shockWave) Draw() {
	n := (p.count / p.Speed) % len(imgShockWave)
	if p.showWave && n >= 0 {
		vx, vy := battlecommon.ViewPos(p.x, p.y)
		if p.Direct == common.DirectLeft {
			flag := int32(dxlib.TRUE)
			dxopts := dxlib.DrawRotaGraphOption{
				ReverseXFlag: &flag,
			}
			dxlib.DrawRotaGraph(vx, vy, 1, 0, imgShockWave[n], dxlib.TRUE, dxopts)
		} else if p.Direct == common.DirectRight {
			dxlib.DrawRotaGraph(vx, vy, 1, 0, imgShockWave[n], dxlib.TRUE)
		}
	}

	if p.ShowPick {
		n = (p.count / delayPick)
		if n < len(imgPick) {
			px, py := objanim.GetObjPos(p.OwnerID)
			vx, vy := battlecommon.ViewPos(px, py)
			dxlib.DrawRotaGraph(vx, vy-15, 1, 0, imgPick[n], dxlib.TRUE)
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
		sound.On(sound.SEShockWave)
		if p.Direct == common.DirectLeft {
			p.x--
		} else if p.Direct == common.DirectRight {
			p.x++
		}
		damage.New(damage.Damage{
			PosX:          p.x,
			PosY:          p.y,
			Power:         int(p.Power),
			TTL:           n - 2,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
			ShowHitArea:   true,
			BigDamage:     true,
		})
	}
	p.count++

	if p.x < 0 || p.x > field.FieldNumX {
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
