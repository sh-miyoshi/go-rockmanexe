package app

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
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

var (
	moveTable = [][]int{
		{1, 2},
		{0, 2},
		{0, 1},
		{0, 0},
		{0, 1},
		{1, 1},
		{2, 1},
		{2, 0},
		{1, 0},
		{1, 1},
		{1, 2},
	}
	moveCount = 0
)

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
	var num int

	switch a.Type {
	case -1:
		return false
	case battlecommon.PlayerActMove:
		num = 4

		if a.Count == 2 {
			// a.Object.X = rand.Intn(3)
			// a.Object.Y = rand.Intn(3)
			a.Object.X = moveTable[moveCount][0]
			a.Object.Y = moveTable[moveCount][1]
			moveCount = (moveCount + 1) % len(moveTable)
		}
	case battlecommon.PlayerActCannon:
		num = 6
	case battlecommon.PlayerActShot:
		num = 6
	case battlecommon.PlayerActBomb:
		num = 5
	case battlecommon.PlayerActSword:
		num = 7
	case battlecommon.PlayerActBuster:
		num = 6

		if a.Count == 1 {
			dm := []damage.Damage{}
			y := a.Object.Y
			for x := a.Object.X + 1; x < 6; x++ {
				dm = append(dm, damage.Damage{
					ID:            uuid.New().String(),
					ClientID:      a.Object.ClientID,
					PosX:          x,
					PosY:          y,
					Power:         1,
					TTL:           1,
					TargetType:    damage.TargetOtherClient,
					HitEffectType: effect.TypeHitSmallEffect,
					ViewOfsX:      int(rand.Intn(2*5) - 5),
					ViewOfsY:      int(rand.Intn(2*5) - 5),
				})
			}
			netconn.SendDamages(dm)
			logger.Info("Add buster damage: %+v", dm)
		}
	default:
		panic(fmt.Sprintf("Invalid player anim type %d was specified.", a.Type))
	}

	a.Count++
	num += a.Opts.KeepCount
	delay := object.ImageDelays[getObjType(a.Type)]
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
