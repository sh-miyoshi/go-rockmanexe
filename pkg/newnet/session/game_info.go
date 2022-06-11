package session

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
)

type PanelInfo struct {
	OwnerClientID string
	ShowHitArea   bool
	Status        int
}

type GameInfo struct {
	CurrentTime time.Time
	Objects     map[string]object.Object
	Panels      [config.FieldNumX][config.FieldNumY]PanelInfo
	Sounds      []int
}

func NewGameInfo() *GameInfo {
	res := &GameInfo{
		CurrentTime: time.Now(),
		Objects:     make(map[string]object.Object),
	}

	for x := 0; x < config.FieldNumX; x++ {
		for y := 0; y < config.FieldNumY; y++ {
			res.Panels[x][y] = PanelInfo{}
		}
	}

	return res
}

func (g *GameInfo) UpdateObject(obj object.Object) {
	g.Objects[obj.ID] = obj
}

func (g *GameInfo) RemoveObject(id string) {
	delete(g.Objects, id)
}

func (g *GameInfo) AddSkill() {
	// TODO
}

func (g *GameInfo) AddDamage() {
	// TODO
}

func (g *GameInfo) AddEffect() {
	// TODO
}

func (g *GameInfo) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(g)
	return buf.Bytes()
}

func (g *GameInfo) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(g)
}
