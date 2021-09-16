package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
	"github.com/stretchr/stew/slice"
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
	skills   map[string]Skill
	clientID string
)

func Init(cid string) {
	skills = make(map[string]Skill)
	clientID = cid
}

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
	case skill.SkillVulcan1:
		skills[id] = newVulcan(arg.X, arg.Y, 3)
	case skill.SkillMiniBomb:
		skills[id] = newMiniBomb(arg.X, arg.Y, 50)
	case skill.SkillThunderBall:
		skills[id] = newThunderBall(arg.X, arg.Y, 30)
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

func getEnemies() []object.Object {
	res := []object.Object{}
	rockmanObj := []int{
		object.TypeRockmanStand,
		object.TypeRockmanMove,
		object.TypeRockmanDamage,
		object.TypeRockmanShot,
		object.TypeRockmanCannon,
		object.TypeRockmanSword,
		object.TypeRockmanBomb,
		object.TypeRockmanBuster,
		object.TypeRockmanPick,
	}

	myClientID := clientID
	finfo := netconn.GetFieldInfo()
	for _, obj := range finfo.Objects {
		if obj.ClientID != myClientID && slice.Contains(rockmanObj, obj.Type) {
			res = append(res, obj)
		}
	}

	return res
}
