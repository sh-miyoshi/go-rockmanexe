package enemy

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	deleteanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/delete"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayCirkillMove   = 8
	delayCirkillAttack = 2

	cirkillMoveNextStepCount = 80
)

type cirKillAttack struct {
	count     int
	attacking bool
}

type enemyCirKill struct {
	pm        EnemyParam
	atk       cirKillAttack
	imgMove   []int
	imgAttack []int
	count     int

	movePoint   [6][2]int
	movePointer int
	moveCount   int

	next point.Point
	prev point.Point
}

func (e *enemyCirKill) Init(objID string) error {
	e.pm.ObjectID = objID
	e.next = e.pm.Pos
	e.prev = e.pm.Pos
	e.count = e.pm.ActNo

	// TODO: check
	for i := 0; i < 6; i++ {
		// x座標
		if e.pm.Pos.X == 3 {
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

		if e.movePoint[i][0] == e.pm.Pos.X && e.movePoint[i][1] == e.pm.Pos.Y {
			e.movePointer = i
		}
	}

	// Load Images
	name, ext := GetStandImageFile(IDCirKill)
	e.imgMove = make([]int, 4)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 96, e.imgMove); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	e.imgAttack = make([]int, 2)
	fname = name + "_atk" + ext
	if res := dxlib.LoadDivGraph(fname, 2, 2, 1, 140, 104, e.imgAttack); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyCirKill) End() {
	// Delete Images
	for _, img := range e.imgMove {
		dxlib.DeleteGraph(img)
	}
	for _, img := range e.imgAttack {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyCirKill) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		deleteanim.New(*img, e.pm.Pos, false)
		localanim.AnimNew(effect.Get(resources.EffectTypeExplode, e.pm.Pos, 0))
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}

	if e.pm.ParalyzedCount > 0 {
		e.pm.ParalyzedCount--
		return false, nil
	}

	e.count++

	if e.atk.attacking {
		e.atk.Process()
		return false, nil
	}

	const waitCount = 20
	if e.count < waitCount {
		return false, nil
	}

	// TODO: move

	return false, nil
}

func (e *enemyCirKill) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	defaultOfsX := -10
	defaultOfsY := 15

	view := battlecommon.ViewPos(e.pm.Pos)
	xflip := int32(dxlib.TRUE)
	img := e.getCurrentImagePointer()

	if e.atk.attacking {
		dxlib.DrawRotaGraph(view.X+defaultOfsX, view.Y+defaultOfsY, 1, 0, *img, true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
		return
	}

	c := e.count % cirkillMoveNextStepCount
	ofsx := battlecommon.GetOffset(e.next.X, e.pm.Pos.X, e.prev.X, c, cirkillMoveNextStepCount, battlecommon.PanelSize.X)
	ofsy := battlecommon.GetOffset(e.next.Y, e.pm.Pos.Y, e.prev.Y, c, cirkillMoveNextStepCount, battlecommon.PanelSize.Y)

	dxlib.DrawRotaGraph(view.X+ofsx+defaultOfsX, view.Y+ofsy+defaultOfsY, 1, 0, *img, true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
	drawParalysis(view.X+ofsx+defaultOfsX, view.Y+ofsy+defaultOfsY, *img, e.pm.ParalyzedCount)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}

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

func (e *enemyCirKill) getCurrentImagePointer() *int {
	if e.atk.attacking {
		n := (e.count / delayLarkAtk)
		if n >= len(e.imgAttack) {
			n = len(e.imgAttack) - 1
		}
		return &e.imgAttack[n]
	}

	n := (e.count / delayLarkMove) % len(e.imgMove)
	return &e.imgMove[n]
}

func (a *cirKillAttack) Process() {
	// TODO: attack

	a.count++
}
