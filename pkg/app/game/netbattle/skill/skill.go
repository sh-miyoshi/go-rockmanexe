package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
)

type Argument struct {
	X     int
	Y     int
	Power int
}

type Skill interface {
	Process() (bool, error)
	RemoveObject()
}

var (
	skills = make(map[string]Skill)
)

func Process() error {
	for id, s := range skills {
		end, err := s.Process()
		if err != nil {
			return fmt.Errorf("skill %s process failed: %w", id, err)
		}

		if end {
			skills[id].RemoveObject()
			delete(skills, id)
		}
	}
	return nil
}

func Add(skillID int, arg Argument) string {
	id := uuid.New().String()

	switch skillID {
	case skill.SkillCannon:
		skills[id] = newCannon(arg.X, arg.Y, arg.Power, skill.TypeNormalCannon)
	case skill.SkillHighCannon:
		skills[id] = newCannon(arg.X, arg.Y, arg.Power, skill.TypeHighCannon)
	case skill.SkillMegaCannon:
		skills[id] = newCannon(arg.X, arg.Y, arg.Power, skill.TypeMegaCannon)
	case skill.SkillSword:
		skills[id] = newSword(arg.X, arg.Y, arg.Power, skill.TypeSword)
	case skill.SkillWideSword:
		skills[id] = newSword(arg.X, arg.Y, arg.Power, skill.TypeWideSword)
	case skill.SkillLongSword:
		skills[id] = newSword(arg.X, arg.Y, arg.Power, skill.TypeLongSword)
	default:
		panic(fmt.Sprintf("Invalid skill id: %d", skillID))
	}

	return id
}
