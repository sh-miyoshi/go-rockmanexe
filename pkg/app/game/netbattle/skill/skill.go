package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

type Skill interface {
	Process() (bool, error)
	GetObjects() []field.Object
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
			delete(skills, id)
		}
	}
	return nil
}

func Add(skillID int) string {
	id := uuid.New().String()

	// debug
	x := 1
	y := 1

	switch skillID {
	case skill.SkillCannon:
		skills[id] = newCannon(id, x, y)
	default:
		panic(fmt.Sprintf("Invalid skill id: %d", skillID))
	}

	for _, obj := range skills[id].GetObjects() {
		netconn.SendObject(obj)
	}

	return id
}
