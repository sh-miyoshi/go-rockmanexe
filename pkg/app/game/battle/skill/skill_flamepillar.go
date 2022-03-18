package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	flamePillarStateWakeup int = iota
	flamePillarStateDoing
	flamePillarStateEnd
)

const (
	flamePillarTypeRandom int = iota
	flamePillarTypeTracking
)

type flamePillar struct {
	OwnerID    string
	Power      uint
	TargetType int

	count int
	state int
	point common.Point
}

type flamePillarManager struct {
	ID string

	arg       Argument
	skillType int
	isPlayer  bool
	pillars   []flamePillar
}

func newFlamePillar(objID string, arg Argument, skillType int) *flamePillarManager {
	res := &flamePillarManager{
		ID: objID,

		arg:       arg,
		skillType: skillType,
		isPlayer:  arg.TargetType == damage.TargetEnemy,
		pillars:   []flamePillar{},
	}

	switch skillType {
	case flamePillarTypeRandom:
		panic("not implemented yet")
	case flamePillarTypeTracking:
		pos := objanim.GetObjPos(arg.OwnerID)
		if res.isPlayer {
			pos.X++
		} else {
			pos.X--
		}

		res.pillars = append(res.pillars, flamePillar{
			OwnerID:    arg.OwnerID,
			Power:      arg.Power,
			TargetType: arg.TargetType,
			state:      flamePillarStateWakeup,
			point:      pos,
		})
	}

	return res
}

func (p *flamePillarManager) Draw() {
	for _, pillar := range p.pillars {
		pillar.Draw()
	}
}

func (p *flamePillarManager) Process() (bool, error) {
	switch p.skillType {
	case flamePillarTypeRandom:
		panic("not implemented yet")
	case flamePillarTypeTracking:
		end, err := p.pillars[0].Process()
		if err != nil {
			return false, fmt.Errorf("flame pillar process failed: %w", err)
		}
		if end {
			x := p.pillars[0].point.X
			if p.isPlayer {
				x++
				if x >= field.FieldNum.X {
					return true, nil
				}
			} else {
				x--
				if x < 0 {
					return true, nil
				}
			}
			y := p.pillars[0].point.Y
			// TODO

			p.pillars = append([]flamePillar{}, flamePillar{
				OwnerID:    p.arg.OwnerID,
				Power:      p.arg.Power,
				TargetType: p.arg.TargetType,
				state:      flamePillarStateWakeup,
				point:      common.Point{X: x, Y: y},
			})
		}
	}

	return false, nil
}

func (p *flamePillarManager) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *flamePillar) Draw() {
	view := battlecommon.ViewPos(p.point)

	n := 0
	switch p.state {
	case flamePillarStateWakeup:
		n = p.count / delayFlamePillar
		if n >= len(imgFlamePillar) {
			n = len(imgFlamePillar) - 1
		}
	case flamePillarStateDoing:
		t := (p.count / delayFlamePillar) % 2
		n = len(imgFlamePillar) - (t + 1)
	case flamePillarStateEnd:
		n = len(imgFlamePillar) - (1 + p.count/delayFlamePillar)
		if n < 0 {
			n = 0
		}
	}

	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, imgFlamePillar[n], true)
}

func (p *flamePillar) Process() (bool, error) {
	// TODO damage

	switch p.state {
	case flamePillarStateWakeup:
		if p.count > len(imgFlamePillar)*delayFlamePillar {
			p.count = 0
			p.state = flamePillarStateDoing
			return false, nil
		}
	case flamePillarStateDoing:
		num := 3
		if p.count > num*delayFlamePillar {
			p.count = 0
			p.state = flamePillarStateEnd
			return false, nil
		}
	case flamePillarStateEnd:
		if p.count > len(imgFlamePillar)*delayFlamePillar {
			return true, nil
		}
	}

	p.count++
	return false, nil
}
