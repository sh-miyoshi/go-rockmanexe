package enemy

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
)

type enemyColdman struct {
	pm EnemyParam
}

func (e *enemyColdman) Init(objID string) error {
	e.pm.ObjectID = objID

	// Load Images
	return nil
}

func (e *enemyColdman) End() {
	// Delete Images
}

func (e *enemyColdman) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	return false, nil
}

func (e *enemyColdman) Draw() {
	// Show Enemy Images
}

func (e *enemyColdman) DamageProc(dm *damage.Damage) bool {
	return damageProc(dm, &e.pm)
}

func (e *enemyColdman) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyColdman) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyColdman) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}
