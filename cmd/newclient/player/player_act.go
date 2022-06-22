package player

import (
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/newnetconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
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
	count int
	obj   *object.Object
}

func NewActBuster(obj *object.Object) *ActBuster {
	return &ActBuster{
		count: 0,
		obj:   obj,
	}
}

func (a *ActBuster) Process() bool {
	if a.count == 0 {
		a.obj.UpdateBaseTime = true
		a.obj.Type = object.TypeRockmanBuster
		netconn.GetInst().SendObject(*a.obj)
	}

	// TODO damage

	a.count++
	delay := object.ImageDelays[a.obj.Type]
	num := 6
	return a.count > delay*num
}

func (a *ActBuster) Interval() int {
	return 30
}
