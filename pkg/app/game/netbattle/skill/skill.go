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
	StopByPlayer()
}

type SkillManager struct {
	skills map[string]Skill
}

var (
	inst SkillManager
)

func GetInst() *SkillManager {
	return &inst
}

func (m *SkillManager) Init() {
	m.skills = make(map[string]Skill)
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

func (m *SkillManager) Add(skillType int, args Argument) string {
	id := uuid.New().String()

	switch skillType {
	case skill.SkillCannon:
		m.skills[id] = newCannon(args.X, args.Y, args.Power, skill.TypeNormalCannon)
	case skill.SkillMiniBomb:
		m.skills[id] = newMiniBomb(args.X, args.Y, args.Power)
	case skill.SkillRecover:
		m.skills[id] = newRecover(args.X, args.Y, args.Power)
	default:
		panic(fmt.Sprintf("skill %d is not implemented yet", skillType))
	}

	return id
}

func (m *SkillManager) StopByPlayer(id string) {
	if s, ok := m.skills[id]; ok {
		s.StopByPlayer()
	}
}
