package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
)

type Argument struct {
	OwnerID    string
	Power      uint
	TargetType int

	GameInfo *gameinfo.GameInfo
}

type SkillAnim interface {
	anim.Anim

	StopByOwner()
}

func GetByChip(chipID int, arg Argument) SkillAnim {
	objID := uuid.New().String()

	switch chipID {
	case chip.IDCannon:
		return newCannon(objID, TypeNormalCannon, arg)
	}
	return nil
}
