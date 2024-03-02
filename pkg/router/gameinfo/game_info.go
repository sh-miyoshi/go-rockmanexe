package gameinfo

import (
	"bytes"
	"encoding/gob"

	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	QueueTypeAction int = iota
	QueueTypeEffect
	QueueTypeSound

	QueueTypeMax
)

type PanelInfo struct {
	OwnerClientID string
	ObjectID      string
	Status        int
	HoleCount     int
	ObjExists     bool
}

type Object struct {
	ID            string
	Type          int
	OwnerClientID string
	HP            int
	Pos           point.Point
	ActCount      int
	IsReverse     bool
	IsInvincible  bool
}

type Anim struct {
	ObjectID  string
	Pos       point.Point
	DrawType  int
	AnimType  int
	ActCount  int
	DrawParam []byte
}

type Effect struct {
	ID            string
	Pos           point.Point
	Type          int
	RandRange     int
	OwnerClientID string
}

type Sound struct {
	ID   string
	Type int
}

// Client側に送られるデータ
type GameInfo struct {
	Panels   [config.FieldNumX][config.FieldNumY]PanelInfo
	Objects  []Object
	Anims    []Anim
	Effects  []Effect
	Sounds   []Sound
	ClientID string
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

func (p *GameInfo) GetPanelInfo(pos point.Point) battlecommon.PanelInfo {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return battlecommon.PanelInfo{}
	}

	pn := p.Panels[pos.X][pos.Y]
	return battlecommon.PanelInfo{
		Type:      p.getPanelType(pn.OwnerClientID),
		ObjectID:  pn.ObjectID,
		Status:    pn.Status,
		HoleCount: pn.HoleCount,
	}
}

func (p *GameInfo) PanelBreak(pos point.Point) {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return
	}

	if p.Panels[pos.X][pos.Y].Status == battlecommon.PanelStatusHole {
		return
	}

	if p.Panels[pos.X][pos.Y].ObjectID != "" {
		p.Panels[pos.X][pos.Y].Status = battlecommon.PanelStatusCrack
	} else {
		p.Panels[pos.X][pos.Y].Status = battlecommon.PanelStatusHole
		p.Panels[pos.X][pos.Y].HoleCount = battlecommon.DefaultPanelHoleEndCount
	}
}

func (p *GameInfo) getPanelType(clientID string) int {
	if p.ClientID == clientID {
		return battlecommon.PanelTypePlayer
	}
	return battlecommon.PanelTypeEnemy
}
