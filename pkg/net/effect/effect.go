package effect

import (
	"bytes"
	"encoding/gob"
)

const (
	TypeNone int = iota
	TypeHitSmallEffect
	TypeHitBigEffect
	TypeExplodeEffect
	TypeCannonHitEffect
	TypeSpreadHitEffect
	TypeVulcanHit1Effect
	TypeVulcanHit2Effect

	TypeMax
)

var (
	Delays = [TypeMax]int{1, 1, 1, 2, 1, 1, 1, 1}
)

type Effect struct {
	ID       string
	ClientID string
	Type     int
	X        int
	Y        int
	ViewOfsX int32
	ViewOfsY int32

	Count int
}

func Marshal(eff Effect) []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(&eff)
	return buf.Bytes()
}

func Unmarshal(eff *Effect, data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(eff)
}
