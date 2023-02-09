package object

import (
	"bytes"
	"encoding/gob"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
)

type InitParam struct {
	HP int
	X  int
	Y  int
}

type Object struct {
	ID  string
	HP  int
	Pos common.Point
	// TODO(他にも必要だと思うが都度追加していく)
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
