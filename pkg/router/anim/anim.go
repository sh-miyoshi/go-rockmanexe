package anim

import (
	"bytes"
	"encoding/gob"
)

const (
	TypeCannonNormal int = iota
	TypeCannonHigh
	TypeCannonMega
	TypeMiniBomb
	TypeRecover
	TypeShockWave
	TypeSpreadGun
	TypeSpreadHit
	TypeSword
	TypeWideSword
	TypeLongSword
	TypeVulcan
	TypeWideShot

	TypeMax
)

type NetInfo struct {
	AnimType      int
	ActCount      int
	OwnerClientID string
}

func (p *NetInfo) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *NetInfo) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
