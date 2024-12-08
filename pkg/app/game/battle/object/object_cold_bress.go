package object

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
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
	delayColdBress         = 6
	coldBressNextStepCount = 80
)

type ColdBress struct {
	pm       ObjectParam
	images   []int
	count    int
	damageID string
	next     point.Point
	prev     point.Point
}

func (o *ColdBress) Init(ownerID string, initParam ObjectParam) error {
	o.pm = initParam
	o.pm.objectID = uuid.New().String()
	o.pm.xFlip = o.pm.OnwerCharType == objanim.ObjTypePlayer

	o.next = o.pm.Pos
	o.prev = point.Point{X: o.pm.Pos.X + 1, Y: o.pm.Pos.Y}

	// Load Images
	o.images = make([]int, 7)
	fname := config.ImagePath + "battle/skill/コールドマン_ブレス.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 76, 78, o.images); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	return nil
}

func (o *ColdBress) End() {
	// Delete Images
	for _, img := range o.images {
		dxlib.DeleteGraph(img)
	}
}

func (o *ColdBress) Update() (bool, error) {
	if o.pm.HP <= 0 {
		return true, nil
	}

	// キャラにヒット時はダメージを与えて消える
	if o.count%coldBressNextStepCount == 1 {
		if o.damageID != "" {
			if !localanim.DamageManager().Exists(o.damageID) {
				// attack hit to target
				return true, nil
			}
		}
	}

	if o.count%coldBressNextStepCount == 0 {
		if o.count != 0 {
			// Update current pos
			o.prev = o.pm.Pos
			o.pm.Pos = o.next
		}

		if o.pm.Pos.X < 0 || o.pm.Pos.X > battlecommon.FieldNum.X || o.pm.Pos.Y < 0 || o.pm.Pos.Y > battlecommon.FieldNum.Y {
			return true, nil
		}

		// Add damage
		target := damage.TargetPlayer
		if o.pm.OnwerCharType == objanim.ObjTypePlayer {
			target = damage.TargetEnemy
		}

		o.damageID = localanim.DamageManager().New(damage.Damage{
			DamageType:    damage.TypePosition,
			Pos:           o.pm.Pos,
			Power:         10,
			TTL:           coldBressNextStepCount,
			TargetObjType: target,
			HitEffectType: resources.EffectTypeNone,
			ShowHitArea:   false,
			StrengthType:  damage.StrengthHigh,
			Element:       damage.ElementNone,
		})

		// Update next pos
		left := point.Point{X: o.next.X - 1, Y: o.next.Y}
		if o.checkMove(left) {
			o.next = left
		} else {
			up := point.Point{X: o.next.X, Y: o.next.Y - 1}
			down := point.Point{X: o.next.X, Y: o.next.Y + 1}
			if up.Y >= 0 && o.checkMove(up) {
				o.next = up
			} else if down.Y < battlecommon.FieldNum.Y && o.checkMove(down) {
				o.next = down
			}
		}
	}

	o.count++
	return false, nil
}

func (o *ColdBress) Draw() {
	cnt := o.count % coldBressNextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	view := battlecommon.ViewPos(o.pm.Pos)

	n := (o.count / delayColdBress) % len(o.images)
	ofsx := battlecommon.GetOffset(o.next.X, o.pm.Pos.X, o.prev.X, cnt, coldBressNextStepCount, battlecommon.PanelSize.X)
	ofsy := battlecommon.GetOffset(o.next.Y, o.pm.Pos.Y, o.prev.Y, cnt, coldBressNextStepCount, battlecommon.PanelSize.Y)

	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+16+ofsy, 1, 0, o.images[n], true, dxlib.OptXReverse(o.pm.xFlip))
}

func (o *ColdBress) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}

	target := damage.TargetEnemy
	if o.pm.OnwerCharType == objanim.ObjTypePlayer {
		target = damage.TargetPlayer
	}

	if dm.TargetObjType&target != 0 {
		o.pm.HP -= dm.Power

		localanim.AnimNew(effect.Get(dm.HitEffectType, o.pm.Pos, 5))
		return true
	}

	return false
}

func (o *ColdBress) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    o.pm.objectID,
			Pos:      o.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: o.pm.HP,
	}
}

func (o *ColdBress) GetObjectType() int {
	return objanim.ObjTypeNone
}

func (o *ColdBress) MakeInvisible(count int) {}

func (o *ColdBress) checkMove(next point.Point) bool {
	objID := localanim.ObjAnimExistsObject(next)
	if objID == "" {
		return true
	}

	target := objanim.ObjTypePlayer
	if o.pm.OnwerCharType == objanim.ObjTypePlayer {
		target = objanim.ObjTypeEnemy
	}
	objs := localanim.ObjAnimGetObjs(objanim.Filter{ObjType: target})
	for _, obj := range objs {
		if obj.ObjID == objID {
			return true
		}
	}

	return false
}
