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
	wideShotStateBegin int = iota
	wideShotStateMove
)

type wideShot struct {
	ID            string
	Arg           Argument
	Direct        int
	NextStepCount int

	state    int
	count    int
	pos      common.Point
	damageID [3]string
}

func newWideShot(objID string, arg Argument) *wideShot {
	pos := objanim.GetObjPos(arg.OwnerID)
	direct := common.DirectRight
	nextStep := 8
	if arg.TargetType == damage.TargetPlayer {
		direct = common.DirectLeft
		nextStep = 16
	}

	return &wideShot{
		ID:            objID,
		Arg:           arg,
		Direct:        direct,
		NextStepCount: nextStep,
		pos:           pos,
		state:         wideShotStateBegin,
	}
}

func (p *wideShot) Draw() {
	opt := dxlib.DrawRotaGraphOption{}
	ofs := 1
	if p.Direct == common.DirectLeft {
		xflip := int32(dxlib.TRUE)
		opt.ReverseXFlag = &xflip
		ofs = -1
	}

	switch p.state {
	case wideShotStateBegin:
		view := battlecommon.ViewPos(p.pos)
		n := (p.count / delayWideShot)

		if n < len(imgWideShotBody) && p.Arg.TargetType == damage.TargetEnemy {
			dxlib.DrawRotaGraph(view.X+40, view.Y-13, 1, 0, imgWideShotBody[n], true, opt)
		}
		if n >= len(imgWideShotBegin) {
			n = len(imgWideShotBegin) - 1
		}
		dxlib.DrawRotaGraph(view.X+62*ofs, view.Y+20, 1, 0, imgWideShotBegin[n], true, opt)
	case wideShotStateMove:
		view := battlecommon.ViewPos(p.pos)
		n := (p.count / delayWideShot) % len(imgWideShotMove)
		next := p.pos.X + 1
		prev := p.pos.X - 1
		if p.Direct == common.DirectLeft {
			next, prev = prev, next
		}

		c := p.count % p.NextStepCount
		if c != 0 {
			ofsx := battlecommon.GetOffset(next, p.pos.X, prev, c, p.NextStepCount, field.PanelSize.X)
			dxlib.DrawRotaGraph(view.X+ofsx, view.Y+20, 1, 0, imgWideShotMove[n], true, opt)
		}
	}
}

func (p *wideShot) Process() (bool, error) {
	for _, did := range p.damageID {
		if did != "" {
			if !damage.Exists(did) && p.count%p.NextStepCount != 0 {
				// attack hit to target
				return true, nil
			}
		}
	}

	switch p.state {
	case wideShotStateBegin:
		if p.count == 0 {
			sound.On(sound.SEWideShot)
		}

		max := len(imgWideShotBody)
		if len(imgWideShotBegin) > max {
			max = len(imgWideShotBegin)
		}
		max *= delayWideShot
		if p.count > max {
			p.state = wideShotStateMove
			p.count = 0
			return false, nil
		}
	case wideShotStateMove:
		if p.count%p.NextStepCount == 0 {
			if p.Direct == common.DirectRight {
				p.pos.X++
			} else if p.Direct == common.DirectLeft {
				p.pos.X--
			}

			if p.pos.X >= field.FieldNum.X || p.pos.X < 0 {
				return true, nil
			}

			for i := -1; i <= 1; i++ {
				y := p.pos.Y + i
				if y < 0 || y >= field.FieldNum.Y {
					continue
				}

				p.damageID[i+1] = damage.New(damage.Damage{
					Pos:           common.Point{X: p.pos.X, Y: y},
					Power:         int(p.Arg.Power),
					TTL:           p.NextStepCount,
					TargetType:    p.Arg.TargetType,
					HitEffectType: effect.TypeNone,
					BigDamage:     true,
				})
			}
		}
	}

	p.count++
	return false, nil
}

func (p *wideShot) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *wideShot) AtDelete() {
	if p.Arg.AtDelete != nil {
		p.Arg.AtDelete()
	}
}
