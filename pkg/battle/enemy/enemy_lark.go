package enemy

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
)

const (
	delayLarkMove = 8
	delayLarkAtk  = 2

	larkMoveNextStepCount = 80
)

type larkAtk struct {
	ownerID   string
	count     int
	attacking bool
	images    []int32
}

type enemyLark struct {
	pm      EnemyParam
	imgMove []int32
	count   int
	atk     larkAtk

	movePoint   [6][2]int
	movePointer int
	moveCount   int

	nextX int
	nextY int
	prevX int
	prevY int
}

func (e *enemyLark) Init(objID string) error {
	e.pm.ObjectID = objID
	e.atk.ownerID = objID
	e.nextX = e.pm.PosX
	e.nextY = e.pm.PosY
	e.prevX = e.pm.PosX
	e.prevY = e.pm.PosY

	for i := 0; i < 6; i++ {
		// x座標
		if e.pm.PosX == 3 {
			// Pattern 1. 前2列を周回
			e.movePoint[i][0] = (i / 3) + 3
		} else {
			// Pattern 2. 後2列を周回
			e.movePoint[i][0] = (i / 3) + 4
		}

		// y座標
		if i < 3 {
			e.movePoint[i][1] = i
		} else {
			e.movePoint[i][1] = 5 - i
		}

		if e.movePoint[i][0] == e.pm.PosX && e.movePoint[i][1] == e.pm.PosY {
			e.movePointer = i
		}
	}

	// Load Images
	name, ext := GetStandImageFile(IDLark)
	e.imgMove = make([]int32, 8)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 140, 130, e.imgMove); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	e.atk.images = make([]int32, 3)
	fname = name + "_atk" + ext
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 160, 124, e.atk.images); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyLark) End() {
	// Delete Images
	for _, img := range e.imgMove {
		dxlib.DeleteGraph(img)
	}
	for _, img := range e.atk.images {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyLark) Process() (bool, error) {
	e.count++
	// TODO Return true if finished(e.g. hp=0)

	if e.atk.attacking {
		e.atk.Process()
		return false, nil
	}

	const waitCount = 20
	const moveNum = 4

	if e.count < waitCount {
		return false, nil
	}

	cnt := e.count % larkMoveNextStepCount
	np := (e.movePointer + 1) % 6

	if cnt == larkMoveNextStepCount/2 {
		if e.moveCount >= moveNum {
			e.moveCount = 0
			// TODO
			// e.atk.SetAtttack()
			// return false, nil
		}

		// 次の移動地点を決定
		e.moveCount++
		tx := e.movePoint[np][0]
		ty := e.movePoint[np][1]
		if battlecommon.MoveObjectDirect(&e.pm.PosX, &e.pm.PosY, tx, ty, field.PanelTypeEnemy, false) {
			e.nextX = tx
			e.nextY = ty
		}
		return false, nil
	}
	if cnt == 0 {
		// 実際に移動
		e.prevX = e.pm.PosX
		e.prevY = e.pm.PosY
		if battlecommon.MoveObjectDirect(&e.pm.PosX, &e.pm.PosY, e.nextX, e.nextY, field.PanelTypeEnemy, true) {
			e.movePointer = np
		}
	}

	return false, nil
}

func (e *enemyLark) Draw() {
	x, y := battlecommon.ViewPos(e.pm.PosX, e.pm.PosY)
	xflip := int32(dxlib.TRUE)

	if e.atk.attacking {
		n := (e.count / delayLarkAtk)
		if n >= len(e.atk.images) {
			n = len(e.atk.images) - 1
		}
		dxlib.DrawRotaGraph(x+20, y, 1, 0, e.atk.images[n], dxlib.TRUE, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
		return
	}

	c := e.count % larkMoveNextStepCount
	ofsx := battlecommon.GetOffset(e.nextX, e.pm.PosX, e.prevX, c, larkMoveNextStepCount, field.PanelSizeX)
	ofsy := battlecommon.GetOffset(e.nextY, e.pm.PosY, e.prevY, c, larkMoveNextStepCount, field.PanelSizeY)

	n := (e.count / delayLarkMove) % len(e.imgMove)
	dxlib.DrawRotaGraph(x+20+int32(ofsx), y+int32(ofsy), 1, 0, e.imgMove[n], dxlib.TRUE, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(x, y+40, int32(e.pm.HP), draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyLark) DamageProc(dm *damage.Damage) {
	if dm == nil {
		return
	}
	if dm.TargetType|damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY, 5))
	}
}

func (e *enemyLark) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.TypeObject,
		ObjType:  anim.ObjTypeEnemy,
	}
}

func (a *larkAtk) SetAtttack() {
	a.count = 0
	a.attacking = true
}

func (a *larkAtk) Process() {
	// TODO damage登録

	a.count++

	if a.count > len(a.images)*delayBillyAtk {
		// Reset params
		a.count = 0
		a.attacking = false
	}
}
