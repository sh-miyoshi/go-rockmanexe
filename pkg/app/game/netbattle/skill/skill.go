package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
	"github.com/stretchr/stew/slice"
)

type Argument struct {
	X     int
	Y     int
	Power int
}

type Skill interface {
	Process() (bool, error)
	RemoveObject()
	StopByPlayer()
}

var (
	skills   = make(map[string]Skill)
	playerID string
)

func Init(plyrID string) {
	skills = make(map[string]Skill)
	playerID = plyrID
}

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
	case skill.SkillVulcan1:
		skills[id] = newVulcan(arg.X, arg.Y, 3)
	case skill.SkillWideShot:
		skills[id] = newWideShot(arg.X, arg.Y, arg.Power, 8)
	case skill.SkillSpreadGun:
		skills[id] = newSpreadGun(arg.X, arg.Y, arg.Power)
	case skill.SkillPlayerShockWave, skill.SkillShockWave:
		skills[id] = newShockWave(arg.X, arg.Y, arg.Power)
	case skill.SkillThunderBall:
		skills[id] = newThunderBall(arg.X, arg.Y, arg.Power)
	case skill.SkillRecover:
		skills[id] = newRecover(arg.X, arg.Y, arg.Power)
	case skill.SkillMiniBomb:
		skills[id] = newMiniBomb(arg.X, arg.Y, arg.Power)
	default:
		panic(fmt.Sprintf("Invalid skill id: %d", skillID))
	}

	return id
}

func StopByPlayer(id string) {
	if s, ok := skills[id]; ok {
		s.StopByPlayer()
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

	myClientID := config.Get().Net.ClientID
	finfo := netconn.GetFieldInfo()
	for _, obj := range finfo.Objects {
		if obj.ClientID != myClientID && slice.Contains(rockmanObj, obj.Type) {
			res = append(res, obj)
		}
	}

	return res
}

func isObjectHit(x, y int) bool {
	finfo := netconn.GetFieldInfo()
	for _, obj := range finfo.Objects {
		if obj.Hittable && obj.X == x && obj.Y == y && obj.ID != playerID {
			return true
		}
	}
	return false
}
