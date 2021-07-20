package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type Argument struct {
	X int
	Y int
}

type Skill interface {
	Process() (bool, error)
	GetObjects() []object.Object
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
			removeObjects(id)
			delete(skills, id)
		}
	}
	return nil
}

func Add(skillID int, arg Argument, clientID string) string {
	id := uuid.New().String()

	switch skillID {
	case skill.SkillCannon:
		skills[id] = newCannon(arg.X, arg.Y, clientID)
	default:
		panic(fmt.Sprintf("Invalid skill id: %d", skillID))
	}

	for _, obj := range skills[id].GetObjects() {
		netconn.SendObject(obj)
	}

	return id
}

func removeObjects(id string) {
	for _, obj := range skills[id].GetObjects() {
		netconn.RemoveObject(obj.ID)
	}
}
