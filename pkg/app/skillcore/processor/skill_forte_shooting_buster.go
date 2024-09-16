package processor

import (
	"math/rand"

	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	forteShootingBusterEndCount = 16
	forteShootingBusterInitWait = 10
)

var (
	forteShootingBusterActQueue []point.Point
)

type ForteShootingBuster struct {
	Arg skillcore.Argument

	count int
	pos   point.Point
}

func (p *ForteShootingBuster) Init() {
	p.count = 0

	objType := objanim.ObjTypePlayer
	panelType := battlecommon.PanelTypePlayer
	if p.Arg.TargetType == damage.TargetEnemy {
		objType = objanim.ObjTypeEnemy
		panelType = battlecommon.PanelTypeEnemy
	}

	if len(forteShootingBusterActQueue) == 0 {
		// 真下に作成する
		objs := p.Arg.GetObjects(objanim.Filter{ObjType: objType})
		if len(objs) > 0 {
			p.pos = objs[0].Pos
		}
	} else {
		attackable := []point.Point{}
		for x := 0; x < battlecommon.FieldNum.X; x++ {
			for y := 0; y < battlecommon.FieldNum.Y; y++ {
				pos := point.Point{X: x, Y: y}
				if p.Arg.GetPanelInfo(pos).Type == panelType {
					ok := true
					for _, p := range forteShootingBusterActQueue {
						if pos.Equal(p) {
							ok = false
							break
						}
					}
					if ok {
						attackable = append(attackable, pos)
					}
				}
			}
		}
		if len(attackable) > 0 {
			n := rand.Intn(len(attackable))
			p.pos = attackable[n]
		}
	}

	forteShootingBusterActQueue = append(forteShootingBusterActQueue, p.pos)
}

func (p *ForteShootingBuster) Process() (bool, error) {
	p.count++
	if p.count == forteShootingBusterInitWait {
		if objID := p.Arg.GetPanelInfo(p.pos).ObjectID; objID != "" {
			p.Arg.DamageMgr.New(damage.Damage{
				OwnerClientID: p.Arg.OwnerClientID,
				TargetObjID:   objID,
				DamageType:    damage.TypeObject,
				Power:         int(p.Arg.Power),
				TargetObjType: p.Arg.TargetType,
				HitEffectType: resources.EffectTypeNone,
				BigDamage:     true,
				Element:       damage.ElementNone,
			})
		}
	}
	if p.count >= forteShootingBusterEndCount {
		// Remove from act queue
		for i, pos := range forteShootingBusterActQueue {
			if p.pos.Equal(pos) {
				forteShootingBusterActQueue = append(forteShootingBusterActQueue[:i], forteShootingBusterActQueue[i+1:]...)
				break
			}
		}
		return true, nil
	}

	return false, nil
}

func (p *ForteShootingBuster) GetCount() int {
	return p.count
}

func (p *ForteShootingBuster) GetPos() point.Point {
	return p.pos
}

func (p *ForteShootingBuster) GetInitWait() int {
	return forteShootingBusterInitWait
}
