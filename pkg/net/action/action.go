package action

import (
	"bytes"
	"encoding/gob"
)

const (
	MoveTypeDirect int = iota
	MoveTypeAbs
)

type Move struct {
	Type    int
	Direct  int
	AbsPosX int
	AbsPosY int
}

type Buster struct {
	Power     int
	IsCharged bool
}

type UseChip struct {
	AnimID           string // anim id must uniq at global
	ChipID           int
	ChipUserClientID string
}

func (p *Move) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *Move) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}

func (p *Buster) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *Buster) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}

func (p *UseChip) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *UseChip) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
