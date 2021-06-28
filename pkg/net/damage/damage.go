package damage

import (
	"bytes"
	"encoding/gob"
)

type Damage struct {
	PosX       int
	PosY       int
	Power      int
	TTL        int
	TargetType int
}

func Marshal(dm []Damage) []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(&dm)
	return buf.Bytes()
}

func Unmarshal(dm *[]Damage, data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(dm)
}
