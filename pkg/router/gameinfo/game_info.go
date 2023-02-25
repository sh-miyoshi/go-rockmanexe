package gameinfo

import (
	"bytes"
	"encoding/gob"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
)

type PanelInfo struct {
	OwnerClientID string
	ObjectID      string
}

type Object struct {
	ID            string
	OwnerClientID string
	HP            int
	Pos           common.Point
	// TODO(他にも必要だと思うが都度追加していく)
}

type Anim struct {
	ObjectID      string
	OwnerClientID string
	Pos           common.Point
	AnimType      int
}

type GameInfo struct {
	Panels          [config.FieldNumX][config.FieldNumY]PanelInfo
	Objects         []Object
	Anims           []Anim
	ReverseClientID string
}

func (p *GameInfo) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *GameInfo) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}

func (p *GameInfo) GetObject(id string) *Object {
	for i, obj := range p.Objects {
		if obj.ID == id {
			return &p.Objects[i]
		}
	}
	return nil
}
