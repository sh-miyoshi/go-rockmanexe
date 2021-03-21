package damage

const (
	TargetPlayer int = 1 << iota
	TargetEnemy
)

type Damage struct {
	PosX          int
	PosY          int
	Power         int
	TTL           int
	TargetType    int
	HitEffectType int
	ShowHitArea   bool
	// TODO: のけぞり, インビジ貫通
}

var (
	damages []Damage
)

func New(dm Damage) {
	damages = append(damages, dm)
}

func MgrProcess() {
	newList := []Damage{}
	for _, d := range damages {
		d.TTL--
		if d.TTL > 0 {
			newList = append(newList, d)
		}
	}
	damages = newList
}

func Get(x, y int) *Damage {
	for _, d := range damages {
		if d.PosX == x && d.PosY == y {
			return &d
		}
	}
	return nil
}
