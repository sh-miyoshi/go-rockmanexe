package player

import (
	"fmt"

	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	netfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

const (
	actTypeMove int = iota
	actTypeBuster
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
	switch a.Type {
	case -1:
		return false
	case actTypeMove:
		if a.Count == 2 {
			battlecommon.MoveObject(&a.Object.X, &a.Object.Y, a.Opts.MoveDirect, appfield.PanelTypePlayer, true, netfield.GetPanelInfo)
			logger.Debug("Moved to (%d, %d)", a.Object.X, a.Object.Y)
		}
	case actTypeBuster:
		// TODO add buster damage
	default:
		panic(fmt.Sprintf("Invalid player anim type %d was specified.", a.Type))
	}

	a.Count++
	num, delay := draw.GetImageInfo(getObjType(a.Type))
	num += a.Opts.KeepCount
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
	case actTypeMove:
		return field.ObjectTypeRockmanMove
	case actTypeBuster:
		return field.ObjectTypeRockmanBuster
	}

	panic(fmt.Sprintf("Undefined object type for act %d", actType))
}
