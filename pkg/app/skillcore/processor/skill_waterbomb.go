package processor

import (
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	waterBombEndCount = 60
)

type WaterBomb struct {
	Arg skillcore.Argument

	count  int
	pos    point.Point
	target point.Point
	hits   []point.Point
}

func (p *WaterBomb) Init() {
	p.pos = p.Arg.GetObjectPos(p.Arg.OwnerID)
	p.target = point.Point{X: p.pos.X + 3, Y: p.pos.Y}
	objType := objanim.ObjTypePlayer
	if p.Arg.TargetType == damage.TargetEnemy {
		objType = objanim.ObjTypeEnemy
	}

	objs := p.Arg.GetObjects(objanim.Filter{ObjType: objType})
	if len(objs) > 0 {
		p.target = objs[0].Pos
	}
}

func (p *WaterBomb) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		p.Arg.SoundOn(resources.SEBombThrow)
	}

	if p.count == waterBombEndCount {
		pn := p.Arg.GetPanelInfo(p.target)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		p.Arg.SoundOn(resources.SEWaterLanding)
		p.hits = append(p.hits, p.target)
		if objID := p.Arg.GetPanelInfo(p.target).ObjectID; objID != "" {
			p.Arg.DamageMgr.New(damage.Damage{
				DamageType:    damage.TypeObject,
				Power:         int(p.Arg.Power),
				TargetObjType: p.Arg.TargetType,
				HitEffectType: resources.EffectTypeNone,
				BigDamage:     true,
				Element:       damage.ElementWater,
				TargetObjID:   objID,
			})
		}
		return true, nil
	}
	return false, nil
}

func (p *WaterBomb) GetCount() int {
	return p.count
}

func (p *WaterBomb) GetEndCount() int {
	return waterBombEndCount
}

func (p *WaterBomb) GetPointParams() (current, target point.Point) {
	return p.pos, p.target
}

func (p *WaterBomb) PopHits() []point.Point {
	if len(p.hits) > 0 {
		res := append([]point.Point{}, p.hits...)
		p.hits = []point.Point{}
		return res
	}
	return []point.Point{}
}
