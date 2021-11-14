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
	endCount = 60
)

type miniBomb struct {
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

func newMiniBomb(objID string, arg Argument) *miniBomb {
	px, py := objanim.GetObjPos(arg.OwnerID)
	return &miniBomb{
		ID:         objID,
		OwnerID:    arg.OwnerID,
		Power:      arg.Power,
		TargetType: arg.TargetType,
		targetX:    px + 3,
		targetY:    py,
		x:          px,
		y:          py,
	}
}

func (p *miniBomb) Draw() {
	imgNo := (p.count / delayMiniBomb) % len(imgMiniBomb)
	vx, vy := battlecommon.ViewPos(p.x, p.y)

	// y = ax^2 + bx + c
	// (0,0), (d/2, ymax), (d, 0)
	// y = (4 * ymax / d^2)x^2 + (4 * ymax / d)x
	size := field.PanelSizeX * (p.targetX - p.x)
	ofsx := size * p.count / endCount
	ymax := 100
	ofsy := ymax*4*ofsx*ofsx/(size*size) - ymax*4*ofsx/size

	if p.targetY != p.y {
		size = field.PanelSizeY * (p.targetY - p.y)
		dy := size * p.count / endCount
		ofsy += dy
	}

	dxlib.DrawRotaGraph(vx+int32(ofsx), vy+int32(ofsy), 1, 0, imgMiniBomb[imgNo], dxlib.TRUE)
}

func (p *miniBomb) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(sound.SEBombThrow)
	}

	if p.count == endCount {
		// TODO 不発処理(画面外やパネル状況など)
		sound.On(sound.SEExplode)
		anim.New(effect.Get(effect.TypeExplode, p.targetX, p.targetY, 0))
		damage.New(damage.Damage{
			PosX:          p.targetX,
			PosY:          p.targetY,
			Power:         int(p.Power),
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
			BigDamage:     true,
		})
		return true, nil
	}
	return false, nil
}

func (p *miniBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}
