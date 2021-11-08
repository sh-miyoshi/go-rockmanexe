package damage

import "github.com/google/uuid"

const (
	TargetPlayer int = 1 << iota
	TargetEnemy
)

type Damage struct {
	ID            string
	PosX          int
	PosY          int
	Power         int
	TTL           int
	TargetType    int
	HitEffectType int
	ShowHitArea   bool
	BigDamage     bool
	// TODO: のけぞり, インビジ貫通
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

func Get(x, y int) *Damage {
	for _, d := range damages {
		if d.PosX == x && d.PosY == y {
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
