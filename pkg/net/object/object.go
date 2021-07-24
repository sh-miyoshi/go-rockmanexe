package object

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
)

const (
	TypeRockmanStand int = iota
	TypeRockmanMove
	TypeRockmanDamage
	TypeRockmanShot
	TypeRockmanCannon
	TypeRockmanSword
	TypeRockmanBomb
	TypeRockmanBuster
	TypeRockmanPick

	TypeNormalCannonAtk
	TypeNormalCannonBody
	TypeHighCannonAtk
	TypeHighCannonBody
	TypeMegaCannonAtk
	TypeMegaCannonBody
	TypeSword
	TypeWideSword
	TypeLongSword
	TypeMiniBomb
	TypeRecover
	TypeSpreadGunAtk
	TypeSpreadGunBody
	TypeVulcan
	TypePick
	TypeThunderBall
	TypeWideShotBody
	TypeWideShotBegin
	TypeWideShotMove
	TypeShockWave

	TypeMax
)

var (
	ImageDelays = [TypeMax]int{
		1, 1, 2, 2, 6, 3, 4, 1, 4, // Rockman
		2, 6, 2, 6, 2, 6, 3, 3, 3, 4, 1, 2, 2, 2, 3, 6, 4, 4, 4, 3, // Skills
	}
)

type Object struct {
	ClientID       string
	ID             string
	Type           int
	HP             int
	X              int
	Y              int
	Chips          []int
	BaseTime       time.Time
	UpdateBaseTime bool
	ViewOfsX       int32
	ViewOfsY       int32
	DamageChecked  bool
	HitDamage      damage.Damage
	Count          int
	Hittable       bool
	Speed          int
	TargetX        int
	TargetY        int
}

func Marshal(obj Object) []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(&obj)
	return buf.Bytes()
}

func Unmarshal(obj *Object, data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(obj)
}
