package enemy

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	cirkillActWait int = iota
	cirkillActMove
	cirkillActAttack
)

type cirKillAct struct {
	current   int
	waitCount int
	next      int
}

type enemyCirKill struct {
	pm  EnemyParam
	act cirKillAct
}

func (e *enemyCirKill) Init(objID string) error {
	e.pm.ObjectID = objID
	e.act.current = cirkillActWait
	e.act.next = cirkillActMove
	e.act.waitCount = 60

	// Load Images
	// TODO

	return nil
}

func (e *enemyCirKill) End() {
	// Delete Images
	// TODO
}

func (e *enemyCirKill) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic

	// TODO
	/*
		Initial Wait
		外周を回る
		if y line == player
			攻撃モーション
	*/

	return false, nil
}

func (e *enemyCirKill) Draw() {
	// Show Enemy Images
	// TODO
}

func (e *enemyCirKill) DamageProc(dm *damage.Damage) bool {
	return damageProc(dm, &e.pm)
}

func (e *enemyCirKill) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyCirKill) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyCirKill) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (a *cirKillAct) Process() {
	switch a.current {
	case cirkillActWait:
		a.waitCount--
		if a.waitCount <= 0 {
			a.set(a.next)
			return
		}
	}
}

func (a *cirKillAct) set(next int) {
	logger.Info("CirKill act change from %d to %d", a.current, next)
	a.current = next
}
