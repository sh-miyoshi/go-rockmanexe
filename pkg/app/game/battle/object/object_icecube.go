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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayIceCubeCreate = 2
)

type IceCube struct {
	pm     ObjectParam
	images []int
	count  int
}

func (o *IceCube) Init(ownerID string, initParam ObjectParam) error {
	o.pm = initParam
	o.pm.objectID = uuid.New().String()
	o.pm.xFlip = o.pm.OnwerCharType == objanim.ObjTypePlayer

	// Load Images
	o.images = make([]int, 6)
	fname := common.ImagePath + "battle/character/アイスキューブ.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 76, 90, o.images); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (o *IceCube) End() {
	// Delete Images
	for _, img := range o.images {
		dxlib.DeleteGraph(img)
	}
}

func (o *IceCube) Process() (bool, error) {
	if o.pm.HP <= 0 {
		// TODO delete animation
		return true, nil
	}

	o.count++

	return false, nil
}

func (o *IceCube) Draw() {
	view := battlecommon.ViewPos(o.pm.Pos)

	opt := dxlib.DrawRotaGraphOption{}
	if o.pm.xFlip {
		f := int32(dxlib.TRUE)
		opt.ReverseXFlag = &f
	}

	n := o.count / delayIceCubeCreate
	if n > len(o.images)-1 {
		n = len(o.images) - 1
	}
	dxlib.DrawRotaGraph(view.X, view.Y+16, 1, 0, o.images[n], true, opt)
}

func (o *IceCube) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}

	target := damage.TargetEnemy
	if o.pm.OnwerCharType == objanim.ObjTypePlayer {
		target = damage.TargetPlayer
	}

	if dm.TargetType&target != 0 {
		o.pm.HP -= dm.Power

		for i := 0; i < dm.PushLeft; i++ {
			if !battlecommon.MoveObject(&o.pm.Pos, common.DirectLeft, battlecommon.PanelTypeEnemy, true, field.GetPanelInfo) {
				break
			}
		}
		for i := 0; i < dm.PushRight; i++ {
			if !battlecommon.MoveObject(&o.pm.Pos, common.DirectRight, battlecommon.PanelTypeEnemy, true, field.GetPanelInfo) {
				break
			}
		}

		localanim.AnimNew(effect.Get(dm.HitEffectType, o.pm.Pos, 5))
		return true
	}

	return false
}

func (o *IceCube) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    o.pm.objectID,
			Pos:      o.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: o.pm.HP,
	}
}

func (o *IceCube) GetObjectType() int {
	return objanim.ObjTypeNone
}

func (o *IceCube) MakeInvisible(count int) {}