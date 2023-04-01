package damage

import (
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

var (
	damages = make(map[string]*Damage)
)

func New(dm Damage) string {
	dm.ID = uuid.New().String()
	damages[dm.ID] = &dm
	logger.Debug("Add damage: %+v to damage manager", dm)
	return dm.ID
}

func MgrProcess() {
	for id, d := range damages {
		d.TTL--
		if d.TTL <= 0 {
			delete(damages, id)
		}
	}
}

func Get(pos common.Point) *Damage {
	for _, d := range damages {
		if d.Pos.X == pos.X && d.Pos.Y == pos.Y {
			return d
		}
	}
	return nil
}

func Exists(id string) bool {
	_, ok := damages[id]
	return ok
}

func Remove(id string) {
	delete(damages, id)
}

func RemoveAll() {
	damages = make(map[string]*Damage)
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
