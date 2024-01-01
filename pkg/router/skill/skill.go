package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
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

func GetByChip(chipID int, arg Argument) SkillAnim {
	switch chipID {
	case chip.IDCannon, chip.IDHighCannon, chip.IDMegaCannon:
		a := skillcore.Argument{
			Power:        arg.Power,
			OwnerID:      arg.OwnerObjectID,
			TargetType:   arg.TargetType,
			GetPanelInfo: arg.GameInfo.GetPanelInfo,
		}
		core := routeranim.SKillManager(arg.OwnerClientID).Get(skillcore.SkillCannon, a)
		switch chipID {
		case chip.IDCannon:
			return newCannon(TypeNormalCannon, arg, core)
		case chip.IDHighCannon:
			return newCannon(TypeHighCannon, arg, core)
		case chip.IDMegaCannon:
			return newCannon(TypeMegaCannon, arg, core)
		}
	case chip.IDMiniBomb:
		return newMiniBomb(arg)
	case chip.IDRecover10, chip.IDRecover30:
		return newRecover(arg)
	case chip.IDShockWave:
		return newShockWave(arg)
	case chip.IDSpreadGun:
		return newSpreadGun(arg)
	case chip.IDSword:
		return newSword(resources.SkillTypeSword, arg)
	case chip.IDWideSword:
		return newSword(resources.SkillTypeWideSword, arg)
	case chip.IDLongSword:
		return newSword(resources.SkillTypeLongSword, arg)
	case chip.IDVulcan1:
		return newVulcan(3, arg)
	case chip.IDWideShot:
		return newWideShot(arg)
	}
	panic(fmt.Sprintf("chip %d is not implemented yet", chipID))
}
