package skill

import (
	"fmt"

	skillanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/skill"
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
	skillanim.Anim

	StopByOwner()
}

func Get(id int, arg Argument) SkillAnim {
	changePanelStatus := func(pos point.Point, crackType int, endCount int) {
		arg.FieldFuncs.ChangePanelStatus(arg.OwnerClientID, pos, crackType, endCount)
	}
	changePanelType := func(pos point.Point, pnType int, endCount int) {
		arg.FieldFuncs.ChangePanelType(arg.OwnerClientID, pos, pnType, endCount)
	}

	coreArg := skillcore.Argument{
		OwnerID:       arg.OwnerObjectID,
		OwnerClientID: arg.OwnerClientID,
		Power:         arg.Power,
		TargetType:    arg.TargetType,

		DamageMgr:         arg.Manager.DamageMgr(),
		GetObjectPos:      arg.Manager.ObjAnimGetObjPos,
		SoundOn:           arg.Manager.SoundOn,
		GetObjects:        arg.Manager.ObjAnimGetObjs,
		GetPanelInfo:      arg.FieldFuncs.GetPanelInfo,
		ChangePanelStatus: changePanelStatus,
		Cutin: func(skillName string, count int) {
			arg.Manager.Cutin(skillName, count, arg.OwnerClientID, arg.AnimObjID)
		},
		ChangePanelType: changePanelType,
		MakeInvisible:   arg.Manager.ObjAnimMakeInvisible,
		AddBarrier:      arg.Manager.ObjAnimAddBarrier,
		SetCustomGaugeMax: func(objID string) {
			// WIP
		},
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
	case resources.SkillVulcan1, resources.SkillVulcan2, resources.SkillVulcan3:
		return newVulcan(arg, core)
	case resources.SkillPlayerWideShot:
		return newWideShot(arg, core)
	case resources.SkillHeatShot, resources.SkillHeatV, resources.SkillHeatSide:
		return newHeatShot(arg, core)
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
	case resources.SkillAreaSteal:
		return newAreaSteal(arg, core)
	case resources.SkillBubbleShot, resources.SkillBubbleSide, resources.SkillBubbleV:
		return newBubbleShot(arg, core)
	default:
		system.SetError(fmt.Sprintf("skill %d is not implemented yet", id))
		return nil
	}
}
