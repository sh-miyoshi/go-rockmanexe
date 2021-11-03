package enemy

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
)

const (
	delayBoomerMove = 32

	boomerMoveNextStepCount = 120
)

type enemyBoomer struct {
	pm      EnemyParam
	imgMove []int32
	count   int
	direct  int
	nextY   int
	prevY   int
}

func (e *enemyBoomer) Init(objID string) error {
	e.pm.ObjectID = objID
	e.nextY = e.pm.PosY
	e.prevY = e.pm.PosY

	e.direct = common.DirectUp

	// Load Images
	name, ext := GetStandImageFile(IDBoomer)
	e.imgMove = make([]int32, 4)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 114, 102, e.imgMove); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyBoomer) End() {
	// Delete Images
	for _, img := range e.imgMove {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyBoomer) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)

	// Enemy Logic
	e.count++
	if e.count%boomerMoveNextStepCount == 0 {
		// Update current pos
		e.prevY = e.pm.PosY
		e.pm.PosY = e.nextY

		if e.direct == common.DirectUp {
			if e.nextY > 0 {
				e.nextY--
			}

			if e.nextY == 0 {
				e.direct = common.DirectDown
			}
		} else { // Down
			if e.nextY < field.FieldNumY-1 {
				e.nextY++
			}

			if e.nextY == field.FieldNumY-1 {
				e.direct = common.DirectUp
			}
		}
	}

	return false, nil
}

func (e *enemyBoomer) Draw() {
	// Show Enemy Images
	x, y := battlecommon.ViewPos(e.pm.PosX, e.pm.PosY)
	img := e.getCurrentImagePointer()

	c := e.count % boomerMoveNextStepCount
	ofsy := battlecommon.GetOffset(e.nextY, e.pm.PosY, e.prevY, c, boomerMoveNextStepCount, field.PanelSizeY)

	dxlib.DrawRotaGraph(x, y+int32(ofsy), 1, 0, *img, dxlib.TRUE)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(x, y+40, int32(e.pm.HP), draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyBoomer) DamageProc(dm *damage.Damage) bool {
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

func (e *enemyBoomer) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.AnimTypeObject,
	}
}

func (e *enemyBoomer) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyBoomer) getCurrentImagePointer() *int32 {
	n := (e.count / delayBoomerMove) % len(e.imgMove)
	return &e.imgMove[n]
}
