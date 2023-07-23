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
	delayColdBress = 3
)

type ColdBress struct {
	pm     ObjectParam
	images []int
	count  int
}

func (o *ColdBress) Init(ownerID string, initParam ObjectParam) error {
	o.pm = initParam
	o.pm.objectID = uuid.New().String()
	o.pm.xFlip = o.pm.OnwerCharType == objanim.ObjTypePlayer

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
		// TODO delete animation
		return true, nil
	}

	o.count++

	return false, nil
}

func (o *ColdBress) Draw() {
	view := battlecommon.ViewPos(o.pm.Pos)

	opt := dxlib.DrawRotaGraphOption{}
	if o.pm.xFlip {
		f := int32(dxlib.TRUE)
		opt.ReverseXFlag = &f
	}

	n := o.count / delayColdBress
	if n > len(o.images)-1 {
		n = len(o.images) - 1
	}
	dxlib.DrawRotaGraph(view.X, view.Y+16, 1, 0, o.images[n], true, opt)
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
