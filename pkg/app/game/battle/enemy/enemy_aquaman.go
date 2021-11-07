package enemy

import (
	"fmt"
	"math/rand"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
)

const (
	aquamanActTypeStand = iota
	aquamanActTypeMove

	aquamanActTypeMax
)

var (
	aquamanDelays = [aquamanActTypeMax]int{8, 4}
)

type enemyAquaman struct {
	pm        EnemyParam
	images    [aquamanActTypeMax][]int32
	count     int
	state     int
	nextState int

	actStand aquamanActStand
}

type aquamanActStand struct {
	count int
}

func (e *enemyAquaman) Init(objID string) error {
	e.pm.ObjectID = objID
	e.state = aquamanActTypeStand
	e.actStand.Init(10)
	e.nextState = aquamanActTypeMove

	// Load Images
	name, ext := GetStandImageFile(IDAquaman)

	fname := name + "_stand" + ext
	e.images[aquamanActTypeStand] = make([]int32, 9)
	if res := dxlib.LoadDivGraph(fname, 9, 9, 1, 62, 112, e.images[aquamanActTypeStand]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = name + "_move" + ext
	e.images[aquamanActTypeMove] = make([]int32, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 64, 126, e.images[aquamanActTypeMove]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyAquaman) End() {
	// Delete Images
	for i := 0; i < aquamanActTypeMax; i++ {
		for j := 0; j < len(e.images[i]); j++ {
			dxlib.DeleteGraph(e.images[i][j])
		}
		e.images[i] = []int32{}
	}
}

func (e *enemyAquaman) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	switch e.state {
	case aquamanActTypeStand:
		if e.actStand.Process() {
			e.state = e.nextState
			e.count = 0
			return false, nil
		}
	case aquamanActTypeMove:
		if e.count == 5*aquamanDelays[aquamanActTypeMove] {
			// move
			for i := 0; i < 10; i++ {
				if battlecommon.MoveObjectDirect(
					&e.pm.PosX,
					&e.pm.PosY,
					rand.Intn(field.FieldNumX/2)+field.FieldNumX/2,
					rand.Intn(field.FieldNumY),
					field.PanelTypeEnemy,
					true,
					field.GetPanelInfo,
				) {
					break
				}
			}
			// standへ
			e.actStand.Init(60)
			e.state = aquamanActTypeStand
		}
	}
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
	n := (e.count / aquamanDelays[e.state]) % len(e.images[e.state])
	return &e.images[e.state][n]
}

func (a *aquamanActStand) Init(waitCount int) {
	a.count = waitCount
}

func (a *aquamanActStand) Process() bool {
	if a.count <= 0 {
		return true
	}
	a.count--
	return false
}
