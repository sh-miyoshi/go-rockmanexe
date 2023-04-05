package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
)

type Argument struct {
	AnimObjID  string
	OwnerID    string
	Power      uint
	TargetType int

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
	}
	return nil
}
