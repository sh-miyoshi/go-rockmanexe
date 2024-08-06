package enemy

import (
	"fmt"
	"math/rand"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayShrimpyMove        = 16
	delayShrimpyAttack      = 4
	shrimpyActNextStepCount = 90
)

const (
	shrimpyStateWait = iota
	shrimpyStateMove
	shrimpyStateAtk
)

type shrimpyAttack struct {
	ownerID string
	count   int
	images  []int
}

type enemyShrimpy struct {
	pm        EnemyParam
	imgMove   []int
	count     int
	atk       shrimpyAttack
	waitCount int
	state     int
	nextState int
	nextY     int
	prevY     int
	direct    int
	moveCount int
	prevOfsY  int
}

func (e *enemyShrimpy) Init(objID string) error {
	e.pm.ObjectID = objID
	e.waitCount = 20
	e.state = shrimpyStateWait
	e.nextState = shrimpyStateMove
	e.nextY = e.pm.Pos.Y
	e.prevY = e.pm.Pos.Y
	e.direct = config.DirectUp
	e.atk.ownerID = objID
	e.setMoveCount()

	// Load Images
	name, ext := GetStandImageFile(IDShrimpy)
	e.imgMove = make([]int, 4)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 110, 112, e.imgMove); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}
	fname = name + "_atk" + ext
	e.atk.images = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 110, 128, e.atk.images); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}
	return nil
}

func (e *enemyShrimpy) End() {
	// Delete Images
	for _, img := range e.imgMove {
		dxlib.DeleteGraph(img)
	}
	for _, img := range e.atk.images {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyShrimpy) Process() (bool, error) {
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

	// Enemy Logic
	switch e.state {
	case shrimpyStateWait:
		e.waitCount--
		if e.waitCount <= 0 {
			e.setState(e.nextState)
			return false, nil
		}
	case shrimpyStateMove:
		if e.count == 0 {
			e.count = shrimpyActNextStepCount/2 + 1
		}

		cnt := e.count % shrimpyActNextStepCount
		if cnt == 0 {
			// Update current pos
			e.prevY = e.pm.Pos.Y
			if battlecommon.MoveObjectDirect(&e.pm.Pos, point.Point{X: e.pm.Pos.X, Y: e.nextY}, battlecommon.PanelTypeEnemy, true, field.GetPanelInfo) {
				e.moveCount--
			}
		} else if cnt == shrimpyActNextStepCount/2 {
			// 次の行動を決定
			if e.moveCount <= 0 {
				e.setMoveCount()
				e.waitCount = 10
				e.nextState = shrimpyStateAtk
				e.setState(shrimpyStateWait)
				return false, nil
			}

			if e.direct == config.DirectUp {
				if e.nextY > 0 {
					e.nextY--
				}

				if e.nextY == 0 {
					e.direct = config.DirectDown
				}
			} else { // Down
				if e.nextY < battlecommon.FieldNum.Y-1 {
					e.nextY++
				}

				if e.nextY == battlecommon.FieldNum.Y-1 {
					e.direct = config.DirectUp
				}
			}
		}
	case shrimpyStateAtk:
		if e.count == 0 {
			e.atk.Set()
		}

		if e.atk.Process() {
			e.waitCount = 40
			e.nextState = shrimpyStateMove
			e.setState(shrimpyStateWait)
			return false, nil
		}
	}

	e.count++
	return false, nil
}

func (e *enemyShrimpy) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()
	var ofsy int
	if e.state == shrimpyStateMove {
		c := e.count % shrimpyActNextStepCount
		if c == 0 || c == shrimpyActNextStepCount/2 {
			ofsy = e.prevOfsY
		} else {
			ofsy = battlecommon.GetOffset(e.nextY, e.pm.Pos.Y, e.prevY, c, shrimpyActNextStepCount, battlecommon.PanelSize.Y)
			e.prevOfsY = ofsy
		}
	}
	dxlib.DrawRotaGraph(view.X, view.Y+ofsy, 1, 0, *img, true)

	drawParalysis(view.X, view.Y, *img, e.pm.ParalyzedCount)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40+ofsy, e.pm.HP, draw.NumberOption{
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
	if e.state == shrimpyStateAtk {
		n := (e.atk.count / delayShrimpyAttack)
		if n >= len(e.atk.images) {
			n = len(e.atk.images) - 1
		}
		return &e.atk.images[n]
	}

	n := (e.count / delayShrimpyMove) % len(e.imgMove)
	return &e.imgMove[n]
}

func (e *enemyShrimpy) setMoveCount() {
	e.moveCount = 3 + rand.Intn(2)
}

func (e *enemyShrimpy) setState(state int) {
	e.state = state
	e.count = 0
}

func (a *shrimpyAttack) Set() {
	a.count = 0
	localanim.AnimNew(skill.Get(
		resources.SkillShrimpyAttack,
		skillcore.Argument{
			OwnerID:    a.ownerID,
			Power:      20,
			TargetType: damage.TargetPlayer,
		},
	))
}

func (a *shrimpyAttack) Process() bool {
	a.count++
	return a.count >= len(a.images)*delayShrimpyAttack
}
