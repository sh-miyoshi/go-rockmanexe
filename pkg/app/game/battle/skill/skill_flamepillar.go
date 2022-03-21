package skill

import (
	"fmt"
	"sort"

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

const (
	delayFlamePillar = 4
)

type flamePillar struct {
	Arg Argument

	count int
	state int
	point common.Point
}

type flamePillarManager struct {
	ID string

	Arg       Argument
	skillType int
	isPlayer  bool
	pillars   []flamePillar
}

func newFlamePillar(objID string, arg Argument, skillType int) *flamePillarManager {
	res := &flamePillarManager{
		ID: objID,

		Arg:       arg,
		skillType: skillType,
		isPlayer:  arg.TargetType == damage.TargetEnemy,
		pillars:   []flamePillar{},
	}

	switch skillType {
	case flamePillarTypeRandom:
		panic("TODO: not implemented yet")
	case flamePillarTypeTracking:
		pos := objanim.GetObjPos(arg.OwnerID)
		if res.isPlayer {
			pos.X++
		} else {
			pos.X--
		}

		res.pillars = append(res.pillars, flamePillar{
			Arg:   arg,
			state: flamePillarStateWakeup,
			point: pos,
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
		panic("TODO: not implemented yet")
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
			objType := objanim.ObjTypePlayer
			if p.Arg.TargetType == damage.TargetEnemy {
				objType = objanim.ObjTypeEnemy
			}
			objs := objanim.GetObjs(objanim.Filter{ObjType: objType})
			if len(objs) > 0 {
				sort.Slice(objs, func(i, j int) bool {
					return objs[i].Pos.X < objs[j].Pos.X
				})

				if objs[0].Pos.Y > y {
					y++
				} else if objs[0].Pos.Y < y {
					y--
				}
			}

			p.pillars = append([]flamePillar{}, flamePillar{
				Arg:   p.Arg,
				state: flamePillarStateWakeup,
				point: common.Point{X: x, Y: y},
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

func (p *flamePillarManager) AtDelete() {
	if p.Arg.RemoveObject != nil {
		p.Arg.RemoveObject(p.ID)
	}
}

func (p *flamePillarManager) StopByOwner() {
	anim.Delete(p.ID)
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
	switch p.state {
	case flamePillarStateWakeup:
		if p.count == 3*delayFlamePillar {
			// Add damage
			damage.New(damage.Damage{
				Pos:         p.point,
				Power:       int(p.Arg.Power),
				TTL:         7 * delayFlamePillar,
				TargetType:  p.Arg.TargetType,
				ShowHitArea: true,
				BigDamage:   true,
			})
		}

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
