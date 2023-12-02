package skill

import (
	"fmt"
	"sort"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type flamePillar struct {
	Arg Argument

	count  int
	state  int
	point  point.Point
	drawer skilldraw.DrawFlamePiller
}

type flamePillarManager struct {
	ID string

	Arg       Argument
	count     int
	skillType int
	isPlayer  bool
	pillars   []*flamePillar
	drawer    skilldraw.DrawFlamePillerManager
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
	case resources.SkillFlamePillarTypeRandom:
		common.SetError("TODO: not implemented yet")
	case resources.SkillFlamePillarTypeTracking:
		pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
		if res.isPlayer {
			pos.X++
		} else {
			pos.X--
		}

		res.pillars = append(res.pillars, &flamePillar{
			Arg:   arg,
			state: resources.SkillFlamePillarStateWakeup,
			point: pos,
		})
	case resources.SkillFlamePillarTypeLine:
		posX := localanim.ObjAnimGetObjPos(arg.OwnerID).X
		if res.isPlayer {
			posX += 2
		} else {
			posX -= 2
		}

		for y := 0; y < battlecommon.FieldNum.Y; y++ {
			res.pillars = append(res.pillars, &flamePillar{
				Arg:   arg,
				state: resources.SkillFlamePillarStateWakeup,
				point: point.Point{X: posX, Y: y},
			})
		}
	}

	return res
}

func (p *flamePillarManager) Draw() {
	for _, pillar := range p.pillars {
		pillar.Draw()
	}

	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	state := resources.SkillFlamePillarStateEnd
	if len(p.pillars) > 0 {
		state = p.pillars[0].state
	}
	p.drawer.Draw(view, p.count, p.skillType, state)
}

func (p *flamePillarManager) Process() (bool, error) {
	switch p.skillType {
	case resources.SkillFlamePillarTypeRandom:
		common.SetError("TODO: not implemented yet")
	case resources.SkillFlamePillarTypeTracking:
		end, err := p.pillars[0].Process()
		if err != nil {
			return false, fmt.Errorf("flame pillar process failed: %w", err)
		}
		if end {
			// 穴パネルなどで進めなかったら終わり
			if p.pillars[0].state == resources.SkillFlamePillarStateDeleted {
				return true, nil
			}

			x := p.pillars[0].point.X
			if p.isPlayer {
				x++
				if x >= battlecommon.FieldNum.X {
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
			objs := localanim.ObjAnimGetObjs(objanim.Filter{ObjType: objType})
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
				state: resources.SkillFlamePillarStateWakeup,
				point: point.Point{X: x, Y: y},
			})
		}
	case resources.SkillFlamePillarTypeLine:
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
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *flamePillarManager) StopByOwner() {
	localanim.AnimDelete(p.ID)
}

func (p *flamePillar) Draw() {
	view := battlecommon.ViewPos(p.point)
	p.drawer.Draw(view, p.count, p.state)
}

func (p *flamePillar) Process() (bool, error) {
	switch p.state {
	case resources.SkillFlamePillarStateWakeup:
		if p.count == 0 {
			pn := field.GetPanelInfo(p.point)
			if pn.Status == battlecommon.PanelStatusHole {
				p.state = resources.SkillFlamePillarStateDeleted
				return true, nil
			}

			sound.On(resources.SEFlameAttack)
		}

		if p.count == 3*resources.SkillFlamePillarDelay {
			// Add damage
			localanim.DamageManager().New(damage.Damage{
				DamageType:    damage.TypePosition,
				Pos:           p.point,
				Power:         int(p.Arg.Power),
				TTL:           7 * resources.SkillFlamePillarDelay,
				TargetObjType: p.Arg.TargetType,
				ShowHitArea:   true,
				BigDamage:     true,
				Element:       damage.ElementFire,
			})
		}

		if p.count > resources.SkillFlamePillarEndCount {
			p.count = 0
			p.state = resources.SkillFlamePillarStateDoing
			return false, nil
		}
	case resources.SkillFlamePillarStateDoing:
		num := 3
		if p.count > num*resources.SkillFlamePillarDelay {
			p.count = 0
			p.state = resources.SkillFlamePillarStateEnd
			return false, nil
		}
	case resources.SkillFlamePillarStateEnd:
		if p.count > resources.SkillFlamePillarEndCount {
			return true, nil
		}
	case resources.SkillFlamePillarStateDeleted:
		// Nothing to do
	}

	p.count++
	return false, nil
}
