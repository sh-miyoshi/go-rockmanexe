package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
)

type Argument struct {
	AnimObjID     string
	OwnerObjectID string
	OwnerClientID string
	Power         uint
	TargetType    int

	GameInfo *gameinfo.GameInfo
	QueueIDs []string
}

type SkillAnim interface {
	anim.Anim

	StopByOwner()
	GetEndCount() int
}

func Get(id int, arg Argument) SkillAnim {
	coreArg := skillcore.Argument{
		OwnerID:       arg.OwnerObjectID,
		OwnerClientID: arg.OwnerClientID,
		Power:         arg.Power,
		TargetType:    arg.TargetType,
		GetPanelInfo:  arg.GameInfo.GetPanelInfo,
	}
	core := routeranim.SkillManager(arg.OwnerClientID).Get(id, coreArg)

	switch id {
	case resources.SkillCannon, resources.SkillHighCannon, resources.SkillMegaCannon:
		return newCannon(arg, core)
	case resources.SkillMiniBomb:
		return newMiniBomb(arg, core)
	case resources.SkillRecover:
		return newRecover(arg, core)
	case resources.SkillEnemyShockWave:
		return newShockWave(arg, core)
	case resources.SkillSpreadGun:
		return newSpreadGun(arg, core)
	case resources.SkillSword, resources.SkillWideSword, resources.SkillLongSword:
		return newSword(id, arg, core)
	case resources.SkillVulcan1:
		return newVulcan(arg, core)
	case resources.SkillPlayerWideShot:
		return newWideShot(arg, core)
	case resources.SkillHeatShot, resources.SkillHeatV, resources.SkillHeatSide:
		return newHeatShot(id, arg, core)
	case resources.SkillFlamePillarLine:
		return newFlameLine(arg, core)
	case resources.SkillPlayerShockWave:
		return newShockWave(arg, core)
	default:
		panic(fmt.Sprintf("skill %d is not implemented yet", id))
	}
}
