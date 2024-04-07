package enemy

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
)

type enemyShrimpy struct {
	pm EnemyParam
}

func (e *enemyShrimpy) Init(objID string) error {
	e.pm.ObjectID = objID

	// Load Images
	return nil
}

func (e *enemyShrimpy) End() {
	// Delete Images
}

func (e *enemyShrimpy) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	return false, nil
}

func (e *enemyShrimpy) Draw() {
	// Show Enemy Images
}

func (e *enemyShrimpy) DamageProc(dm *damage.Damage) bool {
	return damageProc(dm, &e.pm)
}

func (e *enemyShrimpy) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyShrimpy) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyShrimpy) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}
