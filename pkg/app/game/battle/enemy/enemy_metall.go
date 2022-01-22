package enemy

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayMetallAtk = 3
)

var (
	metallActQueue = []string{}
)

type metallAtk struct {
	id      string
	ownerID string
	count   int
	images  []int
}

type enemyMetall struct {
	pm              EnemyParam
	imgMove         []int
	count           int
	moveFailedCount int
	atkID           string
	atk             metallAtk
}

func (e *enemyMetall) Init(objID string) error {
	name, ext := GetStandImageFile(IDMetall)

	e.pm.ObjectID = objID
	e.atk.id = uuid.New().String()
	e.imgMove = make([]int, 1)
	fname := name + "_move" + ext
	e.imgMove[0] = dxlib.LoadGraph(fname)
	if e.imgMove[0] == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}
	e.atk.images = make([]int, 15)
	fname = name + "_atk" + ext
	if res := dxlib.LoadDivGraph(fname, 15, 15, 1, 100, 140, e.atk.images); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	metallActQueue = append(metallActQueue, objID)

	return nil
}

func (e *enemyMetall) End() {
	dxlib.DeleteGraph(e.imgMove[0])
	for _, img := range e.atk.images {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyMetall) Process() (bool, error) {
	if e.pm.HP <= 0 {
		// Delete Animation
		img := &e.imgMove[0]
		if e.atkID != "" {
			img = &e.atk.images[e.atk.GetImageNo()]
		}
		battlecommon.NewDelete(*img, e.pm.Pos, false)
		anim.New(effect.Get(effect.TypeExplode, e.pm.Pos, 0))
		*img = -1 // DeleteGraph at delete animation

		// Delete from act queue
		for i, id := range metallActQueue {
			if e.pm.ObjectID == id {
				metallActQueue = append(metallActQueue[:i], metallActQueue[i+1:]...)
				break
			}
		}
		return true, nil
	}

	const waitCount = 1 * 60
	const actionInterval = 1 * 60
	const forceAttackCount = 3

	if e.atkID != "" {
		// Anim end
		if !anim.IsProcessing(e.atkID) {
			metallActQueue = metallActQueue[1:]
			metallActQueue = append(metallActQueue, e.pm.ObjectID)

			e.atkID = ""
			e.count = 0
		}
		return false, nil
	}

	if metallActQueue[0] != e.pm.ObjectID {
		// other metall is acting
		return false, nil
	}

	e.count++

	// Metall Actions
	if e.count < waitCount {
		return false, nil
	}

	if e.count%actionInterval == 0 {
		pos := objanim.GetObjPos(e.pm.PlayerID)
		if pos.Y == e.pm.Pos.Y || e.moveFailedCount >= forceAttackCount {
			// Attack
			e.atk.count = 0
			e.atk.ownerID = e.pm.ObjectID
			e.atkID = anim.New(&e.atk)
			e.moveFailedCount = 0
		} else {
			// Move
			moved := false
			if pos.Y > e.pm.Pos.Y {
				moved = battlecommon.MoveObject(&e.pm.Pos, common.DirectDown, field.PanelTypeEnemy, true, field.GetPanelInfo)
			} else {
				moved = battlecommon.MoveObject(&e.pm.Pos, common.DirectUp, field.PanelTypeEnemy, true, field.GetPanelInfo)
			}
			if moved {
				e.moveFailedCount = 0
			} else {
				e.moveFailedCount++
			}
		}
	}

	return false, nil
}

func (e *enemyMetall) Draw() {
	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.imgMove[0]
	if e.atkID != "" {
		img = e.atk.images[e.atk.GetImageNo()]
	}
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, img, true)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyMetall) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}
	if dm.TargetType&damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.Pos, 5))
		return true
	}
	return false
}

func (e *enemyMetall) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		Pos:      e.pm.Pos,
		AnimType: anim.AnimTypeObject,
	}
}

func (e *enemyMetall) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (a *metallAtk) Draw() {
	// Nothing to do
}

func (a *metallAtk) Process() (bool, error) {
	a.count++

	if a.count == delayMetallAtk*10 {
		anim.New(skill.Get(skill.SkillShockWave, skill.Argument{
			OwnerID:    a.ownerID,
			Power:      10,
			TargetType: damage.TargetPlayer,
		}))
	}

	return a.count >= (len(a.images) * delayMetallAtk), nil
}

func (a *metallAtk) GetParam() anim.Param {
	return anim.Param{
		ObjID:    a.id,
		AnimType: anim.AnimTypeEffect,
	}
}

func (a *metallAtk) GetImageNo() int {
	n := a.count / delayMetallAtk
	if n >= len(a.images) {
		n = len(a.images) - 1
	}
	return n
}
