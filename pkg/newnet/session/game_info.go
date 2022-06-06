package session

import (
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
)

type GameInfo struct {
	CurrentTime time.Time
	Objects     map[string]object.Object
	Sounds      []int
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
