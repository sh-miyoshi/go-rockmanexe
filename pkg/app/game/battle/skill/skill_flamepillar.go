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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	flamePillarStateWakeup int = iota
	flamePillarStateDoing
	flamePillarStateEnd
	flamePillarStateDeleted
)

const (
	flamePillarTypeRandom int = iota
	flamePillarTypeTracking
	flamePillarTypeLine
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
	count     int
	skillType int
	isPlayer  bool
	pillars   []*flamePillar
}

func newFlamePillar(objID string, arg Argument, skillType int) *flamePillarManager {
	res := &flamePillarManager{
		ID: objID,

		Arg:       arg,
		count:     0,
		skillType: skillType,
		isPlayer:  arg.TargetType == damage.TargetEnemy,
		pillars:   []*flamePillar{},
	}

	switch skillType {
	case flamePillarTypeRandom:
		common.SetError("TODO: not implemented yet")
	case flamePillarTypeTracking:
		pos := objanim.GetObjPos(arg.OwnerID)
		if res.isPlayer {
			pos.X++
		} else {
			pos.X--
		}

		res.pillars = append(res.pillars, &flamePillar{
			Arg:   arg,
			state: flamePillarStateWakeup,
			point: pos,
		})
	case flamePillarTypeLine:
		posX := objanim.GetObjPos(arg.OwnerID).X
		if res.isPlayer {
			posX += 2
		} else {
			posX -= 2
		}

		for y := 0; y < field.FieldNum.Y; y++ {
			res.pillars = append(res.pillars, &flamePillar{
				Arg:   arg,
				state: flamePillarStateWakeup,
				point: common.Point{X: posX, Y: y},
			})
		}
	}

	return res
}

func (p *flamePillarManager) Draw() {
	for _, pillar := range p.pillars {
		pillar.Draw()
	}

	if p.skillType == flamePillarTypeLine {
		if len(p.pillars) > 0 && p.pillars[0].state != flamePillarStateEnd {
			imageNo := p.count / 4
			if imageNo >= len(imgFlameLineBody) {
				imageNo = len(imgFlameLineBody) - 1
			}

			pos := objanim.GetObjPos(p.Arg.OwnerID)
			view := battlecommon.ViewPos(pos)

			// Show body
			dxlib.DrawRotaGraph(view.X+35, view.Y-15, 1, 0, imgFlameLineBody[imageNo], true)
		}
	}
}

func (p *flamePillarManager) Process() (bool, error) {
	switch p.skillType {
	case flamePillarTypeRandom:
		common.SetError("TODO: not implemented yet")
	case flamePillarTypeTracking:
		end, err := p.pillars[0].Process()
		if err != nil {
			return false, fmt.Errorf("flame pillar process failed: %w", err)
		}
		if end {
			// 穴パネルなどで進めなかったら終わり
			if p.pillars[0].state == flamePillarStateDeleted {
				return true, nil
			}

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

			p.pillars = append([]*flamePillar{}, &flamePillar{
				Arg:   p.Arg,
				state: flamePillarStateWakeup,
				point: common.Point{X: x, Y: y},
			})
		}
	case flamePillarTypeLine:
		remove := []int{}
		for i, pillar := range p.pillars {
			end, err := pillar.Process()
			if err != nil {
				return false, fmt.Errorf("flame pillar process failed: %w", err)
			}
			if end {
				remove = append(remove, i)
			}
		}
		// Remove finished pillars
		sort.Sort(sort.Reverse(sort.IntSlice(remove)))
		for _, index := range remove {
			p.pillars[index] = p.pillars[len(p.pillars)-1]
			p.pillars = p.pillars[:len(p.pillars)-1]
		}
		if len(p.pillars) == 0 {
			return true, nil
		}
	}

	p.count++
	return false, nil
}

func (p *flamePillarManager) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
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
		if p.count == 0 {
			pn := field.GetPanelInfo(p.point)
			if pn.Status == field.PanelStatusHole {
				p.state = flamePillarStateDeleted
				return true, nil
			}

			sound.On(sound.SEFlameAttack)
		}

		if p.count == 3*delayFlamePillar {
			// Add damage
			damage.New(damage.Damage{
				Pos:         p.point,
				Power:       int(p.Arg.Power),
				TTL:         7 * delayFlamePillar,
				TargetType:  p.Arg.TargetType,
				ShowHitArea: true,
				BigDamage:   true,
				DamageType:  damage.TypeFire,
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
	case flamePillarStateDeleted:
		// Nothing to do
	}

	p.count++
	return false, nil
}
