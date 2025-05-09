package enemy

import (
	"math/rand"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
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
	delayGarooMove = 16
	delayGarooAtk  = 8

	garooInitWait = 40
	garooAtkStr   = "set_attack"
)

type garooAtk struct {
	id      string
	ownerID string
	count   int
	images  []int
	animMgr *manager.Manager
}

type enemyGaroo struct {
	pm        EnemyParam
	imgMove   []int
	count     int
	atkID     string
	atk       garooAtk
	moveNum   int
	targetPos point.Point
	waitCount int
	animMgr   *manager.Manager
}

func (e *enemyGaroo) Init(objID string, animMgr *manager.Manager) error {
	e.pm.ObjectID = objID
	e.pm.DamageElement = damage.ElementFire
	e.moveNum = 3
	e.targetPos = point.Point{X: -1, Y: -1}
	e.waitCount = garooInitWait
	e.animMgr = animMgr

	// Load Images
	name, ext := GetStandImageFile(IDGaroo)
	e.atk.id = uuid.New().String()
	e.atk.animMgr = animMgr
	e.imgMove = make([]int, 4)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 168, 132, e.imgMove); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}
	e.atk.images = make([]int, 6)
	fname = name + "_atk" + ext
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 168, 132, e.atk.images); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyGaroo) End() {
	// Delete Images
	for _, img := range e.imgMove {
		dxlib.DeleteGraph(img)
	}
	for _, img := range e.atk.images {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyGaroo) Update() (bool, error) {
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		deleteanim.New(*img, e.pm.Pos, false, e.animMgr)
		e.animMgr.EffectAnimNew(effect.Get(resources.EffectTypeExplode, e.pm.Pos, 0))
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}

	if e.pm.ParalyzedCount > 0 {
		e.pm.ParalyzedCount--
		return false, nil
	}

	if e.waitCount > 0 {
		e.waitCount--
		return false, nil
	}

	if e.atkID != "" {
		if e.atkID == garooAtkStr {
			e.atk.ownerID = e.pm.ObjectID
			e.atkID = e.animMgr.EffectAnimNew(&e.atk)
		}

		// Anim end
		if !e.animMgr.IsAnimProcessing(e.atkID) {
			e.atkID = ""
			e.waitCount = garooInitWait
		}
		return false, nil
	}

	e.count++

	const actionInterval = 70

	if e.count%actionInterval == 0 {
		if e.targetPos.X != -1 && e.targetPos.Y != -1 {
			if battlecommon.MoveObjectDirect(
				&e.pm.Pos,
				e.targetPos,
				battlecommon.PanelTypeEnemy,
				true,
				field.GetPanelInfo,
			) {
				e.targetPos = point.Point{X: -1, Y: -1}

				// Set attack
				e.moveNum = 2 + rand.Intn(3)
				e.atk.count = 0
				e.atkID = garooAtkStr
				e.waitCount = 30
				return false, nil
			}
		}

		for i := 0; i < 10; i++ {
			next := point.Point{
				X: rand.Intn(battlecommon.FieldNum.X/2) + battlecommon.FieldNum.X/2,
				Y: rand.Intn(battlecommon.FieldNum.Y),
			}
			if battlecommon.MoveObjectDirect(
				&e.pm.Pos,
				next,
				battlecommon.PanelTypeEnemy,
				true,
				field.GetPanelInfo,
			) {
				break
			}
		}
		e.moveNum--
		if e.moveNum <= 0 {
			objs := e.animMgr.ObjAnimGetObjs(objanim.Filter{ObjType: objanim.ObjTypePlayer})
			pos := point.Point{X: 1, Y: 1}
			if len(objs) > 0 {
				pos = objs[0].Pos
			}
			// set attack pos to {X: random, Y: same as player}
			rnd := rand.Intn(3)
			for i := 0; i < 3; i++ {
				pos = point.Point{X: (rnd+i)%3 + 3, Y: pos.Y}
				if battlecommon.MoveObjectDirect(&e.pm.Pos, pos, battlecommon.PanelTypeEnemy, false, field.GetPanelInfo) {
					e.targetPos = pos
				}
			}
		}

		return false, nil
	}

	return false, nil
}

func (e *enemyGaroo) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, *img, true, dxlib.OptXReverse(true))

	drawParalysis(view.X, view.Y, *img, e.pm.ParalyzedCount)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyGaroo) DamageProc(dm *damage.Damage) bool {
	return damageProc(dm, &e.pm)
}

func (e *enemyGaroo) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID: e.pm.ObjectID,
			Pos:   e.pm.Pos,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyGaroo) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyGaroo) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (e *enemyGaroo) AddBarrier(hp int) {}

func (e *enemyGaroo) SetCustomGaugeMax() {}

func (e *enemyGaroo) getCurrentImagePointer() *int {
	if e.atkID != "" {
		n := (e.atk.count / delayGarooAtk)
		if n >= len(e.atk.images) {
			n = len(e.atk.images) - 1
		}
		return &e.atk.images[n]
	}

	n := (e.count / delayGarooMove) % len(e.imgMove)
	return &e.imgMove[n]
}

func (a *garooAtk) Draw() {
	// Nothing to do
}

func (a *garooAtk) Update() (bool, error) {
	a.count++

	if a.count == delayGarooAtk*4 {
		a.animMgr.SkillAnimNew(skill.Get(resources.SkillGarooBreath, skillcore.Argument{
			OwnerID:    a.ownerID,
			Power:      10,
			TargetType: damage.TargetPlayer,
		}, a.animMgr))
	}

	return a.count >= (len(a.images) * delayGarooAtk), nil
}

func (a *garooAtk) GetParam() anim.Param {
	return anim.Param{
		ObjID: a.id,
	}
}
