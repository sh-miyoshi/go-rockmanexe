package field

import (
	"bytes"
	"encoding/gob"
	"time"
)

const (
	ObjectTypeRockmanStand int = iota
	ObjectTypeRockmanMove
	ObjectTypeRockmanDamage
	ObjectTypeRockmanShot
	ObjectTypeRockmanCannon
	ObjectTypeRockmanSword
	ObjectTypeRockmanBomb
	ObjectTypeRockmanBuster
	ObjectTypeRockmanPick

	ObjectTypeCannonAtk
	ObjectTypeCannonBody
	ObjectTypeSword
	ObjectTypeMiniBomb
	ObjectTypeRecover
	ObjectTypeSpreadGunAtk
	ObjectTypeSpreadGunBody
	ObjectTypeVulcan
	ObjectTypePick
	ObjectTypeThunderBall
	ObjectTypeWideShot
	ObjectTypeShockWave

	ObjectTypeMax
)

var (
	ImageDelays = [ObjectTypeMax]int{
		1, 1, 2, 2, 6, 3, 4, 1, 4, // Rockman
		2, 6, 3, 4, 1, 2, 2, 3, 6, 4, 5, // Skills
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
	// TODO ...
}

func MarshalObject(obj Object) []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(&obj)
	return buf.Bytes()
}

func UnmarshalObject(obj *Object, data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(obj)
}
