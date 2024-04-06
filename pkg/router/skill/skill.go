package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type Argument struct {
	AnimObjID     string
	OwnerObjectID string
	OwnerClientID string
	Power         uint
	TargetType    int
	Manager       *manager.Manager
	FieldFuncs    gameinfo.FieldFuncs
}

type SkillAnim interface {
	anim.Anim

	StopByOwner()
}

func Get(id int, arg Argument) SkillAnim {
	panelBreak := func(pos point.Point) {
		arg.FieldFuncs.PanelBreak(arg.OwnerClientID, pos)
	}

	coreArg := skillcore.Argument{
		OwnerID:       arg.OwnerObjectID,
		OwnerClientID: arg.OwnerClientID,
		Power:         arg.Power,
		TargetType:    arg.TargetType,

		DamageMgr:    arg.Manager.DamageMgr(),
		GetObjectPos: arg.Manager.ObjAnimGetObjPos,
		SoundOn:      arg.Manager.SoundOn,
		GetObjects:   arg.Manager.ObjAnimGetObjs,
		GetPanelInfo: arg.FieldFuncs.GetPanelInfo,
		PanelBreak:   panelBreak,
		Cutin: func(skillName string, count int) {
			// TODO
		},
		ChangePanelType: func(pos point.Point, pnType int) {
			// TODO
		},
		MakeInvisible: arg.Manager.ObjAnimMakeInvisible,
	}
	core := arg.Manager.SkillGet(id, coreArg)

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
	case resources.SkillTornado:
		return newTornado(arg, core)
	case resources.SkillBoomerang:
		return newBoomerang(arg, core)
	case resources.SkillBambooLance:
		return newBambooLance(arg, core)
	case resources.SkillCrackout, resources.SkillDoubleCrack, resources.SkillTripleCrack:
		return newCrack(arg, core)
	default:
		system.SetError(fmt.Sprintf("skill %d is not implemented yet", id))
		return nil
	}
}
