package sysinfo

import (
	"bytes"
	"encoding/gob"
)

type SystemType int

const (
	TypeCutin SystemType = iota
	TypeActing
)

type SysInfo struct {
	Type SystemType
	Data []byte
}

type Cutin struct {
	Count     int
	SkillName string
}

func (p *Cutin) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *Cutin) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
