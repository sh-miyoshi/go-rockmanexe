package object

import (
	"bytes"
	"encoding/gob"
)

type InitParam struct {
	ID string
	HP int
	X  int
	Y  int
}

func (p *InitParam) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *InitParam) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
