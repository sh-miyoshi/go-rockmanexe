package enemy

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
)

type enemyForte struct {
	pm     EnemyParam
	images []int
	count  int
}

func (e *enemyForte) Init(objID string) error {
	e.pm.ObjectID = objID

	// Load Images
	return nil
}

func (e *enemyForte) End() {
	// Delete Images
}

func (e *enemyForte) Process() (bool, error) {
	// WIP

	if e.pm.ParalyzedCount > 0 {
		e.pm.ParalyzedCount--
		return false, nil
	}

	// Enemy Logic
	if e.pm.InvincibleCount > 0 {
		e.pm.InvincibleCount--
	}

	// WIP

	e.count++
	return false, nil
}

func (e *enemyForte) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	// Show Enemy Images
	// WIP
}

func (e *enemyForte) DamageProc(dm *damage.Damage) bool {
	// WIP
	return false
}

func (e *enemyForte) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyForte) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyForte) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}
