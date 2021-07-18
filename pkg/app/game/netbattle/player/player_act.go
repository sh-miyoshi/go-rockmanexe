package player

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	netfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
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

	Object *object.Object
}

func NewAct(obj *object.Object) *Act {
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
	case battlecommon.PlayerActMove:
		if a.Count == 2 {
			battlecommon.MoveObject(&a.Object.X, &a.Object.Y, a.Opts.MoveDirect, appfield.PanelTypePlayer, true, netfield.GetPanelInfo)
			logger.Debug("Moved to (%d, %d)", a.Object.X, a.Object.Y)
		}
	case battlecommon.PlayerActBuster:
		// Add buster damage
		if a.Count == 1 {
			dm := []damage.Damage{}
			y := a.Object.Y
			for x := a.Object.X + 1; x < appfield.FieldNumX; x++ {
				dm = append(dm, damage.Damage{
					ID:            uuid.New().String(),
					ClientID:      a.Object.ClientID,
					PosX:          x,
					PosY:          y,
					Power:         1, // TODO change power
					TTL:           1,
					TargetType:    damage.TargetOtherClient,
					HitEffectType: effect.TypeHitSmallEffect,
					ViewOfsX:      int32(rand.Intn(2*5) - 5),
					ViewOfsY:      int32(rand.Intn(2*5) - 5),
				})

				// break if object exists
				pn := netfield.GetPanelInfo(x, y)
				if pn.ObjectID != "" {
					break
				}
			}
			netconn.SendDamages(dm)
		}
	case battlecommon.PlayerActCannon, battlecommon.PlayerActSword, battlecommon.PlayerActBomb, battlecommon.PlayerActDamage, battlecommon.PlayerActShot, battlecommon.PlayerActPick:
		// No special action
	default:
		panic(fmt.Sprintf("Invalid player anim type %d was specified.", a.Type))
	}

	a.Count++
	num, delay := draw.GetImageInfo(getObjType(a.Type))
	num += a.Opts.KeepCount
	if a.Count > num*delay {
		// Reset params
		a.Init()
		a.Object.Type = object.TypeRockmanStand
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
