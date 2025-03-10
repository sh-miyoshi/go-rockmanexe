package enemy

import (
	"math/rand"

	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common/deleteanim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
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
	pPos     *point.Point
}

type enemyBilly struct {
	pm        EnemyParam
	imgMove   []int
	imgAttack []int
	count     int
	moveCount int
	act       billyAct
}

func (e *enemyBilly) Init(objID string) error {
	e.pm.ObjectID = objID
	e.pm.DamageElement = damage.ElementElec
	e.act.pPos = &e.pm.Pos
	e.act.typ = -1
	e.act.ownerID = objID

	// Load Images
	name, ext := GetStandImageFile(IDBilly)
	e.imgMove = make([]int, 6)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 112, 114, e.imgMove); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	e.imgAttack = make([]int, 8)
	fname = name + "_atk" + ext
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 124, 114, e.imgAttack); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
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

func (e *enemyBilly) Update() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		deleteanim.New(*img, e.pm.Pos, false)
		localanim.EffectAnimNew(effect.Get(resources.EffectTypeExplode, e.pm.Pos, 0))
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}

	if e.pm.ParalyzedCount > 0 {
		e.pm.ParalyzedCount--
		return false, nil
	}

	if e.act.Update() {
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
				if battlecommon.MoveObject(&e.pm.Pos, e.act.MoveDirect, battlecommon.PanelTypeEnemy, false, field.GetPanelInfo) {
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

func (e *enemyBilly) DamageProc(dm *damage.Damage) bool {
	return damageProc(dm, &e.pm)
}

func (e *enemyBilly) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyBilly) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyBilly) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (e *enemyBilly) AddBarrier(hp int) {}

func (e *enemyBilly) getCurrentImagePointer() *int {
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

func (a *billyAct) Update() bool {
	switch a.typ {
	case -1: // No animation
		return false
	case billyActAttack:
		if a.count == 5*delayBillyAtk {
			localanim.SkillAnimNew(skill.Get(resources.SkillThunderBall, skillcore.Argument{
				OwnerID:    a.ownerID,
				Power:      20,
				TargetType: damage.TargetPlayer,
			}))
		}
	case billyActMove:
		if a.count == 4*delayBillyMove {
			battlecommon.MoveObject(a.pPos, a.MoveDirect, battlecommon.PanelTypeEnemy, true, field.GetPanelInfo)
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
