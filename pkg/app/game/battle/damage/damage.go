package damage

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	TargetPlayer int = 1 << iota
	TargetEnemy
)

const (
	TypeNone int = iota
	TypeFire
	TypeWater
	TypeElec
	TypeWood
)

type Damage struct {
	OwnerClientID string
	ID            string
	Pos           common.Point
	Power         int
	TTL           int
	TargetType    int
	HitEffectType int
	ShowHitArea   bool
	BigDamage     bool
	PushRight     int
	PushLeft      int
	DamageType    int
	// TODO: のけぞり(単体), インビジ貫通
}

type DamageManager struct {
	damages map[string]*Damage
}

func IsWeakness(charType int, dm Damage) bool {
	switch charType {
	case TypeFire:
		return dm.DamageType == TypeWater
	case TypeWater:
		return dm.DamageType == TypeElec
	case TypeElec:
		return dm.DamageType == TypeWood
	case TypeWood:
		return dm.DamageType == TypeFire
	}
	return false
}

func NewManager() *DamageManager {
	return &DamageManager{
		damages: make(map[string]*Damage),
	}
}

func (m *DamageManager) New(dm Damage) string {
	dm.ID = uuid.New().String()
	m.damages[dm.ID] = &dm
	logger.Debug("Add damage: %+v to damage manager", dm)
	return dm.ID
}

func (m *DamageManager) Process() {
	for id, d := range m.damages {
		d.TTL--
		if d.TTL <= 0 {
			delete(m.damages, id)
		}
	}
}

func (m *DamageManager) Get(pos common.Point) *Damage {
	for _, d := range m.damages {
		if d.Pos.Equal(pos) {
			return d
		}
	}
	return nil
}

func (m *DamageManager) Exists(id string) bool {
	_, ok := m.damages[id]
	return ok
}

func (m *DamageManager) Remove(id string) {
	delete(m.damages, id)
}

func (m *DamageManager) RemoveAll() {
	m.damages = make(map[string]*Damage)
}

func (m *DamageManager) PrintAllData() {
	msg := "current damages: { "
	for id, d := range m.damages {
		msg += fmt.Sprintf("%s: %+v, ", id, *d)
	}
	msg += " }"
	logger.Debug(msg)
}
