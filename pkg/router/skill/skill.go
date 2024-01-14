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
		Power:        arg.Power,
		OwnerID:      arg.OwnerObjectID,
		TargetType:   arg.TargetType,
		GetPanelInfo: arg.GameInfo.GetPanelInfo,
	}
	core := routeranim.SKillManager(arg.OwnerClientID).Get(id, coreArg)

	switch id {
	case resources.SkillCannon, resources.SkillHighCannon, resources.SkillMegaCannon:
		return newCannon(id, arg, core)
	case resources.SkillMiniBomb:
		return newMiniBomb(arg)
	case resources.SkillRecover:
		return newRecover(arg)
	case resources.SkillShockWave:
		return newShockWave(arg)
	case resources.SkillSpreadGun:
		return newSpreadGun(arg)
	case resources.SkillSword:
		return newSword(resources.SkillTypeSword, arg)
	case resources.SkillWideSword:
		return newSword(resources.SkillTypeWideSword, arg)
	case resources.SkillLongSword:
		return newSword(resources.SkillTypeLongSword, arg)
	case resources.SkillVulcan1:
		return newVulcan(3, arg)
	case resources.SkillWideShot:
		return newWideShot(arg)
	default:
		panic(fmt.Sprintf("skill %d is not implemented yet", id))
	}
}
