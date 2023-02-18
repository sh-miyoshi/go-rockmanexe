package enemy

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	deleteanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/delete"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayBoomerMove = 32
	delayBoomerAtk  = 2

	boomerActNextStepCount = 120
)

const (
	boomerStateWait = iota
	boomerStateMove
	boomerStateAtk
)

type boomerAtk struct {
	ownerID string
	count   int
	images  []int
	atkID   string
}

type enemyBoomer struct {
	pm        EnemyParam
	imgMove   []int
	count     int
	atk       boomerAtk
	direct    int
	nextY     int
	prevY     int
	state     int
	nextState int
	waitCount int
	prevOfsY  int
}

func (e *enemyBoomer) Init(objID string) error {
	e.pm.ObjectID = objID
	e.pm.DamageType = damage.TypeWood
	e.atk.ownerID = objID
	e.nextY = e.pm.Pos.Y
	e.prevY = e.pm.Pos.Y
	e.direct = common.DirectUp
	e.waitCount = 20
	e.state = boomerStateWait
	e.nextState = boomerStateMove

	// Load Images
	name, ext := GetStandImageFile(IDBoomer)
	e.imgMove = make([]int, 4)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 114, 102, e.imgMove); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	tmp := make([]int, 5)
	fname = name + "_atk" + ext
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 136, 104, tmp); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}
	for i := len(tmp) - 1; i >= 0; i-- {
		e.atk.images = append(e.atk.images, tmp[i])
	}

	return nil
}

func (e *enemyBoomer) End() {
	// Delete Images
	for _, img := range e.imgMove {
		dxlib.DeleteGraph(img)
	}
	for _, img := range e.atk.images {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyBoomer) Process() (bool, error) {
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		deleteanim.New(*img, e.pm.Pos, false)
		anim.New(effect.Get(effect.TypeExplode, e.pm.Pos, 0))
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}

	// Enemy Logic
	switch e.state {
	case boomerStateWait:
		e.waitCount--
		if e.waitCount <= 0 {
			e.setState(e.nextState)
			return false, nil
		}
	case boomerStateMove:
		if e.count == 0 {
			e.count = boomerActNextStepCount/2 + 1
		}

		cnt := e.count % boomerActNextStepCount
		if cnt == 0 {
			// Update current pos
			e.prevY = e.pm.Pos.Y
			e.pm.Pos.Y = e.nextY
		}

		if cnt == boomerActNextStepCount/2 {
			// 次の行動を決定
			if e.pm.Pos.Y == 0 || e.pm.Pos.Y == battlecommon.FieldNum.Y-1 {
				e.state = boomerStateWait
				e.nextState = boomerStateAtk
				e.waitCount = 60
				e.atk.Init()
			}

			if e.direct == common.DirectUp {
				if e.nextY > 0 {
					e.nextY--
				}

				if e.nextY == 0 {
					e.direct = common.DirectDown
				}
			} else { // Down
				if e.nextY < battlecommon.FieldNum.Y-1 {
					e.nextY++
				}

				if e.nextY == battlecommon.FieldNum.Y-1 {
					e.direct = common.DirectUp
				}
			}
		}
	case boomerStateAtk:
		if e.atk.Process() {
			e.waitCount = 60
			e.nextState = boomerStateMove
			e.setState(boomerStateWait)
			return false, nil
		}
	}

	e.count++
	return false, nil
}

func (e *enemyBoomer) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	// Show Enemy Images
	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()

	var ofsy int
	if e.state == boomerStateMove {
		c := e.count % boomerActNextStepCount
		if c == 0 || c == boomerActNextStepCount/2 {
			ofsy = e.prevOfsY
		} else {
			ofsy = battlecommon.GetOffset(e.nextY, e.pm.Pos.Y, e.prevY, c, boomerActNextStepCount, field.PanelSize.Y)
			e.prevOfsY = ofsy
		}
	}
	dxlib.DrawRotaGraph(view.X, view.Y+ofsy, 1, 0, *img, true)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40+ofsy, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyBoomer) DamageProc(dm *damage.Damage) bool {
	return damageProc(dm, &e.pm)
}

func (e *enemyBoomer) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		Pos:      e.pm.Pos,
		AnimType: anim.AnimTypeObject,
	}
}

func (e *enemyBoomer) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyBoomer) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (e *enemyBoomer) getCurrentImagePointer() *int {
	if e.state == boomerStateAtk {
		n := (e.count / delayBoomerAtk)
		if n >= len(e.atk.images) {
			n = len(e.atk.images) - 1
		}
		return &e.atk.images[n]
	}

	n := (e.count / delayBoomerMove) % len(e.imgMove)
	return &e.imgMove[n]
}

func (e *enemyBoomer) setState(next int) {
	e.state = next
	e.count = 0
}

func (a *boomerAtk) Init() {
	a.atkID = ""
	a.count = 0
}

func (a *boomerAtk) Process() bool {
	if a.count == 0 {
		a.atkID = anim.New(skill.Get(
			skill.SkillBoomerang,
			skill.Argument{
				OwnerID:    a.ownerID,
				Power:      20,
				TargetType: damage.TargetPlayer,
			},
		))
	}

	a.count++

	return !anim.IsProcessing(a.atkID)
}
