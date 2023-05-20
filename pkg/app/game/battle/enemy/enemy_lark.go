package enemy

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	deleteanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/delete"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayLarkMove = 8
	delayLarkAtk  = 2

	larkMoveNextStepCount = 80
)

var (
	attacker string = ""
)

type larkAtk struct {
	ownerID   string
	count     int
	attacking bool
	images    []int
}

type enemyLark struct {
	pm      EnemyParam
	imgMove []int
	count   int
	atk     larkAtk

	movePoint   [6][2]int
	movePointer int
	moveCount   int

	next common.Point
	prev common.Point
}

func (e *enemyLark) Init(objID string) error {
	e.pm.ObjectID = objID
	e.pm.DamageType = damage.TypeWater
	e.atk.ownerID = objID
	e.next = e.pm.Pos
	e.prev = e.pm.Pos
	e.count = e.pm.ActNo

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
	name, ext := GetStandImageFile(IDLark)
	e.imgMove = make([]int, 8)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 140, 130, e.imgMove); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	e.atk.images = make([]int, 3)
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
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		deleteanim.New(*img, e.pm.Pos, false)
		localanim.AnimNew(effect.Get(battlecommon.EffectTypeExplode, e.pm.Pos, 0))
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}
	e.count++

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
		if e.moveCount >= moveNum && e.pm.Pos.Y != 1 && attacker == "" {
			attacker = e.pm.ObjectID
			e.moveCount = 0
			e.atk.SetAttack()
			return false, nil
		}

		// 次の移動地点を決定
		e.moveCount++
		t := common.Point{X: e.movePoint[np][0], Y: e.movePoint[np][1]}
		if battlecommon.MoveObjectDirect(&e.pm.Pos, t, battlecommon.PanelTypeEnemy, false, field.GetPanelInfo) {
			e.next = t
		}
		return false, nil
	}
	if cnt == 0 {
		// 実際に移動
		e.prev = e.pm.Pos
		if battlecommon.MoveObjectDirect(&e.pm.Pos, e.next, battlecommon.PanelTypeEnemy, true, field.GetPanelInfo) {
			e.movePointer = np
		}
	}

	return false, nil
}

func (e *enemyLark) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	view := battlecommon.ViewPos(e.pm.Pos)
	xflip := int32(dxlib.TRUE)
	img := e.getCurrentImagePointer()

	if e.atk.attacking {
		dxlib.DrawRotaGraph(view.X+20, view.Y, 1, 0, *img, true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
		return
	}

	c := e.count % larkMoveNextStepCount
	ofsx := battlecommon.GetOffset(e.next.X, e.pm.Pos.X, e.prev.X, c, larkMoveNextStepCount, battlecommon.PanelSize.X)
	ofsy := battlecommon.GetOffset(e.next.Y, e.pm.Pos.Y, e.prev.Y, c, larkMoveNextStepCount, battlecommon.PanelSize.Y)

	dxlib.DrawRotaGraph(view.X+20+ofsx, view.Y+ofsy, 1, 0, *img, true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyLark) DamageProc(dm *damage.Damage) bool {
	return damageProc(dm, &e.pm)
}

func (e *enemyLark) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyLark) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyLark) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (a *larkAtk) SetAttack() {
	a.count = 0
	a.attacking = true
}

func (a *larkAtk) Process() {
	if a.count == 1*delayLarkAtk {
		localanim.AnimNew(skill.Get(
			skill.SkillWideShot,
			skill.Argument{
				OwnerID:    a.ownerID,
				Power:      20,
				TargetType: damage.TargetPlayer,
			},
		))
	}

	a.count++

	if a.count > len(a.images)*delayLarkAtk {
		// Reset params
		a.count = 0
		a.attacking = false
		attacker = ""
	}
}

func (e *enemyLark) getCurrentImagePointer() *int {
	if e.atk.attacking {
		n := (e.count / delayLarkAtk)
		if n >= len(e.atk.images) {
			n = len(e.atk.images) - 1
		}
		return &e.atk.images[n]
	}

	n := (e.count / delayLarkMove) % len(e.imgMove)
	return &e.imgMove[n]
}
