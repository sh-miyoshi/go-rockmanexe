package session

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
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
	HitDamages  []damage.Damage
	Effects     []effect.Effect
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

func (g *GameInfo) InitPanel(myClientID, enemyClientID string) {
	for x := 0; x < config.FieldNumX; x++ {
		id := myClientID
		if x > 2 {
			id = enemyClientID
		}
		for y := 0; y < config.FieldNumY; y++ {
			g.Panels[x][y].OwnerClientID = id
		}
	}
}

func (g *GameInfo) Cleanup() {
	g.Sounds = []int{}
	g.HitDamages = []damage.Damage{}
	g.Effects = []effect.Effect{}
}

func (g *GameInfo) UpdateObject(obj object.Object, isMyObj bool) {
	obj.Count = 0
	if obj.UpdateBaseTime {
		obj.BaseTime = time.Now()
	}

	if !isMyObj {
		obj.X = config.FieldNumX - obj.X - 1
		obj.PrevX = config.FieldNumX - obj.PrevX - 1
		obj.TargetX = config.FieldNumX - obj.TargetX - 1
	}

	if o, ok := g.Objects[obj.ID]; ok {
		if !obj.UpdateBaseTime {
			obj.BaseTime = o.BaseTime
		}
		logger.Debug("Update Object: %+v", obj)
	} else {
		logger.Debug("New Object: %+v", obj)
	}

	g.Objects[obj.ID] = obj
}

func (g *GameInfo) RemoveObject(id string) {
	delete(g.Objects, id)
}

func (g *GameInfo) AddSkill() {
	panic("TODO")
}

func (g *GameInfo) AddDamages(dm []damage.Damage) {
	if len(dm) == 0 {
		return
	}

	g.HitDamages = append(g.HitDamages, dm...)
}

func (g *GameInfo) AddEffect(eff effect.Effect, isMyEff bool) {
	if !isMyEff {
		eff.X = config.FieldNumX - eff.X - 1
		eff.ViewOfsX = config.FieldNumX - eff.ViewOfsX - 1
	}

	g.Effects = append(g.Effects, eff)
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
