package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
)

type Argument struct {
	AnimObjID     string
	OwnerObjectID string
	OwnerClientID string
	Power         uint
	TargetType    int

	GameInfo *gameinfo.GameInfo
}

type SkillAnim interface {
	anim.Anim

	StopByOwner()
	GetEndCount() int
}

func GetByChip(chipID int, arg Argument) SkillAnim {
	switch chipID {
	case chip.IDCannon:
		return newCannon(TypeNormalCannon, arg)
	case chip.IDHighCannon:
		return newCannon(TypeHighCannon, arg)
	case chip.IDMegaCannon:
		return newCannon(TypeMegaCannon, arg)
	case chip.IDMiniBomb:
		return newMiniBomb(arg)
	case chip.IDRecover10, chip.IDRecover30:
		return newRecover(arg)
	case chip.IDShockWave:
		return newShockWave(arg)
	case chip.IDSpreadGun:
		return newSpreadGun(arg)
	case chip.IDSword:
		return newSword(TypeSword, arg)
	case chip.IDWideSword:
		return newSword(TypeWideSword, arg)
	case chip.IDLongSword:
		return newSword(TypeLongSword, arg)
	case chip.IDVulcan1:
		return newVulcan(3, arg)
	default:
		panic(fmt.Sprintf("chip %d is not implemented yet", chipID))
	}
}
