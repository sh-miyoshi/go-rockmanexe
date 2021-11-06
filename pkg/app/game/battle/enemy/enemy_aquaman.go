package enemy

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
)

const (
	aquamanActStand = iota

	aquamanActMax
)

var (
	aquamanDelays = [aquamanActMax]int{8}
)

type enemyAquaman struct {
	pm     EnemyParam
	images [aquamanActMax][]int32
	count  int
}

func (e *enemyAquaman) Init(objID string) error {
	e.pm.ObjectID = objID

	// Load Images
	name, ext := GetStandImageFile(IDAquaman)

	fname := name + "_stand" + ext
	e.images[aquamanActStand] = make([]int32, 9)
	if res := dxlib.LoadDivGraph(fname, 9, 9, 1, 62, 112, e.images[aquamanActStand]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyAquaman) End() {
	// Delete Images
	for i := 0; i < aquamanActMax; i++ {
		for j := 0; j < len(e.images[i]); j++ {
			dxlib.DeleteGraph(e.images[i][j])
		}
		e.images[i] = []int32{}
	}
}

func (e *enemyAquaman) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	// TODO
	e.count++
	return false, nil
}

func (e *enemyAquaman) Draw() {
	// Show Enemy Images
	x, y := battlecommon.ViewPos(e.pm.PosX, e.pm.PosY)
	img := e.getCurrentImagePointer()

	dxlib.DrawRotaGraph(x, y, 1, 0, *img, dxlib.TRUE)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(x, y+40, int32(e.pm.HP), draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyAquaman) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}
	if dm.TargetType&damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY, 5))
		return true
	}
	return false
}

func (e *enemyAquaman) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.AnimTypeObject,
	}
}

func (e *enemyAquaman) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyAquaman) getCurrentImagePointer() *int32 {
	// TODO
	n := (e.count / aquamanDelays[0]) % len(e.images[0])
	return &e.images[0][n]
}
