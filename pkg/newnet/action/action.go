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

func (p *Move) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *Move) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
