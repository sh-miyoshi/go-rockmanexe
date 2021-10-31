package skill

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

const (
	wideShotStateBegin int = iota
	wideShotStateMove
)

type wideShot struct {
	ID            string
	OwnerID       string
	Power         uint
	TargetType    int
	Direct        int
	NextStepCount int

	state    int
	count    int
	x, y     int
	damageID [3]string
}

func (p *wideShot) Draw() {
	opt := dxlib.DrawRotaGraphOption{}
	ofs := int32(1)
	if p.Direct == common.DirectLeft {
		xflip := int32(dxlib.TRUE)
		opt.ReverseXFlag = &xflip
		ofs = -1
	}

	switch p.state {
	case wideShotStateBegin:
		x, y := battlecommon.ViewPos(p.x, p.y)
		n := (p.count / delayWideShot)

		if n < len(imgWideShotBody) && p.TargetType == damage.TargetEnemy {
			dxlib.DrawRotaGraph(x+40, y-13, 1, 0, imgWideShotBody[n], dxlib.TRUE, opt)
		}
		if n >= len(imgWideShotBegin) {
			n = len(imgWideShotBegin) - 1
		}
		dxlib.DrawRotaGraph(x+62*ofs, y+20, 1, 0, imgWideShotBegin[n], dxlib.TRUE, opt)
	case wideShotStateMove:
		x, y := battlecommon.ViewPos(p.x, p.y)
		n := (p.count / delayWideShot) % len(imgWideShotMove)
		next := p.x + 1
		prev := p.x - 1
		if p.Direct == common.DirectLeft {
			next, prev = prev, next
		}

		c := p.count % p.NextStepCount
		if c != 0 {
			ofsx := battlecommon.GetOffset(next, p.x, prev, c, p.NextStepCount, field.PanelSizeX)
			dxlib.DrawRotaGraph(x+int32(ofsx), y+20, 1, 0, imgWideShotMove[n], dxlib.TRUE, opt)
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
				p.x++
			} else if p.Direct == common.DirectLeft {
				p.x--
			}

			if p.x >= field.FieldNumX || p.x < 0 {
				return true, nil
			}

			for i := -1; i <= 1; i++ {
				y := p.y + i
				if y < 0 || y >= field.FieldNumY {
					continue
				}

				p.damageID[i+1] = damage.New(damage.Damage{
					PosX:          p.x,
					PosY:          y,
					Power:         int(p.Power),
					TTL:           p.NextStepCount,
					TargetType:    p.TargetType,
					HitEffectType: effect.TypeNone,
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
