package gameinfo

import (
	"bytes"
	"encoding/gob"

	"github.com/google/uuid"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
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

func (p *GameInfo) Init(myClientID, opponentClientID string) {
	p.ClientID = myClientID
	for y := 0; y < config.FieldNumY; y++ {
		hx := config.FieldNumX / 2
		for x := 0; x < hx; x++ {
			p.Panels[x][y] = PanelInfo{
				OwnerClientID: myClientID,
				Status:        battlecommon.PanelStatusNormal,
				HoleCount:     0,
				ObjExists:     false,
			}
			p.Panels[x+hx][y] = PanelInfo{
				OwnerClientID: opponentClientID,
				Status:        battlecommon.PanelStatusNormal,
				HoleCount:     0,
				ObjExists:     false,
			}
		}
	}
	p.Objects = []Object{}
	p.Anims = []Anim{}
	p.Effects = []Effect{}
	p.Sounds = []Sound{}
}

func (p *GameInfo) GetObject(id string) *Object {
	for i, obj := range p.Objects {
		if obj.ID == id {
			return &p.Objects[i]
		}
	}
	return nil
}

func (p *GameInfo) Update(objects []Object, anims []Anim, effects []Effect, sounds []Sound) {
	p.Objects = append([]Object{}, objects...)
	p.Anims = append([]Anim{}, anims...)
	p.Sounds = append([]Sound{}, sounds...)

	p.Effects = []Effect{}
	for _, e := range effects {
		if p.ClientID != e.OwnerClientID {
			e.Pos.X = battlecommon.FieldNum.X - e.Pos.X - 1
		}
		p.Effects = append(p.Effects, e)
	}

	// Panel updates
	// Cleanup object id at first
	for y := 0; y < battlecommon.FieldNum.Y; y++ {
		for x := 0; x < battlecommon.FieldNum.X; x++ {
			p.Panels[x][y].ObjectID = ""
		}
	}
	for _, obj := range objects {
		p.Panels[obj.Pos.X][obj.Pos.Y].ObjectID = obj.ID
		if p.Panels[obj.Pos.X][obj.Pos.Y].Status == battlecommon.PanelStatusCrack {
			p.Panels[obj.Pos.X][obj.Pos.Y].ObjExists = true
		}
	}

	// Panel status update
	for y := 0; y < battlecommon.FieldNum.Y; y++ {
		for x := 0; x < battlecommon.FieldNum.X; x++ {
			if p.Panels[x][y].HoleCount > 0 {
				p.Panels[x][y].HoleCount--
			}

			switch p.Panels[x][y].Status {
			case battlecommon.PanelStatusHole:
				if p.Panels[x][y].HoleCount <= battlecommon.PanelReturnAnimCount {
					p.Panels[x][y].Status = battlecommon.PanelStatusNormal
				}
			case battlecommon.PanelStatusCrack:
				// Objectが乗って離れたらHole状態へ
				if p.Panels[x][y].ObjExists && p.Panels[x][y].ObjectID == "" {
					p.Sounds = append(p.Sounds, Sound{ID: uuid.New().String(), Type: int(resources.SEPanelBreak)})
					p.Panels[x][y].ObjExists = false
					p.Panels[x][y].Status = battlecommon.PanelStatusHole
					p.Panels[x][y].HoleCount = battlecommon.DefaultPanelHoleEndCount
				}
			}
		}
	}
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
