package anim

import (
	"bytes"
	"encoding/gob"
	"time"
)

const (
	TypeCannonNormal int = iota
	TypeCannonHigh
	TypeCannonMega
	TypeMiniBomb

	TypeMax
)

type NetInfo struct {
	AnimType      int
	StartedAt     time.Time
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
