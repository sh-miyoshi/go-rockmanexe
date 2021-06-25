package app

import (
	"fmt"
	"math/rand"

	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/netconn"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

type ActOption struct {
	KeepCount  int
	MoveDirect int
	Charged    bool
	ShotPower  int
}

type Act struct {
	Type  int
	Count int
	Opts  ActOption

	Object *field.Object
}

func NewAct(obj *field.Object) *Act {
	res := &Act{
		Object: obj,
	}
	res.Init()
	return res
}

func (a *Act) Init() {
	a.Type = -1
	a.Count = 0
	a.Opts = ActOption{}
}

// Process method returns true if processing now
func (a *Act) Process() bool {
	var num int

	switch a.Type {
	case -1:
		return false
	case battlecommon.PlayerActMove:
		num = 4

		if a.Count == 2 {
			a.Object.X = rand.Intn(3)
			a.Object.Y = rand.Intn(3)
		}
	case battlecommon.PlayerActCannon:
		num = 6
	default:
		panic(fmt.Sprintf("Invalid player anim type %d was specified.", a.Type))
	}

	a.Count++
	num += a.Opts.KeepCount
	delay := field.ImageDelays[getObjType(a.Type)]
	if a.Count > num*delay {
		// Reset params
		a.Init()
		a.Object.Type = field.ObjectTypeRockmanStand
		netconn.SendObject(*a.Object)
		return false // finished
	}
	return true // processing now
}

func (a *Act) Set(actType int, opts *ActOption) {
	a.Type = actType
	if opts != nil {
		a.Opts = *opts
	}

	a.Object.UpdateBaseTime = true
	a.Object.Type = getObjType(actType)
	netconn.SendObject(*a.Object)
}

func getObjType(actType int) int {
	switch actType {
	case battlecommon.PlayerActMove:
		return field.ObjectTypeRockmanMove
	case battlecommon.PlayerActBuster:
		return field.ObjectTypeRockmanBuster
	case battlecommon.PlayerActShot:
		return field.ObjectTypeRockmanShot
	case battlecommon.PlayerActBomb:
		return field.ObjectTypeRockmanBomb
	case battlecommon.PlayerActCannon:
		return field.ObjectTypeRockmanCannon
	case battlecommon.PlayerActDamage:
		return field.ObjectTypeRockmanDamage
	case battlecommon.PlayerActPick:
		return field.ObjectTypeRockmanPick
	case battlecommon.PlayerActSword:
		return field.ObjectTypeRockmanSword
	}

	panic(fmt.Sprintf("Undefined object type for act %d", actType))
}
