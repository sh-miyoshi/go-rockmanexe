package processor

import (
	"fmt"
	"sort"

	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	flamePillarEndCount = 20
	flamePillarDelay    = 4
)

const (
	flamePillarTypeTracking int = iota
	flamePillarTypeLine
	// flamePillarTypeRandom // 未実装
)

type FlamePillerParam struct {
	Count int
	State int
	Point point.Point
}

type FlamePillar struct {
	Arg skillcore.Argument

	pm FlamePillerParam
}

type FlamePillarManager struct {
	Arg skillcore.Argument

	count   int
	actType int
	pillars []*FlamePillar
}

func (p *FlamePillarManager) Init(skillID int) {
	p.actType = -1
	isPlayer := p.Arg.TargetType == damage.TargetEnemy

	switch skillID {
	case resources.SkillFlamePillarTracking:
		p.actType = flamePillarTypeTracking

		pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
		if isPlayer {
			pos.X++
		} else {
			pos.X--
		}

		p.pillars = append(p.pillars, &FlamePillar{
			Arg: p.Arg,
			pm: FlamePillerParam{
				State: resources.SkillFlamePillarStateWakeup,
				Point: pos,
			},
		})
	case resources.SkillFlamePillarLine:
		p.actType = flamePillarTypeLine

		pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
		if isPlayer {
			pos.X += 2
		} else {
			pos.X -= 2
		}

		for y := -1; y <= 1; y++ {
			posY := pos.Y + y
			if posY >= 0 && posY < battlecommon.FieldNum.Y {
				p.pillars = append(p.pillars, &FlamePillar{
					Arg: p.Arg,
					pm: FlamePillerParam{
						State: resources.SkillFlamePillarStateWakeup,
						Point: point.Point{X: pos.X, Y: posY},
					},
				})
			}
		}
	}
}

func (p *FlamePillarManager) Process() (bool, error) {
	switch p.actType {
	case flamePillarTypeTracking:
		end, err := p.pillars[0].Process()
		if err != nil {
			return false, fmt.Errorf("flame pillar process failed: %w", err)
		}
		if end {
			// 穴パネルなどで進めなかったら終わり
			if p.pillars[0].pm.State == resources.SkillFlamePillarStateDeleted {
				return true, nil
			}

			var targetObjType int
			x := p.pillars[0].pm.Point.X

			if p.Arg.TargetType == damage.TargetEnemy {
				targetObjType = objanim.ObjTypeEnemy
				x++
				if x >= battlecommon.FieldNum.X {
					return true, nil
				}
			} else {
				targetObjType = objanim.ObjTypePlayer
				x--
				if x < 0 {
					return true, nil
				}
			}

			y := p.pillars[0].pm.Point.Y
			objs := p.Arg.GetObjects(objanim.Filter{ObjType: targetObjType})
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

			p.pillars = append([]*FlamePillar{}, &FlamePillar{
				Arg: p.Arg,
				pm: FlamePillerParam{
					State: resources.SkillFlamePillarStateWakeup,
					Point: point.Point{X: x, Y: y},
				},
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

func (p *FlamePillarManager) GetCount() int {
	return p.count
}

func (p *FlamePillarManager) GetEndCount() int {
	return flamePillarEndCount
}

func (p *FlamePillarManager) GetPillars() []FlamePillerParam {
	res := []FlamePillerParam{}
	for _, pillar := range p.pillars {
		res = append(res, pillar.pm)
	}
	return res
}

func (p *FlamePillarManager) GetDelay() int {
	return flamePillarDelay
}

func (p *FlamePillarManager) IsShowBody() bool {
	if len(p.pillars) > 0 {
		return p.Arg.TargetType == damage.TargetEnemy && p.pillars[0].pm.State != resources.SkillFlamePillarStateEnd
	}
	return false
}

func (p *FlamePillar) Process() (bool, error) {
	switch p.pm.State {
	case resources.SkillFlamePillarStateWakeup:
		if p.pm.Count == 0 {
			pn := p.Arg.GetPanelInfo(p.pm.Point)
			if pn.Status == battlecommon.PanelStatusHole {
				p.pm.State = resources.SkillFlamePillarStateDeleted
				return true, nil
			}

			p.Arg.SoundOn(resources.SEFlameAttack)
		}

		if p.pm.Count == 3*flamePillarDelay {
			// Add damage
			p.Arg.DamageMgr.New(damage.Damage{
				OwnerClientID: p.Arg.OwnerClientID,
				DamageType:    damage.TypePosition,
				Pos:           p.pm.Point,
				Power:         int(p.Arg.Power),
				TTL:           7 * flamePillarDelay,
				TargetObjType: p.Arg.TargetType,
				ShowHitArea:   true,
				BigDamage:     true,
				Element:       damage.ElementFire,
			})
		}

		if p.pm.Count > flamePillarEndCount {
			p.pm.Count = 0
			p.pm.State = resources.SkillFlamePillarStateDoing
			return false, nil
		}
	case resources.SkillFlamePillarStateDoing:
		num := 3
		if p.pm.Count > num*flamePillarDelay {
			p.pm.Count = 0
			p.pm.State = resources.SkillFlamePillarStateEnd
			return false, nil
		}
	case resources.SkillFlamePillarStateEnd:
		if p.pm.Count > flamePillarEndCount {
			return true, nil
		}
	case resources.SkillFlamePillarStateDeleted:
		// Nothing to do
	}

	p.pm.Count++
	return false, nil
}
