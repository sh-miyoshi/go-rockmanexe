package gameinfo

import (
	"bytes"
	"encoding/gob"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
)

type PanelInfo struct {
	OwnerClientID string
	ObjectID      string
}

type Object struct {
	ID            string
	Type          int
	OwnerClientID string
	HP            int
	Pos           common.Point
	ActCount      int
	IsReverse     bool
	// TODO(他にも必要だと思うが都度追加していく)
}

type Anim struct {
	ObjectID string
	Pos      common.Point
	AnimType int
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

func (p *GameInfo) GetPanelInfo(pos common.Point) battlecommon.PanelInfo {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return battlecommon.PanelInfo{}
	}

	pn := p.Panels[pos.X][pos.Y]
	return battlecommon.PanelInfo{
		Type:     p.GetPanelType(pn.OwnerClientID),
		ObjectID: pn.ObjectID,

		// TODO: 適切な値を入れる
		Status:    battlecommon.PanelStatusNormal,
		HoleCount: 0,
	}
}

func (p *GameInfo) GetPanelType(clientID string) int {
	if p.ReverseClientID == clientID {
		return 1
	}
	return 0
}
