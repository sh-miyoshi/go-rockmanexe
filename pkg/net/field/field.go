package field

import (
	"bytes"
	"encoding/gob"
)

const (
	ObjectTypeRockman int = iota
)

type Object struct {
	ID   string
	Type int
	HP   int
	X    int
	Y    int
	// TODO ...
}

type Info struct {
	MyArea    [3][3]Object
	EnemyArea [3][3]Object
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
