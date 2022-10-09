package damage

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
)

const (
	TargetPlayer int = 1 << iota
	TargetEnemy
)

const (
	DamageTypeNone int = iota
	DamageTypeFire
	DamageTypeWater
	DamageTypeElec
	DamageTypeWood
)

type Damage struct {
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
