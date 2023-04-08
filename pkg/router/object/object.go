package object

import (
	"bytes"
	"encoding/gob"
	"time"
)

const (
	TypePlayerStand int = iota
	TypePlayerBomb
	TypePlayerBuster
	TypePlayerCannon
	TypePlayerDamaged
	TypePlayerMove
	TypePlayerShot
	TypePlayerSword
	TypePlayerThrow
	TypePlayerPick

	TypeMax
)

type NetInfo struct {
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
