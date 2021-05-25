package field

import (
	"bytes"
	"encoding/gob"
)

type Object struct {
	ID   string
	Type int
	// TODO ...
}

type Info struct {
	MyArea    [3][3]Object
	EnemyArea [3][3]Object
}

func Marshal(fieldInfo Info) []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(&fieldInfo)
	return buf.Bytes()
}

func Unmarshal(fieldInfo *Info, data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(fieldInfo)
}

func (i *Info) Init() {
	// TODO
}
