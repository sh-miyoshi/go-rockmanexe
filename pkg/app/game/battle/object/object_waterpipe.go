package object

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
)

const (
	delayWaterPipeSet = 3
)

// type WaterPipeAtk struct {
// 	id      string
// 	ownerID string
// 	count   int
// 	images  []int32
// }

type WaterPipe struct {
	pm ObjectParam
	// atk    WaterPipeAtk
	imgSet []int32
	count  int
}

func (o *WaterPipe) Init(ownerID string, initParam ObjectParam) error {
	o.pm = initParam
	o.pm.ObjectID = uuid.New().String()
	// o.atk.ownerID = ownerID

	// Load Images
	o.imgSet = make([]int32, 4)
	fname := common.ImagePath + "battle/character/水道管_set.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 72, 88, o.imgSet); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}
	return nil
}

func (o *WaterPipe) End() {
	// Delete Images
	for _, img := range o.imgSet {
		dxlib.DeleteGraph(img)
	}
}

func (o *WaterPipe) Process() (bool, error) {
	o.count++
	return false, nil
}

func (o *WaterPipe) Draw() {
	x, y := battlecommon.ViewPos(o.pm.PosX, o.pm.PosY)

	// TODO: 攻撃時

	n := o.count / delayWaterPipeSet
	if n > len(o.imgSet)-1 {
		n = len(o.imgSet) - 1
	}
	dxlib.DrawRotaGraph(x-8, y+16, 1, 0, o.imgSet[n], dxlib.TRUE)
}

func (o *WaterPipe) DamageProc(dm *damage.Damage) bool {
	// TODO: はじかれエフェクト
	return false
}

func (o *WaterPipe) GetParam() anim.Param {
	return anim.Param{
		ObjID:    o.pm.ObjectID,
		PosX:     o.pm.PosX,
		PosY:     o.pm.PosY,
		AnimType: anim.AnimTypeObject,
	}
}

func (o *WaterPipe) GetObjectType() int {
	return objanim.ObjTypeNone
}
