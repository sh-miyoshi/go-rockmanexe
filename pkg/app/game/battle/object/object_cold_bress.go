package object

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
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
	next     common.Point
	prev     common.Point
}

func (o *ColdBress) Init(ownerID string, initParam ObjectParam) error {
	o.pm = initParam
	o.pm.objectID = uuid.New().String()
	o.pm.xFlip = o.pm.OnwerCharType == objanim.ObjTypePlayer

	o.next = o.pm.Pos
	o.prev = common.Point{X: o.pm.Pos.X + 1, Y: o.pm.Pos.Y}

	// Load Images
	o.images = make([]int, 7)
	fname := common.ImagePath + "battle/skill/コールドマン_ブレス.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 76, 78, o.images); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (o *ColdBress) End() {
	// Delete Images
	for _, img := range o.images {
		dxlib.DeleteGraph(img)
	}
}

func (o *ColdBress) Process() (bool, error) {
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

		// Update next pos
		// TODO: 基本left, 移動できないなら上下
		o.next.X--
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

	opt := dxlib.DrawRotaGraphOption{}
	if o.pm.xFlip {
		f := int32(dxlib.TRUE)
		opt.ReverseXFlag = &f
	}

	n := (o.count / delayColdBress) % len(o.images)
	ofsx := battlecommon.GetOffset(o.next.X, o.pm.Pos.X, o.prev.X, cnt, coldBressNextStepCount, battlecommon.PanelSize.X)
	ofsy := battlecommon.GetOffset(o.next.Y, o.pm.Pos.Y, o.prev.Y, cnt, coldBressNextStepCount, battlecommon.PanelSize.Y)

	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+16+ofsy, 1, 0, o.images[n], true, opt)
}

func (o *ColdBress) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}

	target := damage.TargetEnemy
	if o.pm.OnwerCharType == objanim.ObjTypePlayer {
		target = damage.TargetPlayer
	}

	if dm.TargetType&target != 0 {
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
