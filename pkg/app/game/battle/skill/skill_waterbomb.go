package skill

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

const (
	waterBombEndCount = 60
)

type waterBomb struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	x       int
	y       int
	count   int
	targetX int
	targetY int
}

func newWaterBomb(objID string, arg Argument) *waterBomb {
	px, py := objanim.GetObjPos(arg.OwnerID)
	tx := px + 3
	ty := py
	objType := objanim.ObjTypePlayer
	if arg.TargetType == damage.TargetEnemy {
		objType = objanim.ObjTypeEnemy
	}

	objs := objanim.GetObjs(objanim.Filter{ObjType: objType})
	if len(objs) > 0 {
		tx = objs[0].PosX
		ty = objs[0].PosY
	}

	return &waterBomb{
		ID:         objID,
		OwnerID:    arg.OwnerID,
		Power:      arg.Power,
		TargetType: arg.TargetType,
		targetX:    tx,
		targetY:    ty,
		x:          px,
		y:          py,
	}
}

func (p *waterBomb) Draw() {
	imgNo := (p.count / delayBombThrow) % len(imgBombThrow)
	vx, vy := battlecommon.ViewPos(p.x, p.y)

	// y = ax^2 + bx + c
	// (0,0), (d/2, ymax), (d, 0)
	// y = (4 * ymax / d^2)x^2 + (4 * ymax / d)x
	size := field.PanelSizeX * (p.targetX - p.x)
	ofsx := size * p.count / waterBombEndCount
	ymax := 100
	ofsy := ymax*4*ofsx*ofsx/(size*size) - ymax*4*ofsx/size

	if p.targetY != p.y {
		size = field.PanelSizeY * (p.targetY - p.y)
		dy := size * p.count / waterBombEndCount
		ofsy += dy
	}

	dxlib.DrawRotaGraph(vx+int32(ofsx), vy+int32(ofsy), 1, 0, imgBombThrow[imgNo], dxlib.TRUE)
}

func (p *waterBomb) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(sound.SEBombThrow)
	}

	if p.count == waterBombEndCount {
		pn := field.GetPanelInfo(p.targetX, p.targetY)
		if pn.Status == field.PanelStatusHole {
			return true, nil
		}

		sound.On(sound.SEWaterLanding)
		anim.New(effect.Get(effect.TypeWaterBomb, p.targetX, p.targetY, 0))
		damage.New(damage.Damage{
			PosX:          p.targetX,
			PosY:          p.targetY,
			Power:         int(p.Power),
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
			BigDamage:     true,
		})
		field.PanelCrack(p.targetX, p.targetY)
		return true, nil
	}
	return false, nil
}

func (p *waterBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}
