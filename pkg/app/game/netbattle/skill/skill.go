package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
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

type SkillManager struct {
	skills      map[string]Skill
	playerObjID string
}

var (
	inst SkillManager
)

func GetInst() *SkillManager {
	return &inst
}

func (m *SkillManager) Init(playerObjID string) {
	m.skills = make(map[string]Skill)
	m.playerObjID = playerObjID
}

func (m *SkillManager) Process() error {
	for id, s := range m.skills {
		end, err := s.Process()
		if err != nil {
			return fmt.Errorf("skill %s process failed: %w", id, err)
		}

		if end {
			m.skills[id].RemoveObject()
			delete(m.skills, id)
		}
	}
	return nil
}

func (m *SkillManager) Exists(id string) bool {
	_, ok := m.skills[id]
	return ok
}

func (m *SkillManager) Add(skillType int, args Argument) string {
	id := uuid.New().String()

	switch skillType {
	case skill.SkillCannon:
		m.skills[id] = newCannon(args.X, args.Y, args.Power, skill.TypeNormalCannon)
	case skill.SkillHighCannon:
		m.skills[id] = newCannon(args.X, args.Y, args.Power, skill.TypeHighCannon)
	case skill.SkillMegaCannon:
		m.skills[id] = newCannon(args.X, args.Y, args.Power, skill.TypeMegaCannon)
	case skill.SkillSword:
		m.skills[id] = newSword(args.X, args.Y, args.Power, skill.TypeSword)
	case skill.SkillWideSword:
		m.skills[id] = newSword(args.X, args.Y, args.Power, skill.TypeWideSword)
	case skill.SkillLongSword:
		m.skills[id] = newSword(args.X, args.Y, args.Power, skill.TypeLongSword)
	case skill.SkillVulcan1:
		m.skills[id] = newVulcan(args.X, args.Y, 3)
	case skill.SkillWideShot:
		m.skills[id] = newWideShot(args.X, args.Y, args.Power, 8)
	case skill.SkillSpreadGun:
		m.skills[id] = newSpreadGun(args.X, args.Y, args.Power)
	case skill.SkillPlayerShockWave, skill.SkillShockWave:
		m.skills[id] = newShockWave(args.X, args.Y, args.Power)
	case skill.SkillThunderBall:
		m.skills[id] = newThunderBall(args.X, args.Y, args.Power)
	case skill.SkillRecover:
		m.skills[id] = newRecover(args.X, args.Y, args.Power)
	case skill.SkillMiniBomb:
		m.skills[id] = newMiniBomb(args.X, args.Y, args.Power)
	default:
		common.SetError(fmt.Sprintf("skill %d is not implemented yet", skillType))
	}

	return id
}

func (m *SkillManager) StopByPlayer(id string) {
	if s, ok := m.skills[id]; ok {
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
	ginfo := net.GetInst().GetGameInfo()
	for _, obj := range ginfo.Objects {
		if obj.ClientID != myClientID && slice.Contains(rockmanObj, obj.Type) {
			res = append(res, obj)
		}
	}

	return res
}

func isObjectHit(x, y int) bool {
	ginfo := net.GetInst().GetGameInfo()
	for _, obj := range ginfo.Objects {
		if obj.Hittable && obj.X == x && obj.Y == y && obj.ID != inst.playerObjID {
			return true
		}
	}
	return false
}
