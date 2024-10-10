package enemy

import (
	"github.com/cockroachdb/errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	bluesActTypeStand = iota
	bluesActTypeMove
	bluesActTypeSword
	bluesActTypeShot
	bluesActTypeThrow
	bluesActTypeThrow2
	bluesActTypeDamage

	bluesActTypeMax
)

type enemyBlues struct {
	pm     EnemyParam
	state  int
	images [bluesActTypeMax][]int
}

func (e *enemyBlues) Init(objID string) error {
	e.pm.ObjectID = objID
	e.state = bluesActTypeStand

	// Load Images
	name, ext := GetStandImageFile(IDBlues)

	fname := name + "_all" + ext
	tmp := make([]int, 36)
	if res := dxlib.LoadDivGraph(fname, 36, 7, 6, 170, 156, tmp); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	e.images[bluesActTypeStand] = make([]int, 1)
	e.images[bluesActTypeStand][0] = tmp[0]

	e.images[bluesActTypeMove] = make([]int, 4)
	for i := 0; i < 4; i++ {
		e.images[bluesActTypeMove][i] = tmp[i]
	}

	e.images[bluesActTypeSword] = make([]int, 6)
	for i := 0; i < 6; i++ {
		e.images[bluesActTypeSword][i] = tmp[i+7]
	}

	e.images[bluesActTypeShot] = make([]int, 5)
	for i := 0; i < 5; i++ {
		e.images[bluesActTypeShot][i] = tmp[i+14]
	}

	e.images[bluesActTypeThrow] = make([]int, 5)
	for i := 0; i < 5; i++ {
		e.images[bluesActTypeThrow][i] = tmp[i+21]
	}

	e.images[bluesActTypeThrow2] = make([]int, 7)
	for i := 0; i < 7; i++ {
		e.images[bluesActTypeThrow2][i] = tmp[i+28]
	}

	// WIP: 使わないイメージを削除

	return nil
}

func (e *enemyBlues) End() {
	// Delete Images
	for _, imgs := range e.images {
		for _, img := range imgs {
			dxlib.DeleteGraph(img)
		}
	}
}

func (e *enemyBlues) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	return false, nil
}

func (e *enemyBlues) Draw() {
	// Show Enemy Images
	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()
	ofs := [bluesActTypeMax]point.Point{
		{X: 0, Y: 0}, // Stand
		{X: 0, Y: 0}, // Move
		{X: 0, Y: 0}, // Sword
		{X: 0, Y: 0}, // Shot
		{X: 0, Y: 0}, // Throw
		{X: 0, Y: 0}, // Throw2
		{X: 0, Y: 0}, // Damage
	}

	dxlib.DrawRotaGraph(view.X+ofs[e.state].X, view.Y+ofs[e.state].Y, 1, 0, *img, true)
}

func (e *enemyBlues) DamageProc(dm *damage.Damage) bool {
	return damageProc(dm, &e.pm)
}

func (e *enemyBlues) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyBlues) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyBlues) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (e *enemyBlues) getCurrentImagePointer() *int {
	// WIP: 実装
	return &e.images[e.state][0]
}
