package enemy

import (
	"fmt"
	"math/rand"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
)

const (
	billyActMove int = iota
	billyActAttack
)

const (
	delayBillyMove = 1
	delayBillyAtk  = 5
)

type billyAct struct {
	MoveDirect int

	ownerID  string
	count    int
	typ      int
	endCount int
	pPosX    *int
	pPosY    *int
}

type enemyBilly struct {
	pm        EnemyParam
	imgMove   []int32
	imgAttack []int32
	count     int
	moveCount int
	act       billyAct
}

func (e *enemyBilly) Init(objID string) error {
	e.pm.ObjectID = objID
	e.act.pPosX = &e.pm.PosX
	e.act.pPosY = &e.pm.PosY
	e.act.typ = -1
	e.act.ownerID = objID

	// Load Images
	name, ext := GetStandImageFile(IDBilly)
	e.imgMove = make([]int32, 6)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 112, 114, e.imgMove); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	e.imgAttack = make([]int32, 8)
	fname = name + "_atk" + ext
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 124, 114, e.imgAttack); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyBilly) End() {
	// Delete Images
	for _, img := range e.imgMove {
		dxlib.DeleteGraph(img)
	}
	for _, img := range e.imgAttack {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyBilly) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		battlecommon.NewDelete(*img, e.pm.PosX, e.pm.PosY, false)
		anim.New(effect.Get(effect.TypeExplode, e.pm.PosX, e.pm.PosY, 0))
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}

	if e.act.Process() {
		return false, nil
	}

	const waitCount = 45
	const actionInterval = 75
	const moveNum = 4

	e.count++

	// Billy Actions
	if e.count < waitCount {
		return false, nil
	}

	if e.count%actionInterval == 0 {
		if e.moveCount < moveNum {
			// decide the direction to move
			// try 20 times to move
			for i := 0; i < 20; i++ {
				e.act.MoveDirect = 1 << rand.Intn(4)
				if battlecommon.MoveObject(&e.pm.PosX, &e.pm.PosY, e.act.MoveDirect, field.PanelTypeEnemy, false) {
					break
				}
			}

			e.act.SetAnim(billyActMove, len(e.imgMove)*delayBillyMove)
			e.moveCount++
		} else {
			// Attack
			e.act.SetAnim(billyActAttack, len(e.imgAttack)*delayBillyAtk)
			e.moveCount = 0
			e.count = 0
		}
	}

	return false, nil
}

func (e *enemyBilly) Draw() {
	x, y := battlecommon.ViewPos(e.pm.PosX, e.pm.PosY)
	img := e.getCurrentImagePointer()
	dxlib.DrawRotaGraph(x, y, 1, 0, *img, dxlib.TRUE)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(x, y-20, int32(e.pm.HP), draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyBilly) DamageProc(dm *damage.Damage) {
	if dm == nil {
		return
	}
	if dm.TargetType|damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY, 7))
	}
}

func (e *enemyBilly) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.TypeObject,
		ObjType:  anim.ObjTypeEnemy,
	}
}

func (e *enemyBilly) getCurrentImagePointer() *int32 {
	img := &e.imgMove[0]
	if e.act.typ != -1 {
		imgs := e.imgMove
		delay := delayBillyMove
		if e.act.typ == billyActAttack {
			imgs = e.imgAttack
			delay = delayBillyAtk
		}

		cnt := e.act.count / delay
		if cnt >= len(imgs) {
			cnt = len(imgs) - 1
		}

		img = &imgs[cnt]
	}
	return img
}

func (a *billyAct) SetAnim(animType int, endCount int) {
	a.count = 0
	a.typ = animType
	a.endCount = endCount
}

func (a *billyAct) Process() bool {
	switch a.typ {
	case -1: // No animation
		return false
	case billyActAttack:
		if a.count == 5*delayBillyAtk {
			anim.New(skill.Get(skill.SkillThunderBall, skill.Argument{
				OwnerID:    a.ownerID,
				Power:      20, // TODO: ダメージ
				TargetType: damage.TargetPlayer,
			}))
		}
	case billyActMove:
		if a.count == 4*delayBillyMove {
			battlecommon.MoveObject(a.pPosX, a.pPosY, a.MoveDirect, field.PanelTypeEnemy, true)
		}
	}

	a.count++

	if a.count > a.endCount {
		// Reset params
		a.typ = -1
		a.count = 0
		return false // finished
	}
	return true // processing now
}
