package gameinfo

import (
	"bytes"
	"encoding/gob"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
)

type PanelInfo struct {
	OwnerClientID string
}

type GameInfo struct {
	Panels  [config.FieldNumX][config.FieldNumY]PanelInfo
	Objects map[string]*object.Object
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
