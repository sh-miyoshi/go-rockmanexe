package enemy

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayShrimpyMove = 16
)

type enemyShrimpy struct {
	pm      EnemyParam
	imgMove []int
	count   int
}

func (e *enemyShrimpy) Init(objID string) error {
	e.pm.ObjectID = objID

	// Load Images
	name, ext := GetStandImageFile(IDShrimpy)
	e.imgMove = make([]int, 4)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 110, 112, e.imgMove); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}
	return nil
}

func (e *enemyShrimpy) End() {
	// Delete Images
	for _, img := range e.imgMove {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyShrimpy) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	e.count++
	return false, nil
}

func (e *enemyShrimpy) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, *img, true)

	drawParalysis(view.X, view.Y, *img, e.pm.ParalyzedCount)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
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

func (e *enemyShrimpy) getCurrentImagePointer() *int {
	n := (e.count / delayShrimpyMove) % len(e.imgMove)
	return &e.imgMove[n]
}
