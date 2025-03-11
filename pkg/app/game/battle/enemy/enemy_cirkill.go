package enemy

import (
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
	delayCirkillMove   = 8
	delayCirkillAttack = 6

	cirkillMoveNextStepCount = 40
	initialWaitCount         = 60
)

type cirKillAttack struct {
	ownerID   string
	count     int
	attacking bool
	images    []int
	animID    string
}

type enemyCirKill struct {
	pm      EnemyParam
	atk     cirKillAttack
	imgMove []int
	count   int

	next point.Point
	prev point.Point
}

func (e *enemyCirKill) Init(objID string) error {
	e.pm.ObjectID = objID
	e.next = e.getNextPos()
	e.prev = e.pm.Pos
	e.count = e.pm.ActNo
	e.atk.ownerID = objID

	// Load Images
	name, ext := GetStandImageFile(IDCirKill)
	e.imgMove = make([]int, 4)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 96, e.imgMove); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	e.atk.images = make([]int, 2)
	fname = name + "_atk" + ext
	if res := dxlib.LoadDivGraph(fname, 2, 2, 1, 140, 104, e.atk.images); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyCirKill) End() {
	// Delete Images
	for _, img := range e.imgMove {
		dxlib.DeleteGraph(img)
	}
	for _, img := range e.atk.images {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyCirKill) Update() (bool, error) {
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

	e.count++

	e.atk.Update()

	if e.count < initialWaitCount {
		return false, nil
	}

	cnt := e.count % cirkillMoveNextStepCount
	if cnt == 0 {
		// 実際に移動
		e.prev = e.pm.Pos
		if battlecommon.MoveObjectDirect(&e.pm.Pos, e.next, battlecommon.PanelTypeEnemy, true, field.GetPanelInfo) {
			e.next = e.getNextPos()
		}
	} else if cnt == cirkillMoveNextStepCount/2 {
		pos := localanim.ObjAnimGetObjPos(e.pm.PlayerID)
		if e.pm.Pos.Y == pos.Y {
			e.atk.Set()
		}
	}

	return false, nil
}

func (e *enemyCirKill) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	ofsx := 0
	ofsy := 0
	if e.count > initialWaitCount {
		c := e.count % cirkillMoveNextStepCount
		ofsx = battlecommon.GetOffset(e.next.X, e.pm.Pos.X, e.prev.X, c, cirkillMoveNextStepCount, battlecommon.PanelSize.X)
		ofsy = battlecommon.GetOffset(e.next.Y, e.pm.Pos.Y, e.prev.Y, c, cirkillMoveNextStepCount, battlecommon.PanelSize.Y)
	}

	ofsx -= 10
	ofsy += 15

	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()

	if e.atk.attacking {
		dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, *img, true, dxlib.OptXReverse(true))
	} else {
		dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, *img, true, dxlib.OptXReverse(true))
		drawParalysis(view.X+ofsx, view.Y+ofsy, *img, e.pm.ParalyzedCount)
	}

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X+ofsx+10, view.Y+ofsy+25, e.pm.HP, draw.NumberOption{
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
			ObjID: e.pm.ObjectID,
			Pos:   e.pm.Pos,
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

func (e *enemyCirKill) AddBarrier(hp int) {}

func (e *enemyCirKill) getCurrentImagePointer() *int {
	if e.atk.attacking {
		n := (e.count / delayCirkillAttack)
		if n >= len(e.atk.images) {
			n = len(e.atk.images) - 1
		}
		return &e.atk.images[n]
	}

	n := (e.count / delayCirkillMove) % len(e.imgMove)
	return &e.imgMove[n]
}

func (e *enemyCirKill) getNextPos() point.Point {
	// 外周時計回り
	if e.pm.Pos.X == battlecommon.FieldNum.X/2 {
		if e.pm.Pos.Y != 0 {
			return point.Point{X: e.pm.Pos.X, Y: e.pm.Pos.Y - 1}
		}
	}
	if e.pm.Pos.X == battlecommon.FieldNum.X-1 {
		if e.pm.Pos.Y != battlecommon.FieldNum.Y-1 {
			return point.Point{X: e.pm.Pos.X, Y: e.pm.Pos.Y + 1}
		}
	}
	if e.pm.Pos.Y == 0 {
		return point.Point{X: e.pm.Pos.X + 1, Y: e.pm.Pos.Y}
	}
	if e.pm.Pos.Y == battlecommon.FieldNum.Y-1 {
		return point.Point{X: e.pm.Pos.X - 1, Y: e.pm.Pos.Y}
	}

	return point.Point{}
}

func (a *cirKillAttack) Set() {
	if a.animID == "" {
		a.attacking = true
		a.count = 0
	}
}

func (a *cirKillAttack) Update() {
	if a.animID != "" && !localanim.AnimIsProcessing(a.animID) {
		a.animID = ""
	}

	if a.attacking {
		if a.count == 0 {
			a.animID = localanim.SkillAnimNew(skill.Get(resources.SkillCirkillShot, skillcore.Argument{
				OwnerID:    a.ownerID,
				Power:      10,
				TargetType: damage.TargetPlayer,
			}))
		}

		a.count++
		if a.count > delayCirkillAttack*len(a.images) {
			a.attacking = false
		}
	}
}
