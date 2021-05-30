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

	ObjectTypeMax
)

var (
	ImageDelays = [ObjectTypeMax]int{1, 1, 2, 2, 6, 3, 4, 1, 4}
)

type Object struct {
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

type Info struct {
	CurrentTime time.Time
	MyArea      [3][3]Object
	EnemyArea   [3][3]Object
}

func Marshal(fieldInfo *Info) []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(fieldInfo)
	return buf.Bytes()
}

func Unmarshal(fieldInfo *Info, data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(fieldInfo)
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

func (i *Info) Init() {
}
