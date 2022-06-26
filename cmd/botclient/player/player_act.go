package player

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type Act interface {
	Process() bool
	Interval() int
}

type ActWait struct {
	waitFrame int
}

func NewActWait(waitFrame int) *ActWait {
	return &ActWait{
		waitFrame: waitFrame,
	}
}

func (a *ActWait) Process() bool {
	return true
}

func (a *ActWait) Interval() int {
	return a.waitFrame
}

type ActMove struct {
	targetX int
	targetY int
	obj     *object.Object
}

func NewActMove(obj *object.Object, targetX, targetY int) *ActMove {
	return &ActMove{
		targetX: targetX,
		targetY: targetY,
		obj:     obj,
	}
}

func (a *ActMove) Process() bool {
	logger.Debug("Move to (%d, %d)", a.targetX, a.targetY)
	a.obj.X = a.targetX
	a.obj.Y = a.targetY
	netconn.GetInst().SendObject(*a.obj)
	return true
}

func (a *ActMove) Interval() int {
	return 30
}

type ActBuster struct {
	count     int
	obj       *object.Object
	shotPower uint
	charged   bool
}

func NewActBuster(obj *object.Object) *ActBuster {
	return &ActBuster{
		count: 0,
		obj:   obj,

		// debug
		shotPower: 1,
		charged:   false,
	}
}

func (a *ActBuster) Process() bool {
	if a.count == 0 {
		a.obj.UpdateBaseTime = true
		a.obj.Type = object.TypeRockmanBuster
		netconn.GetInst().SendObject(*a.obj)
	}

	if a.count == 1 {
		s := a.shotPower
		eff := effect.TypeHitSmall
		if a.charged {
			s *= 10
			eff = effect.TypeHitBig
		}

		y := a.obj.Y
		for x := a.obj.X + 1; x < field.FieldNum.X; x++ {
			// logger.Debug("Rock buster damage set %d to (%d, %d)", s, x, *a.pPosY)
			pos := common.Point{X: x, Y: y}
			if field.GetPanelInfo(pos).ObjectID != "" {
				damage.New(damage.Damage{
					Pos:           pos,
					Power:         int(s),
					TTL:           1,
					TargetType:    damage.TargetEnemy,
					HitEffectType: eff,
				})
				break
			}
		}
	}

	a.count++
	delay := object.ImageDelays[a.obj.Type]
	num := 6
	return a.count > delay*num
}

func (a *ActBuster) Interval() int {
	return 30
}
