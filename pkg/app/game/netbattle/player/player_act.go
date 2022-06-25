package player

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	netfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/newnetconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
)

type ActOption struct {
	KeepCount  int
	MoveDirect int
	Charged    bool
	ShotPower  int
}

type Act struct {
	Type   int
	Count  int
	Opts   ActOption
	Object *object.Object

	drawMgr *draw.DrawManager
}

func NewAct(drawMgr *draw.DrawManager, obj *object.Object) *Act {
	res := &Act{
		Object:  obj,
		drawMgr: drawMgr,
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
	case battlecommon.PlayerActMove:
		if a.Count == 2 {
			pos := common.Point{X: a.Object.X, Y: a.Object.Y}
			battlecommon.MoveObject(&pos, a.Opts.MoveDirect, appfield.PanelTypePlayer, true, netfield.GetPanelInfo)
			a.Object.X = pos.X
			a.Object.Y = pos.Y
			logger.Debug("Moved to (%d, %d)", a.Object.X, a.Object.Y)
		}
	case battlecommon.PlayerActBuster:
		// TODO damages
	default:
		panic(fmt.Sprintf("Invalid player anim type %d was specified.", a.Type))
	}

	a.Count++
	num, delay := a.drawMgr.GetImageInfo(getObjType(a.Type))
	num += a.Opts.KeepCount
	if a.Count > num*delay {
		// Reset params
		a.Init()
		a.Object.Type = object.TypeRockmanStand
		netconn.GetInst().SendObject(*a.Object)
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
	netconn.GetInst().SendObject(*a.Object)
}

func getObjType(actType int) int {
	switch actType {
	case battlecommon.PlayerActMove:
		return object.TypeRockmanMove
	case battlecommon.PlayerActBuster:
		return object.TypeRockmanBuster
	case battlecommon.PlayerActShot:
		return object.TypeRockmanShot
	case battlecommon.PlayerActBomb:
		return object.TypeRockmanBomb
	case battlecommon.PlayerActCannon:
		return object.TypeRockmanCannon
	case battlecommon.PlayerActDamage:
		return object.TypeRockmanDamage
	case battlecommon.PlayerActPick:
		return object.TypeRockmanPick
	case battlecommon.PlayerActSword:
		return object.TypeRockmanSword
	}

	panic(fmt.Sprintf("Undefined object type for act %d", actType))
}
